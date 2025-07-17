package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPolardbDatabasesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_databases.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_databases.default.id}_fake"]`,
		}),
	}

	data_base_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_polardb_databases.default.id}"]`,
			"data_base_instance_id": `"${alibabacloudstack_polardb_databases.default.DataBaseInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_polardb_databases.default.id}_fake"]`,
			"data_base_instance_id": `"${alibabacloudstack_polardb_databases.default.DataBaseInstanceId}_fake"`,
		}),
	}

	data_base_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_polardb_databases.default.id}"]`,
			"data_base_name": `"${alibabacloudstack_polardb_databases.default.DataBaseName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_polardb_databases.default.id}_fake"]`,
			"data_base_name": `"${alibabacloudstack_polardb_databases.default.DataBaseName}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_polardb_databases.default.id}"]`,
			"status": `"${alibabacloudstack_polardb_databases.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_polardb_databases.default.id}_fake"]`,
			"status": `"${alibabacloudstack_polardb_databases.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_databases.default.id}"]`,

			"data_base_instance_id": `"${alibabacloudstack_polardb_databases.default.DataBaseInstanceId}"`,
			"data_base_name":        `"${alibabacloudstack_polardb_databases.default.DataBaseName}"`,
			"status":                `"${alibabacloudstack_polardb_databases.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_databases.default.id}_fake"]`,

			"data_base_instance_id": `"${alibabacloudstack_polardb_databases.default.DataBaseInstanceId}_fake"`,
			"data_base_name":        `"${alibabacloudstack_polardb_databases.default.DataBaseName}_fake"`,
			"status":                `"${alibabacloudstack_polardb_databases.default.Status}_fake"`}),
	}

	AlibabacloudstackPolardbDatabasesCheckInfo.dataSourceTestCheck(t, rand, idsConf, data_base_instance_idConf, data_base_nameConf, statusConf, allConf)
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

func testAccCheckAlibabacloudstackPolardbDatabasesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackPolardbDatabases%d"
}


	resource "alibabacloudstack_polardb_instance" "instance" {
		engine            = "MySQL"
		engine_version    = "5.7"
		instance_name = "${var.name}"
		db_instance_storage_type= "local_ssd"
		db_instance_storage = 5
		db_instance_class = "rds.mysql.t1.small"
		zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}
	resource "alibabacloudstack_polardb_database" "default" {
		data_base_instance_id = "resource.alibabacloudstack_polardb_instance.instance.id"
		data_base_description = "自动化生成测试"
		data_base_name        = "tftest"
		character_set_name = "utf8"
	}



data "alibabacloudstack_polardb_databases" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
