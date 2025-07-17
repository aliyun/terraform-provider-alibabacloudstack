package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_accounts.default"
	name := fmt.Sprintf("tf-testAccPolardbAccounts%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbAccountsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}",
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}_fake",
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
	}

	name_regex_Conf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}",
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}_fake",
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.instance.id}",
		}),
	}

	AlibabacloudstackPolardbAccountsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, account_nameConf, name_regex_Conf)
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

func dataSourcePolardbAccountsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
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
	data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	account_description = "test"
	account_name        = "polardb_account"
	account_password = "%s"
	account_type ="Normal"
}
	
`, name, getAccTestPassword(12))
}
