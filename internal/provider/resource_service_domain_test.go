package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccServiceDomainResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServiceDomainResourceConfigDefault("terraform-tester"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", "terraform-tester"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "environment_id", "d0519b29-5d12-4857-a5dd-76fa7418336c"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "service_id", "39da7e07-fa3a-42fd-b695-d229319f2993"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "project_id", "0bb01547-570d-4109-a5e8-138691f6a2d1"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", "terraform-tester.up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_service_domain.test",
				ImportState:       true,
				ImportStateId:     "39da7e07-fa3a-42fd-b695-d229319f2993:staging:terraform-tester.up.railway.app",
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccServiceDomainResourceConfigDefault("terraform-tester"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", "terraform-tester"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "environment_id", "d0519b29-5d12-4857-a5dd-76fa7418336c"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "service_id", "39da7e07-fa3a-42fd-b695-d229319f2993"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "project_id", "0bb01547-570d-4109-a5e8-138691f6a2d1"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", "terraform-tester.up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// Update with default values
			{
				Config: testAccServiceDomainResourceConfigDefault("terraform-tester-2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", "terraform-tester-2"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "environment_id", "d0519b29-5d12-4857-a5dd-76fa7418336c"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "service_id", "39da7e07-fa3a-42fd-b695-d229319f2993"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "project_id", "0bb01547-570d-4109-a5e8-138691f6a2d1"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "domain", "terraform-tester-2.up.railway.app"),
					resource.TestCheckResourceAttr("railway_service_domain.test", "suffix", "up.railway.app"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_service_domain.test",
				ImportState:       true,
				ImportStateId:     "39da7e07-fa3a-42fd-b695-d229319f2993:staging:terraform-tester-2.up.railway.app",
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServiceDomainResourceExternalDeletion(t *testing.T) {
	projectName := "tf-acc-domain-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	subdomain := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	config := testAccServiceDomainResourceExternalDeletionConfig(projectName, subdomain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					testAccDeleteServiceDomain,
				),
			},
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_service_domain.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_service_domain.test", "subdomain", subdomain),
				),
			},
		},
	})
}

func testAccDeleteServiceDomain(state *terraform.State) error {
	domain, ok := state.RootModule().Resources["railway_service_domain.test"]
	if !ok {
		return fmt.Errorf("service domain resource is missing from state")
	}
	_, err := deleteServiceDomain(context.Background(), testAccClient(), domain.Primary.ID)
	return err
}

func testAccServiceDomainResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "railway_service_domain" "test" {
  subdomain = "%s"
  environment_id = "d0519b29-5d12-4857-a5dd-76fa7418336c"
  service_id = "39da7e07-fa3a-42fd-b695-d229319f2993"
}
`, name)
}

func testAccServiceDomainResourceExternalDeletionConfig(projectName string, subdomain string) string {
	return fmt.Sprintf(`
resource "railway_project" "test" {
  name = "%s"
}

resource "railway_service" "test" {
  name       = "domain-test"
  project_id = railway_project.test.id
}

resource "railway_service_domain" "test" {
  subdomain     = "%s"
  environment_id = railway_project.test.default_environment.id
  service_id    = railway_service.test.id
}
`, projectName, subdomain)
}
