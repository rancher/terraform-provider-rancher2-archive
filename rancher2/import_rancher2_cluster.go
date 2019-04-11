package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRancher2ClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2ClusterRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
