package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscmLogonPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourcealibabacloudstackascmLogonPolicies,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_logon_policies.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_logon_policies.default", "policies.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_logon_policies.default", "policies.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_logon_policies.default", "policies.login_policy_id"),
				),
			},
		},
	})
}

const datasourcealibabacloudstackascmLogonPolicies = `
resource "alibabacloudstack_ascm_logon_policy" "default" {
  name="test_foo1"
  description="testing purpose"
  rule="ALLOW"
}
data "alibabacloudstack_ascm_logon_policies" "default"{
	name_regex = alibabacloudstack_ascm_logon_policy.default.name
}
`
