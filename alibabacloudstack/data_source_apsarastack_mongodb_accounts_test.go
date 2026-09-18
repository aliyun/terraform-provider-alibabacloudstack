package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMongodbAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_mongodb_accounts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackMongodbAccounts%d", rand),
		dataSourceMongodbAccountsConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":                "1",
			"accounts.0.account_name":   CHECKSET,
			"accounts.0.character_type": CHECKSET,
			"accounts.0.instance_id":    CHECKSET,
			"accounts.0.status":         CHECKSET,
			"accounts.0.account_type":   CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_mongodb_account.default.id}"},
				"instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_mongodb_account.default.id}_fake"},
				"instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":        "${alibabacloudstack_mongodb_instance.default.id}",
				"account_name_regex": "${alibabacloudstack_mongodb_account.default.account_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":        "${alibabacloudstack_mongodb_instance.default.id}",
				"account_name_regex": "${alibabacloudstack_mongodb_account.default.account_name}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":                []string{"${alibabacloudstack_mongodb_account.default.id}"},
				"instance_id":        "${alibabacloudstack_mongodb_instance.default.id}",
				"account_name_regex": "${alibabacloudstack_mongodb_account.default.account_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":                []string{"${alibabacloudstack_mongodb_account.default.id}"},
				"instance_id":        "${alibabacloudstack_mongodb_instance.default.id}",
				"account_name_regex": "${alibabacloudstack_mongodb_account.default.account_name}_fake",
			}),
		},
	)
}

func dataSourceMongodbAccountsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	engine_version       = "3.0"
	db_instance_class    = "dds.mongo.mid"
	db_instance_storage  = "10"
	name                 = "${var.name}"
	storage_engine       = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor   = "3"
}

resource "alibabacloudstack_mongodb_account" "default" {
	account_name     = "testaccountv1"
	account_password = "TestAcc12345!"
	instance_id      = "${alibabacloudstack_mongodb_instance.default.id}"
}
`, name, VSwitchCommonTestCase)
}
