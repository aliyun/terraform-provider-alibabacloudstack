package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPolardbxDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_polardbx_db_%d", rand)
	polardbxInstanceIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_polardbx_database.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id": `"polardbxusrztw1cfake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":   `"${alibabacloudstack_polardbx_database.default.instance_id}"`,
			"database_name": fmt.Sprintf(`"%s"`, name),
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":   `"${alibabacloudstack_polardbx_database.default.instance_id}"`,
			"database_name": `"tf_acc_polardbx_db_fake"`,
		}),
	}
	namesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":    `"${alibabacloudstack_polardbx_database.default.instance_id}"`,
			"database_names": fmt.Sprintf(`["%s"]`, name),
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":    `"${alibabacloudstack_polardbx_database.default.instance_id}"`,
			"database_names": `["tf_acc_polardbx_db_fake"]`,
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"databases.#":               CHECKSET,
			"databases.0.database_name": CHECKSET,
			"databases.0.encode":        CHECKSET,
			"databases.0.instance_id":   CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"databases.#": "0",
			"names.#":     "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_polardbx_databases.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, polardbxInstanceIdRegexConf, idsConf, namesConf)
}

func testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s
	
resource "alibabacloudstack_polardbx_database" "default" {
    instance_id  = "${local.polardbx_instance.id}"
	database_name = "${var.name}"
	encode = "utf8mb4"
	mode = "auto"
}
	
data "alibabacloudstack_polardbx_databases" "default" {
  %s
}
`, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase(), strings.Join(pairs, "\n  "))
	return config
}
