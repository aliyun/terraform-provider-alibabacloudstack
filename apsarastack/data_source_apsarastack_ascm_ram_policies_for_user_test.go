package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmRamPoliciesForUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceapsarastackascmRamPoliciesForUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_ram_policies_for_user.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.policy_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.policy_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.attach_date"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.default_version"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_policies_for_user.default", "policies.policy_document"),
				),
			},
		},
	})
}

const datasourceapsarastackascmRamPoliciesForUser = `
data "apsarastack_ascm_ram_policies_for_user" "default" {
  login_name = "test_admin"
}
`
