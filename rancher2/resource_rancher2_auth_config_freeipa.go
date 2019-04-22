package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2AuthConfigFreeIpa() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigFreeIpaCreate,
		Read:   resourceRancher2AuthConfigFreeIpaRead,
		Update: resourceRancher2AuthConfigFreeIpaUpdate,
		Delete: resourceRancher2AuthConfigFreeIpaDelete,

		Schema: authConfigFreeIpaFields(),
	}
}

func resourceRancher2AuthConfigFreeIpaCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AuthConfigFreeIpaName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s %s", AuthConfigFreeIpaName, auth.Name)

	authFreeIpa, err := expandAuthConfigFreeIpa(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AuthConfigFreeIpaName, err)
	}

	// Checking if other auth config is enabled
	if authFreeIpa.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AuthConfigFreeIpaName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AuthConfigFreeIpaName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.FreeIpaConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authFreeIpa, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AuthConfigFreeIpaName, err)
	}

	return resourceRancher2AuthConfigFreeIpaRead(d, meta)
}

func resourceRancher2AuthConfigFreeIpaRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AuthConfigFreeIpaName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigFreeIpaName)
			d.SetId("")
			return nil
		}
		return err
	}

	authFreeIpa, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigFreeIpa(d, authFreeIpa.(*managementClient.LdapConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigFreeIpaUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AuthConfigFreeIpaName)

	return resourceRancher2AuthConfigFreeIpaCreate(d, meta)
}

func resourceRancher2AuthConfigFreeIpaDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AuthConfigFreeIpaName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AuthConfigFreeIpaName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AuthConfigFreeIpaName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AuthConfigFreeIpaName, err)
		}
	}

	d.SetId("")
	return nil
}
