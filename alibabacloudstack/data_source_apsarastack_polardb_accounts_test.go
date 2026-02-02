package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_accounts.default"
	name := fmt.Sprintf("tf_acc%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbAccountsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
	}

	account_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}",
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}_fake",
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
	}

	name_regex_Conf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}",
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_polardb_account.default.id}_fake"},
			"account_name":   "${alibabacloudstack_polardb_account.default.account_name}_fake",
			"db_instance_id": "${local.polardb_dbinstance_id}",
		}),
	}
	var AlibabacloudstackPolardbAccountsDataCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existAlibabacloudstackPolardbAccountsDataMapFunc,
		fakeMapFunc:       fakeAlibabacloudstackPolardbAccountsDataMapFunc,
		ExternalProviders: testAccExternalProviders,
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

func dataSourcePolardbAccountsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
}
%s

%s

resource "alibabacloudstack_polardb_account" "default" {
	data_base_instance_id = "${local.polardb_dbinstance_id}"
	account_description = "test"
	account_name        = var.name
	account_password = random_password.password.0.result
	account_type ="Normal"
}
	
`, name, PolarDBCommonTestCase("MySQL", false), RandomPasswordTestCase(12, 1))
}
