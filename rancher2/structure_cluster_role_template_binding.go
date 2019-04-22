package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterRoleTemplateBinding(d *schema.ResourceData, in *managementClient.ClusterRoleTemplateBinding) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("cluster_id", in.ClusterID)
	if err != nil {
		return err
	}

	err = d.Set("role_template_id", in.RoleTemplateID)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("group_id", in.GroupID)
	if err != nil {
		return err
	}

	err = d.Set("group_principal_id", in.GroupPrincipalID)
	if err != nil {
		return err
	}

	err = d.Set("user_id", in.UserID)
	if err != nil {
		return err
	}

	err = d.Set("user_principal_id", in.UserPrincipalID)
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

func expandClusterRoleTemplateBinding(in *schema.ResourceData) *managementClient.ClusterRoleTemplateBinding {
	obj := &managementClient.ClusterRoleTemplateBinding{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.RoleTemplateID = in.Get("role_template_id").(string)
	obj.Name = in.Get("name").(string)
	obj.GroupID = in.Get("group_id").(string)
	obj.GroupPrincipalID = in.Get("group_principal_id").(string)
	obj.UserID = in.Get("user_id").(string)
	obj.UserPrincipalID = in.Get("user_principal_id").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
