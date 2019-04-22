package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2AuthConfigFreeIpaType   = "rancher2_auth_config_freeipa"
	testAccRancher2AuthConfigFreeIpaConfig = `
resource "rancher2_auth_config_freeipa" "freeipa" {
  servers = ["freeipa.test.local"]
  service_account_distinguished_name = "uid=admin,dc=test,dc=local"
  service_account_password = "XXXXXXXX"
  user_search_base = "dc=test,dc=local"
  port = 389
  group_dn_attribute = "entrydn"
  group_member_mapping_attribute = "member"
  group_member_user_attribute = "entrydn"
  group_object_class = "groupOfNames"
  user_name_attribute = "givenName"
}
`

	testAccRancher2AuthConfigFreeIpaUpdateConfig = `
resource "rancher2_auth_config_freeipa" "freeipa" {
  servers = ["freeipa.test.local"]
  service_account_distinguished_name = "uid=admin,cn=users,dc=test,dc=local"
  service_account_password = "YYYYYYYY"
  user_search_base = "cn=users,dc=test,dc=local"
  port = 389
  group_dn_attribute = "entrydn"
  group_member_mapping_attribute = "member"
  group_member_user_attribute = "entrydn"
  group_object_class = "groupOfNames"
  user_name_attribute = "givenName-updated"
}
 `

	testAccRancher2AuthConfigFreeIpaRecreateConfig = `
resource "rancher2_auth_config_freeipa" "freeipa" {
  servers = ["freeipa.test.local"]
  service_account_distinguished_name = "uid=admin,dc=test,dc=local"
  service_account_password = "XXXXXXXX"
  user_search_base = "dc=test,dc=local"
  port = 389
  group_dn_attribute = "entrydn"
  group_member_mapping_attribute = "member"
  group_member_user_attribute = "entrydn"
  group_object_class = "groupOfNames"
  user_name_attribute = "givenName"
}
 `
)

func TestAccRancher2AuthConfigFreeIpa_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigFreeIpaDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigFreeIpaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "name", AuthConfigFreeIpaName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_name_attribute", "givenName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_distinguished_name", "uid=admin,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_password", "XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigFreeIpaUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "name", AuthConfigFreeIpaName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_name_attribute", "givenName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_distinguished_name", "uid=admin,cn=users,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_search_base", "cn=users,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_password", "YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigFreeIpaRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "name", AuthConfigFreeIpaName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_name_attribute", "givenName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_distinguished_name", "uid=admin,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, "service_account_password", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigFreeIpa_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigFreeIpaDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigFreeIpaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigFreeIpaType+"."+AuthConfigFreeIpaName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigFreeIpaType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigFreeIpaDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigFreeIpaType {
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
