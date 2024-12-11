package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackAdbAccountsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_accounts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_accounts.default.id}_fake"]`,
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_adb_accounts.default.id}"]`,
			"account_name": `"${alibabacloudstack_adb_accounts.default.AccountName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_adb_accounts.default.id}_fake"]`,
			"account_name": `"${alibabacloudstack_adb_accounts.default.AccountName}_fake"`,
		}),
	}

	account_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_adb_accounts.default.id}"]`,
			"account_type": `"${alibabacloudstack_adb_accounts.default.AccountType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_adb_accounts.default.id}_fake"]`,
			"account_type": `"${alibabacloudstack_adb_accounts.default.AccountType}_fake"`,
		}),
	}

	db_cluster_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_accounts.default.id}"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_accounts.default.DBClusterID}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_adb_accounts.default.id}_fake"]`,
			"db_cluster_id": `"${alibabacloudstack_adb_accounts.default.DBClusterID}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_accounts.default.id}"]`,

			"account_name":  `"${alibabacloudstack_adb_accounts.default.AccountName}"`,
			"account_type":  `"${alibabacloudstack_adb_accounts.default.AccountType}"`,
			"db_cluster_id": `"${alibabacloudstack_adb_accounts.default.DBClusterID}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_adb_accounts.default.id}_fake"]`,

			"account_name":  `"${alibabacloudstack_adb_accounts.default.AccountName}_fake"`,
			"account_type":  `"${alibabacloudstack_adb_accounts.default.AccountType}_fake"`,
			"db_cluster_id": `"${alibabacloudstack_adb_accounts.default.DBClusterID}_fake"`}),
	}

	AlibabacloudstackAdbAccountsCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, account_typeConf, db_cluster_idConf, allConf)
}

var existAlibabacloudstackAdbAccountsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":    "1",
		"accounts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAdbAccountsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackAdbAccountsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_adb_accounts.default",
	existMapFunc: existAlibabacloudstackAdbAccountsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAdbAccountsMapFunc,
}

func testAccCheckAlibabacloudstackAdbAccountsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAdbAccounts%d"
}






data "alibabacloudstack_adb_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
