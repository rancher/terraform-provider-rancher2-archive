package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2AuthConfigAzureADType   = "rancher2_auth_config_azuread"
	testAccRancher2AuthConfigAzureADConfig = `
resource "rancher2_auth_config_azuread" "azuread" {
  application_id = "XXXXXX"
  application_secret = "XXXXXXXX"
  auth_endpoint = "authorize"
  graph_endpoint = "graph"
  rancher_url = "https://RANCHER"
  tenant_id = "XXXXXXXX"
  token_endpoint = "token"
}
`

	testAccRancher2AuthConfigAzureADUpdateConfig = `
resource "rancher2_auth_config_azuread" "azuread" {
  application_id = "XXXXXX"
  application_secret = "YYYYYYYY"
  auth_endpoint = "authorize-updated"
  graph_endpoint = "graph"
  rancher_url = "https://RANCHER-UPDATED"
  tenant_id = "YYYYYYYY"
  token_endpoint = "token"
}
 `

	testAccRancher2AuthConfigAzureADRecreateConfig = `
resource "rancher2_auth_config_azuread" "azuread" {
  application_id = "XXXXXX"
  application_secret = "XXXXXXXX"
  auth_endpoint = "authorize"
  graph_endpoint = "graph"
  rancher_url = "https://RANCHER"
  tenant_id = "XXXXXXXX"
  token_endpoint = "token"
}
 `
)

func TestAccRancher2AuthConfigAzureAD_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigAzureADDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigAzureADConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "name", AzureADConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "auth_endpoint", "authorize"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "rancher_url", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "application_secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "tenant_id", "XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigAzureADUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "name", AzureADConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "auth_endpoint", "authorize-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "rancher_url", "https://RANCHER-UPDATED"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "application_secret", "YYYYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "tenant_id", "YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigAzureADRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "name", AzureADConfigName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "auth_endpoint", "authorize"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "rancher_url", "https://RANCHER"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "application_secret", "XXXXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, "tenant_id", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigAzureAD_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigAzureADDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigAzureADConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigAzureADType+"."+AzureADConfigName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigAzureADType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigAzureADDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigAzureADType {
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
