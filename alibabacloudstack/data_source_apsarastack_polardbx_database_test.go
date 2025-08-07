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
			"instance_id": `"${alibabacloudstack_polardbx_instance.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id": `"polardbxusrztw1cfake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":   `"${alibabacloudstack_polardbx_instance.default.instance_id}"`,
			"database_name": fmt.Sprintf(`"%s"`, name),
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxDatabasesSourceConfig(name, map[string]string{
			"instance_id":   `"${alibabacloudstack_polardbx_instance.default.instance_id}"`,
			"database_name": `"tf_acc_polardbx_db__fake"`,
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
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, polardbxInstanceIdRegexConf, idsConf)
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

variable "password" {
  default = "%s"
}

%s

resource "alibabacloudstack_polardbx_instance" "default" {
    description = "testtf1111"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

resource "alibabacloudstack_polardbx_account" "super" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	account_name = "admin"
	account_type = "Super"
	password     = "${var.password}"
	description  = "Super user"
}
	
resource "alibabacloudstack_polardbx_database" "default" {
    instance_id  = "${alibabacloudstack_polardbx_instance.default.id}"
	database_name = "${var.name}"
	encode = "utf8mb4"
	mode = "auto"
	depends_on = ["alibabacloudstack_polardbx_account.super"]
}
	
data "alibabacloudstack_polardbx_databases" "default" {
  %s
}
`, name, getAccTestPassword(12), VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
