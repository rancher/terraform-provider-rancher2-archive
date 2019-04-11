package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2NodeTemplateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	nodeTemplate := &NodeTemplate{}
	err = client.APIBaseClient.ByID(managementClient.NodeTemplateType, d.Id(), nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNodeTemplate(d, nodeTemplate)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
