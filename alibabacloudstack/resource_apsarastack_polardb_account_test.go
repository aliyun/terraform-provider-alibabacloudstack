package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbAccount_Update(t *testing.T) {
	var account *PolardbDescribeaccountsResponse
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_account%d", rand)
	var basicMap = map[string]string{
		"data_base_instance_id": CHECKSET,
		"account_name":          name,
		"account_type":          "Normal",
	}
	resourceId := "alibabacloudstack_polardb_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &account, serviceFunc, "DescribeDBAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbAccountConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_base_instance_id": "${local.polardb_dbinstance_id}",
					"account_name":          "${var.name}",
					"account_password":      "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password is a sensitive field, it will not be displayed after setting
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "tf test",
					"account_password":    "${random_password.password.2.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "tf test",
					}),
				),
			},
		},
	})
}

func resourcePolardbAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	%s
	variable "creation" {
		default = "PolarDB"
	}
	%s
	
	%s
	
	`, name, VSwitchCommonTestCase, PolarDBCommonTestCase("MySQL", true), RandomPasswordTestCase(12, 3))
}
