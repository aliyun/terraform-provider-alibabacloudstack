package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDrdsAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_db_%d", rand)
	drdsInstanceIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_database.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id": `"drdsusrztw1cfake"`,
		}),
	}
	namesRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":       `["${var.name}",]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":       `["${var.name}_fake",]`,
		}),
	}
	defaultAccountTypeRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":        `["${var.name}_db",]`,
			"account_type": "0",
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":        `["${var.name}_db",]`,
			"account_type": "1",
		}),
	}
	userAccountTypeRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":        `["${var.name}",]`,
			"account_type": "1",
		}),
		fakeConfig: testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name, map[string]string{
			"instance_id":  `"${alibabacloudstack_drds_database.default.instance_id}"`,
			"names":        `["${var.name}",]`,
			"account_type": "0",
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":                   CHECKSET,
			"accounts.0.drds_account_name": CHECKSET,
			"accounts.0.instance_id":       CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_drds_databases.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, drdsInstanceIdRegexConf, namesRegexConf, defaultAccountTypeRegexConf, userAccountTypeRegexConf)
}

func testAccCheckAlibabacloudStackDrdsAccountsSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}

	variable "instance_series" {
		default = "drds.sn2.4c16g"
	}

	resource "random_password" "password" {
		length           = 12
		special          = true
		override_special = "_"
		min_lower        = 1
		min_upper        = 1
		min_numeric      = 1
	}

	%s

	resource "alibabacloudstack_drds_instance" "default" {
		description          = "${var.name}"
		zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
		instance_series      = "${var.instance_series}"
		instance_charge_type = "PostPaid"
		vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
		specification        = "drds.sn2.4c16g.8C32G"
	}

	resource "alibabacloudstack_drds_rds_instance" "default" {
		zone_id             = data.alibabacloudstack_zones.default.zones.0.id
		db_instance_storage = "20"
		storage_type        = "local_ssd"
		category            = "HighAvailability"
		db_instance_class   = "rds.mysql.s1.small"
		drds_instance_id    = alibabacloudstack_drds_instance.default.id
	}

	resource "alibabacloudstack_drds_database" "default" {
		instance_id        = "${alibabacloudstack_drds_instance.default.id}"
		drds_database_name = "${var.name}_db"
		password           = random_password.password.result
		rds_instance_ids   = [alibabacloudstack_drds_rds_instance.default.rds_instance_id,]
	}
	
	resource "alibabacloudstack_drds_account" "default" {
		instance_id       = alibabacloudstack_drds_instance.default.id
		drds_account_name = var.name
		password          = random_password.password.result
		description       = var.name
		db_privileges {
			db_name   = alibabacloudstack_drds_database.default.drds_database_name
			privilege = "R"
		}
	}
	
data "alibabacloudstack_drds_accounts" "default" {
  %s
}
`, name, VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
