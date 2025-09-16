package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbClusterAccount_basic(t *testing.T) {
	var account map[string]interface{}
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbaccount-%d", rand)
	var basicMap = map[string]string{
		"data_base_instance_id": CHECKSET,
		"account_name":          "tftestnormal",
		"account_password":      "${random_password.password.0.result}",
		"account_type":          "Normal",
	}
	resourceId := "alibabacloudstack_polardb_cluster_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &account, serviceFunc, "DescribePolardbClusterAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterAccountConfigDependence)
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
					"db_cluster_id":    "${alibabacloudstack_polardb_cluster_instance.instance.id}",
					"account_name":     "tftestnormal",
					"account_password": "${random_password.password.0.result}",
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
					testAccCheck(map[string]string{
						"account_password": "${random_password.password.1.result}",
					}),
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
						"account_password":    "${random_password.password.2.result}",
					}),
				),
			},
		},
	})
}

func resourcePolardbClusterAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	%s
	%s
	resource "alibabacloudstack_polardb_cluster_instance" "instance" {
		db_cluster_description 	= "${var.name}"
		db_type            		= "MySQL"
		db_version    			= "5.7"
		instance_name 			= "${var.name}"
		storage_type			= "ESSDPL1"
		storage_space 			= 20
		db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
		db_node_num 			= "1"
		zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
		vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
		sub_category 			= "General"
	}
	`, name, RandomPasswordTestCase(12, 3), VSwitchCommonTestCase)
}
