package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbAccountUpdate(t *testing.T) {
	var account *PolardbDescribeaccountsResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"data_base_instance_id": CHECKSET,
		"account_name":          "tftestnormal",
		"account_password":      "inputYourCodeHere",
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

		Providers: testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"data_base_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
					"account_name":          "tftestnormal",
					"account_password":      "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
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
					"account_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_description": "tf test",
					"account_password":    "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_description": "tf test",
						"account_password":    "inputYourCodeHere",
					}),
				),
			},
		},
	})
}

func resourcePolardbAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	resource "alibabacloudstack_polardb_dbinstance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	`, RdsCommonTestCase, name)
}
