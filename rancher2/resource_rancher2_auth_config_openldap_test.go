package rancher2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2AuthConfigOpenLdapType   = "rancher2_auth_config_openldap"
	testAccRancher2AuthConfigOpenLdapConfig = `
resource "rancher2_auth_config_openldap" "openldap" {
  servers = ["openldap.test.local"]
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

	testAccRancher2AuthConfigOpenLdapUpdateConfig = `
resource "rancher2_auth_config_openldap" "openldap" {
  servers = ["openldap.test.local"]
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

	testAccRancher2AuthConfigOpenLdapRecreateConfig = `
resource "rancher2_auth_config_openldap" "openldap" {
  servers = ["openldap.test.local"]
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

func TestAccRancher2AuthConfigOpenLdap_basic(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigOpenLdapDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigOpenLdapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "name", AuthConfigOpenLdapName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_name_attribute", "givenName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_distinguished_name", "uid=admin,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_password", "XXXXXXXX"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigOpenLdapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "name", AuthConfigOpenLdapName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_name_attribute", "givenName-updated"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_distinguished_name", "uid=admin,cn=users,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_search_base", "cn=users,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_password", "YYYYYYYY"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2AuthConfigOpenLdapRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, authConfig),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "name", AuthConfigOpenLdapName),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_name_attribute", "givenName"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_distinguished_name", "uid=admin,dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "user_search_base", "dc=test,dc=local"),
					resource.TestCheckResourceAttr(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, "service_account_password", "XXXXXXXX"),
				),
			},
		},
	})
}

func TestAccRancher2AuthConfigOpenLdap_disappears(t *testing.T) {
	var authConfig *managementClient.AuthConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2AuthConfigOpenLdapDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2AuthConfigOpenLdapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2AuthConfigExists(testAccRancher2AuthConfigOpenLdapType+"."+AuthConfigOpenLdapName, authConfig),
					testAccRancher2AuthConfigDisappears(authConfig, testAccRancher2AuthConfigOpenLdapType),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckRancher2AuthConfigOpenLdapDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2AuthConfigOpenLdapType {
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
