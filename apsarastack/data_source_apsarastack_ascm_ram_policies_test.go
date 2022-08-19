package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmRamPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceapsarastackascmRamPolicies,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_ram_policies.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.ctime"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.policy_document"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.region"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies.default", "policies.cuser_id"),
				),
			},
		},
	})
}

const datasourceapsarastackascmRamPolicies = `
resource "apsarastack_ascm_ram_policy" "default" {
  name = "TestingRamPolicy"
  description = "Testing Policy"
  policy_document = "{\"Statement\":[{\"Action\":\"ecs:*\",\"Effect\":\"Allow\",\"Resource\":\"*\"}],\"Version\":\"1\"}"
}

data "apsarastack_ascm_ram_policies" "default" {
  name_regex = apsarastack_ascm_ram_policy.default.name
}
`
