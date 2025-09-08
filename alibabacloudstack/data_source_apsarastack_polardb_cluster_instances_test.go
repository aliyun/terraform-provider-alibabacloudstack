package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbDbClusterInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_cluster_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc-%d", rand),
		dataSourcePolardbDbClusterInstancesDependence)

	db_cluster_descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardb_cluster_instance.default.db_cluster_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardb_cluster_instance.default.db_cluster_description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_cluster_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_cluster_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardb_cluster_instance.default.db_cluster_description}",
			"ids":        []string{"${alibabacloudstack_polardb_cluster_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardb_cluster_instance.default.db_cluster_description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_polardb_cluster_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existPolardbDbClusterInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"ids.0":                  CHECKSET,
			"db_cluster_instances.#": "1",
			"db_cluster_instances.0.db_cluster_description": fmt.Sprintf("tf-testAcc-%d", rand),
			"db_cluster_instances.0.create_time":            CHECKSET,
			"db_cluster_instances.0.db_cluster_id":          CHECKSET,
			"db_cluster_instances.0.db_nodes.#":             "2",
			"db_cluster_instances.0.engine":                 CHECKSET,
			"db_cluster_instances.0.architecture":           CHECKSET,
			"db_cluster_instances.0.db_cluster_status":      CHECKSET,
			"db_cluster_instances.0.ascm_create_user":       CHECKSET,
		}
	}

	var fakePolardbDbClusterInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"db_cluster_instances.#": "0",
		}
	}

	var PolardbDbClusterInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbDbClusterInstancesMapFunc,
		fakeMapFunc:  fakePolardbDbClusterInstancesMapFunc,
	}

	PolardbDbClusterInstancesCheckInfo.dataSourceTestCheck(t, rand, db_cluster_descriptionRegexConf, idsConf, allConf)
}

func dataSourcePolardbDbClusterInstancesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "8.0"
}

%s

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by = "CPU"
  cpu_type = "hygon"
  sub_category = "General"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
	db_cluster_description 	=  "${var.name}"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id			  	= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	db_node_num 			= "2"
	sub_category 			= "General"
	storage_type 			= "ESSDPL1"
	tde_enabled 			= "true"
}

 `, name, VSwitchCommonTestCase)
}
