package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2AuthConfigPingType   = "rancher2_auth_config_ping"
	testAccRancher2AuthConfigPingConfig = `
resource "rancher2_auth_config_ping" "ping" {
  display_name_field = "displayName"
  groups_field = "memberOf"
  uid_field = "distinguishedName"
  user_name_field = "sAMAccountName"
  idp_metadata_content = "XXXXXXXX"
  rancher_api_host = "https://RANCHER"
  sp_cert = "XXXXXX"
  sp_key = "XXXXXXXX"
}
`

	testAccRancher2AuthConfigPingUpdateConfig = `
resource "rancher2_auth_config_ping" "ping" {
  display_name_field = "displayName"
  groups_field = "memberOf"
  uid_field = "distinguishedName"
  user_name_field = "sAMAccountName-updated"
  idp_metadata_content = "YYYYYYYY"
  rancher_api_host = "https://RANCHER-UPDATED"
  sp_cert = "XXXXXX"
  sp_key = "YYYYYYYY"
}
 `

	testAccRancher2AuthConfigPingRecreateConfig = `
resource "rancher2_auth_config_ping" "ping" {
  display_name_field = "displayName"
  groups_field = "memberOf"
  uid_field = "distinguishedName"
  user_name_field = "sAMAccountName"
  idp_metadata_content = "XXXXXXXX"
  rancher_api_host = "https://RANCHER"
  sp_cert = "XXXXXX"
  sp_key = "XXXXXXXX"
}
 `
)

func TestAccRancher2AuthConfigPing_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigPingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigPingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+PingConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "name", PingConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigPingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+PingConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "name", PingConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "user_name_field", "sAMAccountName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "rancher_api_host", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "sp_key", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "idp_metadata_content", "YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigPingRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+PingConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "name", PingConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "user_name_field", "sAMAccountName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "rancher_api_host", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "sp_key", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigPingType+"."+PingConfigName, "idp_metadata_content", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigPing_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigPingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigPingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigPingType+"."+PingConfigName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigPingType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigPingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigPingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		auth, err := client.AuthConfig.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}

		if auth.Enabled == true {
			err = client.Post(auth.Actions["disable"], nil, nil)
			if err != nil {
				return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", auth.ID, err)
			}
		}
		return nil
	}
	return nil
}
