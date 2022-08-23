package apsarastack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_PasswordPolicy(t *testing.T) {
	var v *PasswordPolicy
	resourceId := "apsarastack_ascm_password_policy.default"
	ra := resourceAttrInit(resourceId, ascmPassword)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApsarastackPasswordPolicyMinimumPasswordLength,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

const testAccApsarastackPasswordPolicyMinimumPasswordLength = `
resource "apsarastack_ascm_password_policy" "default" {
		minimum_password_length = 12
	    require_lowercase_characters = true
	   require_uppercase_characters = true
       require_numbers = true
	    require_symbols = true
	    hard_expiry = true
	    max_password_age = 90
	    password_reuse_prevention = 5
	    max_login_attempts = 5
}
`

var ascmPassword = map[string]string{
	"minimum_password_length":      "12",
	"require_lowercase_characters": "true",
	"require_uppercase_characters": "true",
	"require_numbers":              "true",
	"require_symbols":              "true",
	"hard_expiry":                  "true",
	"max_password_age":             "90",
	"password_reuse_prevention":    "5",
	"max_login_attempts":           "5",
}
