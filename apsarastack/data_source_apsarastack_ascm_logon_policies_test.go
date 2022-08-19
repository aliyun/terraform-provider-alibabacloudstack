package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmLogonPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceapsarastackascmLogonPolicies,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_logon_policies.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_logon_policies.default", "policies.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_logon_policies.default", "policies.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_logon_policies.default", "policies.login_policy_id"),
				),
			},
		},
	})
}

const datasourceapsarastackascmLogonPolicies = `
resource "apsarastack_ascm_logon_policy" "default" {
  name="test_foo1"
  description="testing purpose"
  rule="ALLOW"
}
data "apsarastack_ascm_logon_policies" "default"{
	name_regex = apsarastack_ascm_logon_policy.default.name
}
`
