package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPolardbAccountsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_account.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_account.default.id}_fake"]`,
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_polardb_account.default.id}"]`,
			"account_name": `"${alibabacloudstack_polardb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_polardb_account.default.id}_fake"]`,
			"account_name": `"${alibabacloudstack_polardb_account.default.account_name}_fake"`,
		}),
	}

	data_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_polardb_account.default.id}"]`,
			"data_instance_id": `"${alibabacloudstack_polardb_dbinstance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_polardb_account.default.id}_fake"]`,
			"data_instance_id": `"${alibabacloudstack_polardb_dbinstance.default.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_account.default.id}"]`,

			"account_name":     `"${alibabacloudstack_polardb_account.default.account_name}"`,
			"data_instance_id": `"${alibabacloudstack_polardb_dbinstance.default.id}"`}),
		fakeConfig: testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_polardb_account.default.id}_fake"]`,

			"account_name":     `"${alibabacloudstack_polardb_account.default.account_name}_fake"`,
			"data_instance_id": `"${alibabacloudstack_polardb_dbinstance.default.id}_fake"`}),
	}

	AlibabacloudstackPolardbAccountsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, data_instance_idConf, allConf)
}

var existAlibabacloudstackPolardbAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":    "1",
		"accounts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackPolardbAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackPolardbAccountsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_polardb_accounts.default",
	existMapFunc: existAlibabacloudstackPolardbAccountsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackPolardbAccountsDataMapFunc,
}

func testAccCheckAlibabacloudstackPolardbAccountsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackPolardbAccounts%d"
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
}
resource "alibabacloudstack_polardb_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "tfinstance"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}
resource "alibabacloudstack_polardb_account" "default" {
	data_base_instance_id = "${resource.alibabacloudstack_polardb_dbinstance.instance.id}"
	account_description = "test"
	account_name        = "polardb_account"
	account_password = "Test12345"
	account_type ="Normal"
}
		

data "alibabacloudstack_polardb_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
}
