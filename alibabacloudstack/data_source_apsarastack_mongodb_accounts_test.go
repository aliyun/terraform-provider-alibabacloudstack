package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMongodbAccountsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	account_password := getAccTestPassword(12)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"ids":         `["${alibabacloudstack_mongodb_account.default.id}"]`,
			"instance_id": `"${alibabacloudstack_mongodb_instance.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"ids":         `["${alibabacloudstack_mongodb_account.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_mongodb_instance.default.id}"`,
		}),
	}

	account_nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"instance_id":        `"${alibabacloudstack_mongodb_instance.default.id}"`,
			"account_name_regex": `"${alibabacloudstack_mongodb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"instance_id":        `"${alibabacloudstack_mongodb_instance.default.id}"`,
			"account_name_regex": `"${alibabacloudstack_mongodb_account.default.account_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"ids":                `["${alibabacloudstack_mongodb_account.default.id}"]`,
			"instance_id":        `"${alibabacloudstack_mongodb_instance.default.id}"`,
			"account_name_regex": `"${alibabacloudstack_mongodb_account.default.account_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand, account_password, map[string]string{
			"ids":                `["${alibabacloudstack_mongodb_account.default.id}"]`,
			"instance_id":        `"${alibabacloudstack_mongodb_instance.default.id}"`,
			"account_name_regex": `"${alibabacloudstack_mongodb_account.default.account_name}_fake"`,
		}),
	}
	AlibabacloudstackMongodbAccountsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameRegexConf, allConf)
}

var existAlibabacloudstackMongodbAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":                "1",
		"accounts.0.account_name":   CHECKSET,
		"accounts.0.character_type": CHECKSET,
		"accounts.0.instance_id":    CHECKSET,
		"accounts.0.status":         CHECKSET,
		"accounts.0.account_type":   CHECKSET,
	}
}

var fakeAlibabacloudstackMongodbAccountsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AlibabacloudstackMongodbAccountsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_mongodb_accounts.default",
	existMapFunc: existAlibabacloudstackMongodbAccountsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackMongodbAccountsDataMapFunc,
}

func testAccCheckAlibabacloudstackMongodbAccountsDataSourceConfig(rand int, account_password string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackMongodbAccounts%d"
}

%s
resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
	engine_version      = "3.0"
	db_instance_class   = "dds.mongo.mid"
	db_instance_storage = "10"
	name                = "${var.name}"
	storage_engine      = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor = "3"
  }

resource "alibabacloudstack_mongodb_account" "default" {
	account_name = "testaccountv1"
	account_password = "%s"
	instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}

data "alibabacloudstack_mongodb_accounts" "default" {
	%s
}

`, rand, VSwitchCommonTestCase, account_password, strings.Join(pairs, "\n   "))
	return config
}
