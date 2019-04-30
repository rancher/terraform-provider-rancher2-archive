package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

func resourceRancher2Namespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NamespaceCreate,
		Read:   resourceRancher2NamespaceRead,
		Update: resourceRancher2NamespaceUpdate,
		Delete: resourceRancher2NamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NamespaceImport,
		},

		Schema: namespaceFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	clusterID, err := clusterIDFromProjectID(d.Get("project_id").(string))
	if err != nil {
		return err
	}

	active, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		return err
	}
	if !active {
		return fmt.Errorf("[ERROR] Creating namespace: Cluster ID %s is not active", clusterID)
	}

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns := expandNamespace(d)

	log.Printf("[INFO] Creating Namespace %s on Cluster ID %s", ns.Name, clusterID)

	newNs, err := client.Namespace.Create(ns)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"activating"},
		Target:     []string{"active"},
		Refresh:    namespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be created: %s", newNs.ID, waitErr)
	}

	err = flattenNamespace(d, newNs)
	if err != nil {
		return err
	}

	return resourceRancher2NamespaceRead(d, meta)
}

func resourceRancher2NamespaceRead(d *schema.ResourceData, meta interface{}) error {
	clusterID, _ := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Refreshing Namespace ID %s", d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Namespace ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenNamespace(d, ns)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2NamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterID, projectID := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Updating Namespace ID %s", d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		return err
	}

	readClusterID, readProjectID := splitProjectID(ns.ProjectID)

	if projectID != readProjectID && (clusterID == readClusterID || readClusterID == "") {
		log.Printf("[INFO] Moving Namespace ID %s to project %s", d.Id(), projectID)
		nsMove := &clusterClient.NamespaceMove{
			ProjectID: projectID,
		}

		err = client.Namespace.ActionMove(ns, nsMove)
		if err != nil {
			return err
		}
	}

	resourceQuota := expandNamespaceResourceQuota(d.Get("resource_quota").([]interface{}))

	update := map[string]interface{}{
		"description":   d.Get("description").(string),
		"resourceQuota": resourceQuota,
		"annotations":   toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":        toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNs, err := client.Namespace.Update(ns, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    namespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be updated: %s", newNs.ID, waitErr)
	}

	err = flattenNamespace(d, newNs)
	if err != nil {
		return err
	}

	return resourceRancher2NamespaceRead(d, meta)
}

func resourceRancher2NamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	clusterID, _ := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Deleting Namespace ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return err
	}

	ns, err := client.Namespace.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Namespace ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Namespace.Delete(ns)
	if err != nil {
		return fmt.Errorf("Error removing Namespace: %s", err)
	}

	log.Printf("[DEBUG] Waiting for namespace (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    namespaceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// namespaceStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Namespace.
func namespaceStateRefreshFunc(client *clusterClient.Client, nsID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Namespace.ByID(nsID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
