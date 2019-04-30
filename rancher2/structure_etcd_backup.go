package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenEtcdBackup(d *schema.ResourceData, in *managementClient.EtcdBackup) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	if v, ok := d.Get("backup_config.0.s3_backup_config.0.secret_key").(string); ok && len(v) > 0 {
		in.BackupConfig.S3BackupConfig.SecretKey = d.Get("backup_config.0.s3_backup_config.0.secret_key").(string)
	}

	err := d.Set("backup_config", flattenClusterRKEConfigServicesEtcdBackupConfig(in.BackupConfig))
	if err != nil {
		return err
	}

	err = d.Set("cluster_id", in.ClusterID)
	if err != nil {
		return err
	}

	err = d.Set("filename", in.Filename)
	if err != nil {
		return err
	}

	err = d.Set("manual", in.Manual)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("namespace_id", in.NamespaceId)
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

func expandEtcdBackup(in *schema.ResourceData) *managementClient.EtcdBackup {
	obj := &managementClient.EtcdBackup{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("backup_config").([]interface{}); ok && len(v) > 0 {
		obj.BackupConfig = expandClusterRKEConfigServicesEtcdBackupConfig(v)
	}

	obj.ClusterID = in.Get("cluster_id").(string)

	if v, ok := in.Get("filename").(string); ok && len(v) > 0 {
		obj.Filename = v
	}

	if v, ok := in.Get("manual").(bool); ok {
		obj.Manual = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("namespace_id").(string); ok && len(v) > 0 {
		obj.NamespaceId = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
