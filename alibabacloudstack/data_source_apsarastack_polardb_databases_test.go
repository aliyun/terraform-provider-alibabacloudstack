package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_databases.default"
	name := fmt.Sprintf("tf-testaccpolardbdatabases%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbDatabaseDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}"},
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}_fake"},
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
		}),
	}

	data_base_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}"},
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
			"data_base_name":        "${alibabacloudstack_polardb_database.default.data_base_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alibabacloudstack_polardb_database.default.id}_fake"},
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
			"data_base_name":        "${alibabacloudstack_polardb_database.default.data_base_name}_fake",
		}),
	}

	name_regex_Conf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":            "${alibabacloudstack_polardb_database.default.data_base_name}",
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":            "${alibabacloudstack_polardb_database.default.data_base_name}-fakeTestAcccc",
			"data_base_instance_id": "${local.polardb_dbinstance_id}",
		}),
	}

	AlibabacloudstackPolardbDatabasesCheckInfo.dataSourceTestCheck(t, rand, idsConf, data_base_nameConf, name_regex_Conf)
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

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

%s

resource "alibabacloudstack_polardb_database" "default" {
	data_base_instance_id = "${local.polardb_dbinstance_id}"
	data_base_description = "Automatically generated test"
	data_base_name        = var.name
	character_set_name = "utf8"
}
`, name, PolarDBCommonTestCase("MySQL",false))
}
