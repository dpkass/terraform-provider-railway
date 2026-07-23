package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBucketResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketResourceConfig("terraform-provider-bucket-test", "ams"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("railway_bucket.test", "name", "terraform-provider-bucket-test"),
					resource.TestCheckResourceAttr("railway_bucket.test", "region", "ams"),
				),
			},
			{
				ResourceName:      "railway_bucket.test",
				ImportState:       true,
				ImportStateIdFunc: bucketImportIDFunc,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketResourceConfig("terraform-provider-bucket-renamed", "ams"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("railway_bucket.test", "name", "terraform-provider-bucket-renamed"),
				),
			},
			{
				Config: testAccBucketResourceConfig("terraform-provider-bucket-renamed", "iad"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("railway_bucket.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttr("railway_bucket.test", "region", "iad"),
			},
		},
	})
}

func testAccBucketResourceConfig(name string, region string) string {
	return fmt.Sprintf(`
resource "railway_bucket" "test" {
  name           = "%s"
  project_id     = "0bb01547-570d-4109-a5e8-138691f6a2d1"
  environment_id = "d0519b29-5d12-4857-a5dd-76fa7418336c"
  region         = "%s"
}
`, name, region)
}

func bucketImportIDFunc(state *terraform.State) (string, error) {
	rawState, ok := state.RootModule().Resources["railway_bucket.test"]
	if !ok {
		return "", fmt.Errorf("resource not found")
	}

	return fmt.Sprintf(
		"%s:%s:%s",
		rawState.Primary.Attributes["project_id"],
		rawState.Primary.Attributes["environment_id"],
		rawState.Primary.Attributes["id"],
	), nil
}
