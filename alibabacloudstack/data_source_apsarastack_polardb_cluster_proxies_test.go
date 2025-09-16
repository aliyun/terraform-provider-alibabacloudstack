package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterProxiesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_cluster_proxies.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%d", rand),
		dataSourcePolardbClusterProxiesDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_proxy.default.db_cluster_id}",
		}),
	}
	var existPolardbClusterProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"db_proxy_cluster_id":              CHECKSET,
			"proxy_instances.#":                "2",
			"proxy_instances.0.db_node_status": CHECKSET,
			"proxy_instances.0.db_node_id":     CHECKSET,
			"proxy_instances.0.db_node_class":  CHECKSET,
		}
	}

	var fakePolardbClusterProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"proxy_instances.#": "0",
		}
	}

	var PolardbClusterProxiesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbClusterProxiesMapFunc,
		fakeMapFunc:  fakePolardbClusterProxiesMapFunc,
	}

	PolardbClusterProxiesCheckInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourcePolardbClusterProxiesDependence(name string) string {
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
		db_type = "${var.db_type}"
		db_version = "${var.db_version}"
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
		storage_type			= "ESSDPL1"
		storage_space 			= 20
		db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
		zone_id					= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
		vswitch_id 				= "${alibabacloudstack_vpc_vswitch.default.id}"
		sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	}
	
	resource "alibabacloudstack_polardb_cluster_proxy" "default" {
		db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
		db_proxy_cluster_class = "${data.alibabacloudstack_polardb_cluster_proxy_types.types.proxy_classes.0.id}"
	}
	`, name, VSwitchCommonTestCase)
}
