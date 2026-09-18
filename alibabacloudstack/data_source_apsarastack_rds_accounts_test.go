package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackRdsAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_rds_accounts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackRdsAccounts%d", rand),
		dataSourceRdsAccountsConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":    "1",
			"accounts.0.id": CHECKSET,
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
				"ids": []string{"${alibabacloudstack_rds_accounts.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_rds_accounts.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":          []string{"${alibabacloudstack_rds_accounts.default.id}"},
				"account_name": "${alibabacloudstack_rds_accounts.default.AccountName}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":          []string{"${alibabacloudstack_rds_accounts.default.id}_fake"},
				"account_name": "${alibabacloudstack_rds_accounts.default.AccountName}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_accounts.default.id}"},
				"data_base_instance_id": "${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_accounts.default.id}_fake"},
				"data_base_instance_id": "${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_accounts.default.id}"},
				"account_name":          "${alibabacloudstack_rds_accounts.default.AccountName}",
				"data_base_instance_id": "${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_accounts.default.id}_fake"},
				"account_name":          "${alibabacloudstack_rds_accounts.default.AccountName}_fake",
				"data_base_instance_id": "${alibabacloudstack_rds_accounts.default.DataBaseInstanceId}_fake",
			}),
		},
	)
}

func dataSourceRdsAccountsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
