package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEnvironmentResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEnvironmentResourceConfigDefault("integration"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_environment.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_environment.test", "name", "integration"),
					resource.TestCheckResourceAttr("railway_environment.test", "project_id", "0bb01547-570d-4109-a5e8-138691f6a2d1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "railway_environment.test",
				ImportState:       true,
				ImportStateId:     "0bb01547-570d-4109-a5e8-138691f6a2d1:integration",
				ImportStateVerify: true,
			},
			// Update with default values
			{
				Config: testAccEnvironmentResourceConfigDefault("integration"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_environment.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_environment.test", "name", "integration"),
					resource.TestCheckResourceAttr("railway_environment.test", "project_id", "0bb01547-570d-4109-a5e8-138691f6a2d1"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccEnvironmentResourceClone(t *testing.T) {
	projectName := "tf-acc-env-clone-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentResourceConfigClone(projectName, "integration-clone"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("railway_project.test_clone", "id", uuidRegex()),
					resource.TestMatchResourceAttr("railway_environment.test", "id", uuidRegex()),
					resource.TestCheckResourceAttr("railway_environment.test", "name", "integration-clone"),
					resource.TestCheckResourceAttrPair("railway_environment.test", "project_id", "railway_project.test_clone", "id"),
					resource.TestCheckResourceAttrPair("railway_environment.test", "source_environment_id", "railway_project.test_clone", "default_environment.id"),
					resource.TestCheckResourceAttr("railway_environment.test", "skip_initial_deploys", "true"),
				),
			},
			{
				ResourceName:            "railway_environment.test",
				ImportState:             true,
				ImportStateIdFunc:       testAccEnvironmentCloneImportStateId,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_environment_id", "skip_initial_deploys"},
			},
		},
	})
}

func testAccEnvironmentResourceConfigDefault(name string) string {
	return fmt.Sprintf(`
resource "railway_environment" "test" {
  name = "%s"
  project_id = "0bb01547-570d-4109-a5e8-138691f6a2d1"
}
`, name)
}

func testAccEnvironmentResourceConfigClone(projectName string, name string) string {
	return fmt.Sprintf(`
resource "railway_project" "test_clone" {
  name = "%s"
}

resource "railway_environment" "test" {
  name                  = "%s"
  project_id            = railway_project.test_clone.id
  source_environment_id = railway_project.test_clone.default_environment.id
  skip_initial_deploys  = true
}
`, projectName, name)
}

func testAccEnvironmentCloneImportStateId(state *terraform.State) (string, error) {
	resourceState := state.RootModule().Resources["railway_environment.test"]
	return fmt.Sprintf(
		"%s:%s",
		resourceState.Primary.Attributes["project_id"],
		resourceState.Primary.Attributes["name"],
	), nil
}
