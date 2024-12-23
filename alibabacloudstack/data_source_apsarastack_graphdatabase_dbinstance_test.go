package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackGraphdatabaseDbInstancesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_graphdatabase_db_instances.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_graphdatabase_db_instances.default.id}_fake"]`,
		}),
	}

	db_instance_descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_graphdatabase_db_instances.default.id}"]`,
			"db_instance_description": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceDescription}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_graphdatabase_db_instances.default.id}_fake"]`,
			"db_instance_description": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceDescription}_fake"`,
		}),
	}

	db_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_graphdatabase_db_instances.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_graphdatabase_db_instances.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceId}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_graphdatabase_db_instances.default.id}"]`,
			"status": `"${alibabacloudstack_graphdatabase_db_instances.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_graphdatabase_db_instances.default.id}_fake"]`,
			"status": `"${alibabacloudstack_graphdatabase_db_instances.default.Status}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_graphdatabase_db_instances.default.id}"]`,

			"db_instance_description": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceDescription}"`,
			"db_instance_id":          `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceId}"`,
			"status":                  `"${alibabacloudstack_graphdatabase_db_instances.default.Status}"`}),
		fakeConfig: testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_graphdatabase_db_instances.default.id}_fake"]`,

			"db_instance_description": `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceDescription}_fake"`,
			"db_instance_id":          `"${alibabacloudstack_graphdatabase_db_instances.default.DbInstanceId}_fake"`,
			"status":                  `"${alibabacloudstack_graphdatabase_db_instances.default.Status}_fake"`}),
	}

	AlibabacloudstackGraphdatabaseDbInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, db_instance_descriptionConf, db_instance_idConf, statusConf, allConf)
}

var existAlibabacloudstackGraphdatabaseDbInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "1",
		"instances.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackGraphdatabaseDbInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#": "0",
	}
}

var AlibabacloudstackGraphdatabaseDbInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_graphdatabase_db_instances.default",
	existMapFunc: existAlibabacloudstackGraphdatabaseDbInstancesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackGraphdatabaseDbInstancesMapFunc,
}

func testAccCheckAlibabacloudstackGraphdatabaseDbInstancesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackGraphdatabaseDbInstances%d"
}






data "alibabacloudstack_graphdatabase_db_instances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
