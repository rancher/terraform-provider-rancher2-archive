package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Schemas

func projectRoleTemplateBindingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"role_template_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"group_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"group_principal_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_principal_id": {
			Type:     schema.TypeString,
			Optional: true,
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
