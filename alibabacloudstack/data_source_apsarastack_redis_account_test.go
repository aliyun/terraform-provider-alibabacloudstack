package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackRedisAccountsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_accounts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_accounts.default.id}_fake"]`,
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_accounts.default.id}"]`,
			"account_name": `"${alibabacloudstack_redis_accounts.default.AccountName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_redis_accounts.default.id}_fake"]`,
			"account_name": `"${alibabacloudstack_redis_accounts.default.AccountName}_fake"`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_accounts.default.id}"]`,
			"instance_id": `"${alibabacloudstack_redis_accounts.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_accounts.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_redis_accounts.default.InstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_accounts.default.id}"]`,

			"account_name": `"${alibabacloudstack_redis_accounts.default.AccountName}"`,
			"instance_id":  `"${alibabacloudstack_redis_accounts.default.InstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_accounts.default.id}_fake"]`,

			"account_name": `"${alibabacloudstack_redis_accounts.default.AccountName}_fake"`,
			"instance_id":  `"${alibabacloudstack_redis_accounts.default.InstanceId}_fake"`}),
	}

	AlibabacloudstackRedisAccountsCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, instance_idConf, allConf)
}

var existAlibabacloudstackRedisAccountsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":    "1",
		"accounts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRedisAccountsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackRedisAccountsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_redis_accounts.default",
	existMapFunc: existAlibabacloudstackRedisAccountsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRedisAccountsMapFunc,
}

func testAccCheckAlibabacloudstackRedisAccountsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRedisAccounts%d"
}






data "alibabacloudstack_redis_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
