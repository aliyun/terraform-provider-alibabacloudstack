package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscm_PasswordPolicy(t *testing.T) {
	var v *PasswordPolicy
	resourceId := "alibabacloudstack_ascm_password_policy.default"
	ra := resourceAttrInit(resourceId, ascmPassword)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudstackPasswordPolicyMinimumPasswordLength,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

const testAccAlibabacloudstackPasswordPolicyMinimumPasswordLength = `
resource "alibabacloudstack_ascm_password_policy" "default" {
		minimum_password_length = 10
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
	"minimum_password_length":      "10",
	"require_lowercase_characters": "true",
	"require_uppercase_characters": "true",
	"require_numbers":              "true",
	"require_symbols":              "true",
	"hard_expiry":                  "true",
	"max_password_age":             "90",
	"password_reuse_prevention":    "5",
	"max_login_attempts":           "5",
}
