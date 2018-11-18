package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

func init() {
	descriptions = map[string]string{
		"name": "Name of the k8s namespace managed by rancher v2",

		"cluster_id": "Cluster ID",

		"project_id": "Project ID where k8s namespace belongs",

		"description": "Description of the k8s namespace managed by rancher v2",

		"resource_quota_template_id": "Resource quota template id to apply on k8s namespace",

		"annotations": "Annotations of the k8s namespace managed by rancher v2",

		"labels": "Labels of the k8s namespace managed by rancher v2",
	}
}

//Schemas
func namespaceMoveFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: descriptions["cluster_id"],
		},
		"project_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     nil,
			Description: descriptions["project_id"],
		},
		"name": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: descriptions["name"],
		},
		"description": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: descriptions["description"],
		},
		"annotations": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["annotations"],
		},
		"labels": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: descriptions["labels"],
		},
	}
	return s
}

// Flatteners
func flattenNamespaceMove(clusterID string, d *schema.ResourceData, in *clusterClient.Namespace) error {
	if in == nil {
		return nil
	}

	d.Set("cluster_id", clusterID)
	d.SetId(in.ID)

	err := d.Set("project_id", in.ProjectID)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("description", in.Description)
	if err != nil {
		return err
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandNamespaceMove(in *schema.ResourceData) *clusterClient.Namespace {
	obj := &clusterClient.Namespace{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ProjectID = in.Get("project_id").(string)
	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func resourceRancher2NamespaceMove() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NamespaceMoveCreate,
		Read:   resourceRancher2NamespaceMoveRead,
		Update: resourceRancher2NamespaceMoveUpdate,
		Delete: resourceRancher2NamespaceMoveDelete,
		Schema: namespaceMoveFields(),
	}
}

func resourceRancher2NamespaceMoveUpdateProjectID(newProjectID string, d *schema.ResourceData, meta interface{}) error {

	clusterID := d.Get("cluster_id").(string)

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

	ns, err := client.Namespace.ByID(d.Get("name").(string))
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"projectId":   newProjectID,
		"description": ns.Description,
		"annotations": ns.Annotations,
		"labels":      ns.Labels,
	}

	log.Printf("[INFO] Move Namespace %s from Project %s to %s", ns.Name, ns.ProjectID, newProjectID)

	// Move Namespace by updating Project ID
	newNs, err := client.Namespace.Update(ns, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    namespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for namespace (%s) to be updated: %s", newNs.ID, waitErr)
	}

	err = flattenNamespaceMove(clusterID, d, newNs)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2NamespaceMoveCreate(d *schema.ResourceData, meta interface{}) error {
	newProjectID := d.Get("project_id").(string)
	log.Printf("[INFO] Move Namespace %s into Project %s", d.Get("name"), d.Get("project_id"))
	err := resourceRancher2NamespaceMoveUpdateProjectID(newProjectID, d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceRancher2NamespaceMoveRead(d *schema.ResourceData, meta interface{}) error {

	clusterID := d.Get("cluster_id").(string)

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

	err = flattenNamespaceMove(clusterID, d, ns)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2NamespaceMoveUpdate(d *schema.ResourceData, meta interface{}) error {
	err := resourceRancher2NamespaceMoveCreate(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func resourceRancher2NamespaceMoveDelete(d *schema.ResourceData, meta interface{}) error {
	// Set empty string Project ID to move Namespace into None Project
	newProjectID := ""
	log.Printf("[INFO] Move Namespace %s into Project None", d.Get("name"))
	resourceRancher2NamespaceMoveUpdateProjectID(newProjectID, d, meta)
	return nil
}
