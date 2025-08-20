package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxSuperAccount_basic0(t *testing.T) {
	var v []PolardbxAccount

	resourceId := "alibabacloudstack_polardbx_super_account.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"instance_id": CHECKSET,
	})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeSuperAccountRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_pldbx_account_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxSuperAccountDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":                  "${local.polardbx_instance.id}",
					"admin_account_name":           "admin_user",
					"admin_account_password":       "${random_password.password.0.result}",
					"admin_account_description":    "system user",
					"security_account_name":        "security_user",
					"security_account_description": "security user",
					"security_account_password":    "${random_password.password.0.result}",
					"audit_account_name":           "audit_user",
					"audit_account_description":    "audit user",
					"audit_account_password":       "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"three_roles": "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password无法回读
				ImportStateVerifyIgnore: []string{"admin_account_password", "security_account_password", "audit_account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_account_password":    "${random_password.password.1.result}",
					"security_account_password": "${random_password.password.1.result}",
					"audit_account_password":    "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"admin_account_description":    "system user update",
					"security_account_description": "security user update",
					"audit_account_description":    "audit user update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_account_name": "security_user1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_account_name":        REMOVEKEY,
					"security_account_description": REMOVEKEY,
					"security_account_password":    REMOVEKEY,
					"audit_account_name":           REMOVEKEY,
					"audit_account_description":    REMOVEKEY,
					"audit_account_password":       REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_account_name":        REMOVEKEY,
						"security_account_description": REMOVEKEY,
						"audit_account_name":           REMOVEKEY,
						"audit_account_description":    REMOVEKEY,
						"three_roles":                  "false",
					}),
				),
			},
		},
	})
}

func resourcePolardbxSuperAccountDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

%s

 `, name, RandomPasswordTestCase(12, 2), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
