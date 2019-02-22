package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const clusterRegistrationTokenName = "system"

// Schema

func clusterRegistationTokenFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"cluster_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"insecure_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"manifest_url": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"windows_node_command": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

// Flatteners

func flattenClusterRegistationToken(in *managementClient.ClusterRegistrationToken) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["id"] = in.ID
	obj["cluster_id"] = in.ClusterID
	obj["name"] = clusterRegistrationTokenName
	obj["command"] = in.Command
	obj["insecure_command"] = in.InsecureCommand
	obj["manifest_url"] = in.ManifestURL
	obj["node_command"] = in.NodeCommand
	obj["windows_node_command"] = in.WindowsNodeCommand
	obj["annotations"] = toMapInterface(in.Annotations)
	obj["labels"] = toMapInterface(in.Labels)

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRegistationToken(p []interface{}, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("[ERROR] Expanding Cluster Registration Token: Cluster id is nil")
	}

	obj := &managementClient.ClusterRegistrationToken{}
	obj.ClusterID = clusterID
	obj.Name = clusterRegistrationTokenName

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func findFlattenClusterRegistrationToken(client *managementClient.Client, clusterID string) ([]interface{}, error) {
	clusterReg, err := findClusterRegistrationToken(client, clusterID)
	if err != nil {
		return []interface{}{}, err
	}

	return flattenClusterRegistationToken(clusterReg)
}

func findClusterRegistrationToken(client *managementClient.Client, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	regTokenID := clusterID + ":" + clusterRegistrationTokenName
	regToken, err := client.ClusterRegistrationToken.ByID(regTokenID)

	if err != nil {
		if IsNotFound(err) {
			regToken, err = expandClusterRegistationToken([]interface{}{}, clusterID)
			if err != nil {
				return nil, err
			}
			newRegToken, err := client.ClusterRegistrationToken.Create(regToken)
			if err != nil {
				return nil, err
			}
			return newRegToken, nil
		}
		return nil, err
	}

	return regToken, nil
}
