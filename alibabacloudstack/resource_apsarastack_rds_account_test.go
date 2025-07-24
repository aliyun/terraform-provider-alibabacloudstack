package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBAccountUpdate(t *testing.T) {
	var v *rds.DBInstanceAccount
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"instance_id": CHECKSET,
		"name":        "tftestnormal",
		"password":    CHECKSET,
		"type":        "Normal",
	}
	resourceId := "alibabacloudstack_db_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		ExternalProviders: testAccExternalProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alibabacloudstack_db_instance.default.id}",
					"name":        "tftestnormal",
					"password": "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				// password敏感字段设置后不回显
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf test",
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf test",
					}),
				),
			},
		},
	})
}

func resourceDBAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	
	variable "creation" {
		default = "Rds"
	}
	
	variable "name" {
		default = "%v"
	}
	
	resource "random_password" "password" {
		count            = 2
		length           = 12
		special          = true
		override_special = "!@#$^&*()_"
		min_lower        = 1
		min_upper        = 1
		min_numeric      = 1
	}
	
	%s
	%s

	`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
