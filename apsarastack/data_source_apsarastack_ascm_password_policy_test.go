package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_PasswordPolicies_DataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_PasswordPolicies,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_password_policies.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.hard_expiry"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.require_numbers"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.require_numbers"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.require_symbols"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.require_lowercase_characters"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.require_uppercase_characters"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.max_login_attempts"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.max_password_age"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.minimum_password_length"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_password_policies.default", "policies.password_reuse_prevention"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_PasswordPolicies = `

data "apsarastack_ascm_password_policies" "default" {

}
`
