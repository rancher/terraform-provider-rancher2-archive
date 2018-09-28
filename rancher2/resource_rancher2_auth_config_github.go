package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const GithubConfigName = "github"

//Schemas

func authConfigGithubFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"client_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"client_secret": {
			Type:     schema.TypeString,
			Required: true,
		},
		"code": {
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "github.com",
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}

// Flatteners

func flattenAuthConfigGithub(d *schema.ResourceData, in *managementClient.GithubConfig) error {
	d.SetId(GithubConfigName)

	err := d.Set("name", GithubConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.GithubConfigType)
	if err != nil {
		return err
	}

	err = d.Set("access_mode", in.AccessMode)
	if err != nil {
		return err
	}
	err = d.Set("allowed_principal_ids", toArrayInterface(in.AllowedPrincipalIDs))
	if err != nil {
		return err
	}
	err = d.Set("enabled", in.Enabled)
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

	err = d.Set("client_id", in.ClientID)
	if err != nil {
		return err
	}
	err = d.Set("hostname", in.Hostname)
	if err != nil {
		return err
	}
	err = d.Set("tls", in.TLS)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigGithub(in *schema.ResourceData) (*managementClient.GithubConfig, error) {
	obj := &managementClient.GithubConfig{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding Github Auth Config: Input ResourceData is nil")
	}

	obj.Name = GithubConfigName
	obj.Type = managementClient.GithubConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if v, ok := in.Get("enabled").(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("client_id").(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in.Get("client_secret").(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in.Get("hostname").(string); ok && len(v) > 0 {
		obj.Hostname = v
	}

	if v, ok := in.Get("tls").(bool); ok {
		obj.TLS = v
	}

	return obj, nil
}

func resourceRancher2AuthConfigGithub() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigGithubCreate,
		Read:   resourceRancher2AuthConfigGithubRead,
		Update: resourceRancher2AuthConfigGithubUpdate,
		Delete: resourceRancher2AuthConfigGithubDelete,
		//Importer: &schema.ResourceImporter{
		//	State: resourceRancher2AuthConfigGithubImport,
		//},

		Schema: authConfigGithubFields(),
	}
}

func resourceRancher2AuthConfigGithubCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(GithubConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config Github ID %s err=%s", GithubConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config Github %s", auth.Name)

	authGithub, err := expandAuthConfigGithub(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config Github ID %s err=%s", GithubConfigName, err)
	}

	authGithubTestAndApply := managementClient.GithubConfigApplyInput{
		GithubConfig: authGithub,
		Enabled:      authGithub.Enabled,
		Code:         d.Get("code").(string),
	}

	err = client.Post(auth.Actions["testAndApply"], authGithubTestAndApply, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config Github testAndApply [%s] %s", auth.Actions["testAndApply"], err)
	}

	return resourceRancher2AuthConfigGithubRead(d, meta)
}

func resourceRancher2AuthConfigGithubRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config Github ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(GithubConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config Github ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	authGithub, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigGithub(d, authGithub.(*managementClient.GithubConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigGithubUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config Github ID %s", d.Id())

	return resourceRancher2AuthConfigGithubCreate(d, meta)
}

func resourceRancher2AuthConfigGithubDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config Github ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(GithubConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config Github ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Post(auth.Actions["disable"], nil, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config Github disable [%s] %s", auth.Actions["disable"], err)
	}

	d.SetId("")
	return nil
}

/*
func resourceRancher2AuthConfigGithubImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	auth, err := client.AuthConfig.ByID(GithubConfigName)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	authGithub, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenAuthConfigGithub(d, authGithub.(*managementClient.GithubConfig))
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
*/
