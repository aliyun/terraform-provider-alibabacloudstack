package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_polardb_cluster_databases.default"
	name := fmt.Sprintf("tfaccount%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, datasourcePolardbClusterDatabasesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_database.default.id}"},
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_polardb_cluster_database.default.id}_fake"},
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	db_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_name":       "${alibabacloudstack_polardb_cluster_database.default.db_name}",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_name":       "${alibabacloudstack_polardb_cluster_database.default.db_name}_fake",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	name_regex_Conf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alibabacloudstack_polardb_cluster_database.default.db_name}",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":    "${alibabacloudstack_polardb_cluster_database.default.db_name}_fake",
			"db_cluster_id": "${alibabacloudstack_polardb_cluster_instance.instance.id}",
		}),
	}

	AlibabacloudstackPolardbClusterDatabasesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, db_nameConf, name_regex_Conf)
}

var existAlibabacloudstackPolardbClusterDatabasesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#":                    "1",
		"databases.0.id":                 CHECKSET,
		"databases.0.db_description":     CHECKSET,
		"databases.0.db_name":            fmt.Sprintf("tfaccount%d", rand),
		"databases.0.character_set_name": CHECKSET,
	}
}

var fakeAlibabacloudstackPolardbClusterDatabasesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#": "0",
	}
}

var AlibabacloudstackPolardbClusterDatabasesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_polardb_cluster_databases.default",
	existMapFunc: existAlibabacloudstackPolardbClusterDatabasesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackPolardbClusterDatabasesDataMapFunc,
}

func datasourcePolardbClusterDatabasesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	variable "creation" {
		default = "PolarDB"
	}
	variable "db_type" {
		default = "MySQL"
	}

	variable "db_version" {
		default = "5.7"
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
	
	resource "alibabacloudstack_polardb_cluster_database" "default" {
		db_cluster_id		= "${alibabacloudstack_polardb_cluster_instance.instance.id}"
		db_name 			= "${var.name}"
		character_set_name 	= "utf8"
	}
	`, name, VSwitchCommonTestCase)
}
