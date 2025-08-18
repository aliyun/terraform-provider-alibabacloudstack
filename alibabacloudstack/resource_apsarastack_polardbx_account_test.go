package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxAccount_basic0(t *testing.T) {
	var v *PolardbxAccount

	resourceId := "alibabacloudstack_polardbx_account.default"
	ra := resourceAttrInit(resourceId, PolardbxAccountbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribeAccountRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_pldbx_account_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxAccountDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${local.polardbx_instance.id}",
					"account_name": "${var.name}",
					"password":     "${random_password.password.0.result}",
					"description":  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name": name,
						"description":  name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password无法回读
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("%s_update", name),
					}),
				),
			},
		},
	})
}

func resourcePolardbxAccountDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

%s

 `, name, RandomPasswordTestCase(12,2), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}

var PolardbxAccountbasicMap = map[string]string{
	"instance_id":  CHECKSET,
	"account_name": CHECKSET,
}
