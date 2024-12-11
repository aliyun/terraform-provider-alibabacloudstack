package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackRdsDatabasesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_databases.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_databases.default.id}_fake"]`,
		}),
	}

	data_base_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_rds_databases.default.id}"]`,
			"data_base_instance_id": `"${alibabacloudstack_rds_databases.default.DataBaseInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_rds_databases.default.id}_fake"]`,
			"data_base_instance_id": `"${alibabacloudstack_rds_databases.default.DataBaseInstanceId}_fake"`,
		}),
	}

	data_base_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_rds_databases.default.id}"]`,
			"data_base_name": `"${alibabacloudstack_rds_databases.default.DataBaseName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_rds_databases.default.id}_fake"]`,
			"data_base_name": `"${alibabacloudstack_rds_databases.default.DataBaseName}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_rds_databases.default.id}"]`,
			"status": `"${alibabacloudstack_rds_databases.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_rds_databases.default.id}_fake"]`,
			"status": `"${alibabacloudstack_rds_databases.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_databases.default.id}"]`,

			"data_base_instance_id": `"${alibabacloudstack_rds_databases.default.DataBaseInstanceId}"`,
			"data_base_name":        `"${alibabacloudstack_rds_databases.default.DataBaseName}"`,
			"status":                `"${alibabacloudstack_rds_databases.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_databases.default.id}_fake"]`,

			"data_base_instance_id": `"${alibabacloudstack_rds_databases.default.DataBaseInstanceId}_fake"`,
			"data_base_name":        `"${alibabacloudstack_rds_databases.default.DataBaseName}_fake"`,
			"status":                `"${alibabacloudstack_rds_databases.default.Status}_fake"`}),
	}

	AlibabacloudstackRdsDatabasesCheckInfo.dataSourceTestCheck(t, rand, idsConf, data_base_instance_idConf, data_base_nameConf, statusConf, allConf)
}

var existAlibabacloudstackRdsDatabasesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#":    "1",
		"databases.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRdsDatabasesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"databases.#": "0",
	}
}

var AlibabacloudstackRdsDatabasesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_rds_databases.default",
	existMapFunc: existAlibabacloudstackRdsDatabasesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRdsDatabasesMapFunc,
}

func testAccCheckAlibabacloudstackRdsDatabasesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRdsDatabases%d"
}






data "alibabacloudstack_rds_databases" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
