package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2AuthConfigGithubType   = "rancher2_auth_config_github"
	testAccRancher2AuthConfigGithubConfig = `
resource "rancher2_auth_config_github" "github" {
  client_id = "XXXXXX"
  client_secret = "XXXXXXXX"
}
`

	testAccRancher2AuthConfigGithubUpdateConfig = `
resource "rancher2_auth_config_github" "github" {
  client_id = "YYYYYY"
  client_secret = "YYYYYYYY"
}
 `

	testAccRancher2AuthConfigGithubRecreateConfig = `
resource "rancher2_auth_config_github" "github" {
  client_id = "XXXXXX"
  client_secret = "XXXXXXXX"
}
 `
)

func TestAccRancher2AuthConfigGithub_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigGithubDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigGithubConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "name", AuthConfigGithubName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_id", "XXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_secret", "XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigGithubUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "name", AuthConfigGithubName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_id", "YYYYYY"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_secret", "YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigGithubRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "name", AuthConfigGithubName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_id", "XXXXXX"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, "client_secret", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigGithub_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigGithubDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigGithubConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigGithubType+"."+AuthConfigGithubName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigGithubType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigGithubDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigGithubType {
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
