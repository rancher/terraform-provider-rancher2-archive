package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

const (
	testAccRancher2NamespaceMoveType    = "rancher2_namespace"
	testAccRancher2NamespaceMoveProject = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "local"
  description = "Terraform namespace acceptance test"
}
`

	testAccRancher2NamespaceMoveConfig = testAccRancher2NamespaceMoveProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = "${rancher2_project.foo.id}"
}
`

	testAccRancher2NamespaceMoveUpdateConfig = testAccRancher2NamespaceMoveProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test - updated"
  project_id = "${rancher2_project.foo.id}"
}
 `

	testAccRancher2NamespaceMoveRecreateConfig = testAccRancher2NamespaceMoveProject + `
resource "rancher2_namespace" "foo" {
  name = "foo"
  description = "Terraform namespace acceptance test"
  project_id = "${rancher2_project.foo.id}"
}
 `
)

func TestAccRancher2NamespaceMove_basic(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceMoveDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NamespaceMoveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceMoveExists(testAccRancher2NamespaceMoveType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "description", "Terraform namespace acceptance test"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceMoveUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceMoveExists(testAccRancher2NamespaceMoveType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "description", "Terraform namespace acceptance test - updated"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NamespaceMoveRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceMoveExists(testAccRancher2NamespaceMoveType+".foo", ns),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2NamespaceMoveType+".foo", "description", "Terraform namespace acceptance test"),
				),
			},
		},
	})
}

func TestAccRancher2NamespaceMove_disappears(t *testing.T) {
	var ns *clusterClient.Namespace

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NamespaceMoveDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NamespaceMoveConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NamespaceMoveExists(testAccRancher2NamespaceMoveType+".foo", ns),
					testAccRancher2NamespaceMoveDisappears(ns),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NamespaceMoveDisappears(ns *clusterClient.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NamespaceMoveType {
				continue
			}
			clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
			if err != nil {
				return err
			}
			client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
			if err != nil {
				return err
			}

			ns, err = client.Namespace.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.Namespace.Delete(ns)
			if err != nil {
				return fmt.Errorf("Error removing Namespace: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    namespaceStateRefreshFunc(client, ns.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for namespace (%s) to be removed: %s", ns.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2NamespaceMoveExists(n string, ns *clusterClient.Namespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No namespace ID is set")
		}

		clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}

		client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
		if err != nil {
			return err
		}

		foundNs, err := client.Namespace.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Namespace not found")
			}
			return err
		}

		ns = foundNs

		return nil
	}
}

func testAccCheckRancher2NamespaceMoveDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NamespaceMoveType {
			continue
		}

		clusterID, err := clusterIDFromProjectID(rs.Primary.Attributes["project_id"])
		if err != nil {
			return err
		}
		client, err := testAccProvider.Meta().(*Config).ClusterClient(clusterID)
		if err != nil {
			return err
		}

		_, err = client.Namespace.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Namespace still exists")
	}
	return nil
}
