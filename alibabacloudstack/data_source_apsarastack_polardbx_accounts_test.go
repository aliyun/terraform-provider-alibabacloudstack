package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardbx_accounts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_polardbx_account_%d", rand),
		dataSourcePolardbxAccountsConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":              CHECKSET,
			"accounts.0.account_name": CHECKSET,
			"accounts.0.instance_id":  CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existMapFunc,
		fakeMapFunc:       fakeMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_polardbx_account.default.instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "drdsusrztw1cfake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_polardbx_account.default.instance_id}",
				"names":       []string{"${var.name}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "${alibabacloudstack_polardbx_account.default.instance_id}",
				"names":       []string{"${var.name}_fake"},
			}),
		},
	)
}

func dataSourcePolardbxAccountsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s

%s

resource "alibabacloudstack_polardbx_account" "default" {
	instance_id  = local.polardbx_instance.id
	account_name = var.name
	password     = "${random_password.password.0.result}"
	description  = var.name
}
`, name, RandomPasswordTestCase(12, 1), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
