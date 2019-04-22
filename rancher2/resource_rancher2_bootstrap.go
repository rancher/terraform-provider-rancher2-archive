package rancher2

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRancher2Bootstrap() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2BootstrapCreate,
		Read:   resourceRancher2BootstrapRead,
		Update: resourceRancher2BootstrapUpdate,
		Delete: resourceRancher2BootstrapDelete,
		Schema: bootstrapFields(),
	}
}

func resourceRancher2BootstrapCreate(d *schema.ResourceData, meta interface{}) error {
	if !meta.(*Config).Bootstrap {
		return fmt.Errorf("[ERROR] Resource rancher2_bootstrap just available on bootstrap mode")
	}

	err := bootstrapDoLogin(d, meta)
	if err != nil {
		return err
	}

	// Set user
	d.Set("user", bootstrapDefaultUser)

	// Set rancher url
	url := strings.TrimSuffix(meta.(*Config).URL, "/v3")
	err = meta.(*Config).SetSetting(bootstrapSettingURL, url)
	if err != nil {
		return err
	}

	// Set telemetry option
	telemetry := "out"
	if d.Get("telemetry").(bool) {
		telemetry = "in"
	}

	err = meta.(*Config).SetSetting(bootstrapSettingTelemetry, telemetry)
	if err != nil {
		return err
	}

	// Generate a new token
	tokenID, token, err := meta.(*Config).GenerateUserToken(bootstrapDefaultUser, bootstrapDefaultTokenDesc, d.Get("token_ttl").(int))
	if err != nil {
		return fmt.Errorf("[ERROR] Creating Admin token: %s", err)
	}

	// Update new tokenkey
	d.Set("token_id", tokenID)
	d.Set("token", token)
	err = meta.(*Config).UpdateToken(token)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Admin token: %s", err)
	}

	// Set admin user password
	pass := d.Get("password").(string)
	_, newPass, adminUser, err := meta.(*Config).SetUserPasswordByName(bootstrapDefaultUser, pass)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Admin password: %s", err)
	}

	d.Set("password", newPass)
	d.Set("current_password", newPass)

	// Set resource ID
	d.SetId(adminUser.ID)

	return resourceRancher2BootstrapRead(d, meta)
}

func resourceRancher2BootstrapRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing bootstrap")

	if !meta.(*Config).Bootstrap {
		return fmt.Errorf("[ERROR] Resource rancher2_bootstrap just available on bootstrap mode")
	}

	err := bootstrapDoLogin(d, meta)
	if err != nil {
		return err
	}

	// Check if token is expired
	expiredToken, err := meta.(*Config).IsTokenExpired(d.Get("token_id").(string))
	if err != nil {
		return err
	}

	d.Set("token_update", expiredToken)

	// Get rancher url
	url, err := meta.(*Config).GetSettingValue(bootstrapSettingURL)
	if err != nil {
		return err
	}

	d.Set("url", url)

	// Get telemetry
	telemetry, err := meta.(*Config).GetSettingValue(bootstrapSettingTelemetry)
	if err != nil {
		return err
	}

	if telemetry == "in" {
		d.Set("telemetry", true)
	} else {
		d.Set("telemetry", false)
	}

	return bootstrapCleanUpTempToken(d, meta)
}

func resourceRancher2BootstrapUpdate(d *schema.ResourceData, meta interface{}) error {
	err := bootstrapDoLogin(d, meta)
	if err != nil {
		return err
	}

	// Set user
	d.Set("user", bootstrapDefaultUser)

	// Set rancher url
	url := strings.TrimSuffix(meta.(*Config).URL, "/v3")
	err = meta.(*Config).SetSetting(bootstrapSettingURL, url)
	if err != nil {
		return err
	}

	// Set telemetry option
	telemetry := "out"
	if d.Get("telemetry").(bool) {
		telemetry = "in"
	}

	err = meta.(*Config).SetSetting(bootstrapSettingTelemetry, telemetry)
	if err != nil {
		return err
	}

	// Update admin user password if needed
	pass := d.Get("password").(string)
	changedPass, newPass, adminUser, err := meta.(*Config).SetUserPasswordByName(bootstrapDefaultUser, pass)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Admin password: %s", err)
	}

	if changedPass {
		d.Set("password", newPass)
		d.Set("current_password", newPass)
	}

	// Generate a new token if token_update is set or token is expired
	// Check if token is expired
	expiredToken, err := meta.(*Config).IsTokenExpired(d.Get("token_id").(string))
	if err != nil {
		return err
	}
	if d.Get("token_update").(bool) || expiredToken {
		tokenID, token, err := meta.(*Config).GenerateUserToken(bootstrapDefaultUser, bootstrapDefaultTokenDesc, d.Get("token_ttl").(int))
		if err != nil {
			return fmt.Errorf("[ERROR] Creating Admin token: %s", err)
		}

		// Delete old token
		err = meta.(*Config).DeleteToken(d.Get("token_id").(string))
		if err != nil {
			return fmt.Errorf("[ERROR] Deleting previous Admin token: %s", err)
		}

		// Update new tokenkey
		d.Set("token_id", tokenID)
		d.Set("token", token)
		err = meta.(*Config).UpdateToken(token)
		if err != nil {
			return fmt.Errorf("[ERROR] Updating Admin token: %s", err)
		}
	}

	// Set resource ID
	d.SetId(adminUser.ID)

	return resourceRancher2BootstrapRead(d, meta)
}

func resourceRancher2BootstrapDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}

func bootstrapDoLogin(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Doing login")

	// Try to connect with admin token
	token := d.Get("token").(string)
	err := meta.(*Config).UpdateToken(token)
	if err == nil {
		log.Printf("[INFO] Connecting with token")
		return nil
	}

	// If fails, try to connect with temp token
	token = d.Get("temp_token").(string)
	err = meta.(*Config).UpdateToken(token)
	if err == nil {
		log.Printf("[INFO] Connecting with temp token")
		return nil
	}

	// If fails, try to login with default admin user and current password
	currentPass := d.Get("current_password").(string)
	if len(currentPass) == 0 {
		currentPass = bootstrapDefaultPassword
	}
	tokenID, token, err := DoUserLogin(meta.(*Config).URL, bootstrapDefaultUser, currentPass, bootstrapDefaultTTL, bootstrapDefaultSessionDesc, meta.(*Config).CACerts, meta.(*Config).Insecure)
	if err != nil {
		return fmt.Errorf("[ERROR] Login with %s user: %v", bootstrapDefaultUser, err)
	}

	// Update config token
	err = meta.(*Config).UpdateToken(token)
	if err != nil {
		return fmt.Errorf("[ERROR] Connecting with user/pass: %s", err)
	}
	log.Printf("[INFO] Connecting with user/pass")

	// Delete temp token if exists
	err = meta.(*Config).DeleteToken(d.Get("temp_token_id").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Deleting temp token: %s", err)
	}

	// Update temp token data
	d.Set("temp_token_id", tokenID)
	d.Set("temp_token", token)

	return nil

}

func bootstrapCleanUpTempToken(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cleaning up temp token")

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	tokenID := d.Get("temp_token_id").(string)

	if len(tokenID) == 0 {
		// Clean up temp token data
		d.Set("temp_token_id", "")
		d.Set("temp_token", "")
		return nil
	}

	token, err := client.Token.ByID(tokenID)
	if err != nil {
		if IsNotFound(err) {
			// Clean up temp token data
			d.Set("temp_token_id", "")
			d.Set("temp_token", "")
			return nil
		}
		return err
	}

	// If token is current let temp token data
	if token.Current {
		return nil
	}

	// Delete temp token
	err = client.Token.Delete(token)
	if err != nil {
		return fmt.Errorf("[ERROR] Deleting temp token ID %s: %s", token.ID, err)
	}

	// Clean up temp token data
	d.Set("temp_token_id", "")
	d.Set("temp_token", "")

	return nil

}
