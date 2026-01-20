package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackGpdbAccountsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_gpdb_accounts.default"
	rand := getAccTestRandInt(100, 999)
	name := fmt.Sprintf("tf_account%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackGpdbAccountsDataSourceName)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids": []string{"${alibabacloudstack_gpdb_account.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids": []string{"${alibabacloudstack_gpdb_account.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"name_regex": "${alibabacloudstack_gpdb_account.default.account_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"name_regex": "${alibabacloudstack_gpdb_account.default.account_name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids":    []string{"${alibabacloudstack_gpdb_account.default.id}"},
			"status": "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids":    []string{"${alibabacloudstack_gpdb_account.default.id}"},
			"status": "Creating",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids":        []string{"${alibabacloudstack_gpdb_account.default.id}"},
			"name_regex": "${alibabacloudstack_gpdb_account.default.account_name}",
			"status":     "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${local.gpdb_instance_id}",
			"ids":        []string{"${alibabacloudstack_gpdb_account.default.id}_fake"},
			"name_regex": "${alibabacloudstack_gpdb_account.default.account_name}_fake",
			"status":     "Creating",
		}),
	}
	var existAlibabacloudStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"accounts.#":                     "1",
			"accounts.0.id":                  CHECKSET,
			"accounts.0.account_name":        name,
			"accounts.0.account_description": name,
			"accounts.0.db_instance_id":      CHECKSET,
			"accounts.0.status":              "Active",
		}
	}
	var fakeAlibabacloudStackGpdbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackGpdbAccountsCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existAlibabacloudStackGpdbAccountsDataSourceNameMapFunc,
		fakeMapFunc:       fakeAlibabacloudStackGpdbAccountsDataSourceNameMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	alibabacloudstackGpdbAccountsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlibabacloudStackGpdbAccountsDataSourceName(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
	%s

	%s
	
	resource "alibabacloudstack_gpdb_account" "default" {
	  account_name        = var.name
	  db_instance_id      = local.gpdb_instance_id
	  account_password    = random_password.password.0.result
	  account_description = var.name
	}
		`, name, GpdbCommonTestCase(), RandomPasswordTestCase(12, 1))
}
