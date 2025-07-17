package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_databases.default"
	name := fmt.Sprintf("tf-testAccPolardbBackups%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbDatabaseDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_database.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_database.default.id}_fake"},
		}),
	}

	data_base_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}"},
			"data_base_instance_id": "${alibabacloudstack_polardb_instance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}_fake"},
			"data_base_instance_id": "${alibabacloudstack_polardb_instance.instance.id}_fake",
		}),
	}

	data_base_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_database.default.id}"},
			"data_base_name": "${alibabacloudstack_polardb_database.default.data_base_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_database.default.id}_fake"},
			"data_base_name": "${alibabacloudstack_polardb_database.default.data_base_name}_fake",
		}),
	}

	AlibabacloudstackPolardbDatabasesCheckInfo.dataSourceTestCheck(t, rand, idsConf, data_base_instance_idConf, data_base_nameConf)
}

var existAlibabacloudstackPolardbDatabasesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#":    "1",
		"databases.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackPolardbDatabasesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#": "0",
	}
}

var AlibabacloudstackPolardbDatabasesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_polardb_databases.default",
	existMapFunc: existAlibabacloudstackPolardbDatabasesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackPolardbDatabasesMapFunc,
}

func dataSourcePolardbDatabaseDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_polardb_instance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}
resource "alibabacloudstack_polardb_database" "default" {
	data_base_instance_id = "${alibabacloudstack_polardb_instance.instance.id}"
	data_base_description = "自动化生成测试"
	data_base_name        = "tftest"
	character_set_name = "utf8"
}
`, name)
}
