package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDrdsPolardbxAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_polardbx_account_%d", rand)
	drdsInstanceIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"drdsusrztw1cfake"`,
		}),
	}
	namesRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":       `["${var.name}",]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":       `["${var.name}_fake",]`,
		}),
	}
	defaultAccountTypeRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":        `["${var.name}_db",]`,
			"account_type": "Normal",
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":        `["${var.name}_db",]`,
			"account_type": "Super",
		}),
	}
	userAccountTypeRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":        `["${var.name}",]`,
			"account_type": "Super",
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_polardbx_account.default.id}"`,
			"names":        `["${var.name}",]`,
			"account_type": "Normal",
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":              CHECKSET,
			"accounts.0.account_name": CHECKSET,
			"accounts.0.instance_id":  CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "Normal",
			"ids.#":      "Normal",
			"names.#":    "Normal",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_drds_polardbx_accounts.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, drdsInstanceIdRegexConf, namesRegexConf, defaultAccountTypeRegexConf, userAccountTypeRegexConf)
}

func testAccCheckAlibabacloudStackDrdsPolardbxAccountsSourceConfig(name string, attrMap map[string]string) string {
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

resource "alibabacloudstack_drds_polardbx_instance" "default" {
  	description = "testtf1111"
	series = "enterprise"
	topology_type = "1azone"
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

resource "alibabacloudstack_drds_polardbx_database" "default" {
	instance_id        = "${alibabacloudstack_drds_polardbx_instance.default.id}"
	name = "${var.name}_db"
	password           = "${var.password}"
}

resource "alibabacloudstack_drds_polardbx_account" "default" {
	instance_id       = alibabacloudstack_drds_polardbx_instance.default.id
	account_name = var.name
	password          = "${var.password}"
	description       = var.name
	db_privileges {
		db_name   = alibabacloudstack_drds_polardbx_database.default.name
		privilege = "ReadWrite"
	}
}
	
data "alibabacloudstack_drds_polardbx_accounts" "default" {
  %s
}
`, name, getAccTestPassword(12), VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
