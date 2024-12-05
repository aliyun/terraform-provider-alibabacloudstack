package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackRdsAccountsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_accounts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_accounts.default.id}_fake"]`,
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_rds_accounts.default.id}"]`,
			"account_name": `"${alibabacloudstack_rds_accounts.default.AccountName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_rds_accounts.default.id}_fake"]`,
			"account_name": `"${alibabacloudstack_rds_accounts.default.AccountName}_fake"`,
		}),
	}

	data_base_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_rds_accounts.default.id}"]`,
			"data_base_instance_id": `"${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_rds_accounts.default.id}_fake"]`,
			"data_base_instance_id": `"${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_accounts.default.id}"]`,

			"account_name":          `"${alibabacloudstack_rds_accounts.default.AccountName}"`,
			"data_base_instance_id": `"${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_rds_accounts.default.id}_fake"]`,

			"account_name":          `"${alibabacloudstack_rds_accounts.default.AccountName}_fake"`,
			"data_base_instance_id": `"${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}_fake"`}),
	}

	AlibabacloudstackRdsAccountsCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, data_base_instance_idConf, allConf)
}

var existAlibabacloudstackRdsAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":    "1",
		"accounts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRdsAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackRdsAccountsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_rds_accounts.default",
	existMapFunc: existAlibabacloudstackRdsAccountsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRdsAccountsDataMapFunc,
}

func testAccCheckAlibabacloudstackRdsAccountsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRdsAccounts%d"
}






data "alibabacloudstack_rds_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

