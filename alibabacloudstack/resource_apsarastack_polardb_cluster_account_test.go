package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbClusterAccount_basic(t *testing.T) {
	var account map[string]interface{}
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tfaccount%d", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
		"account_name":  name,
		"account_type":  "Normal",
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
		ExternalProviders: testAccExternalProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":    "${alibabacloudstack_polardb_cluster_instance.instance.id}",
					"account_name":     "${var.name}",
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
					"account_password":    "${random_password.password.1.result}",
					"account_lock_state":  "Lock",
					"account_description": "tf test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_lock_state":  "Lock",
						"account_description": "tf test",
					}),
				),
			},
		},
	})
}

func resourcePolardbClusterAccountConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	variable "db_type" {
	  default = "PostgreSQL"
	}

	variable "db_version" {
	  default = "14"
	}

	%s

	data "alibabacloudstack_polardb_cluster_instance_types" "default" {
	  db_type = "${var.db_type}"
	  db_version = "${var.db_version}"
	  sorted_by = "CPU"
	  sub_category = "normal_exclusive"
	}

	%s
	
	resource "alibabacloudstack_polardb_cluster_instance" "instance" {
		db_cluster_description 	= "${var.name}"
		db_type            		= "${var.db_type}"
		db_version    			= "${var.db_version}"
		storage_type			= "ESSDPL1"
		storage_space 			= 20
		db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
		zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
		sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	}
	`, name, RandomPasswordTestCase(12, 2), VSwitchCommonTestCase)
}
