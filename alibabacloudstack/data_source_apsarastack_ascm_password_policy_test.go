package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_PasswordPolicies_DataSource(t *testing.T) { //not completed
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_PasswordPolicies,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_password_policies.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.hard_expiry"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.require_numbers"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.require_numbers"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.require_symbols"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.require_lowercase_characters"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.require_uppercase_characters"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.max_login_attempts"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.max_password_age"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.minimum_password_length"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_password_policies.default", "policies.password_reuse_prevention"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_PasswordPolicies = `

data "alibabacloudstack_ascm_password_policies" "default" {

}
`
