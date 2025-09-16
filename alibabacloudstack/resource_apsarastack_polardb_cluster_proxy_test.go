package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbClusterProxy_basic(t *testing.T) {
	var account map[string]interface{}
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tfaccount%d", rand)
	var basicMap = map[string]string{
		"db_cluster_id": CHECKSET,
	}
	resourceId := "alibabacloudstack_polardb_cluster_proxy.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &account, serviceFunc, "DescribePolardbClusterProxy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbClusterProxyConfigDependence)
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
					"db_cluster_id":          "${alibabacloudstack_polardb_cluster_instance.instance.id}",
					"db_proxy_cluster_class": "${data.alibabacloudstack_polardb_cluster_proxy_types.types.proxy_classes.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_cluster_class": "${data.alibabacloudstack_polardb_cluster_proxy_types.types.proxy_classes.1.id}",
				}),
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

func resourcePolardbClusterProxyConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	variable "db_type" {
		default = "PostgreSQL"
	}

	variable "db_version" {
		default = "14"
	}

	data "alibabacloudstack_polardb_cluster_proxy_types" "types" {
		dbtype = "${var.db_type}"
		dbversion = "${var.db_version}"
	}
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
	`, name, VSwitchCommonTestCase)
}
