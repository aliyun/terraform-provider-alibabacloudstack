package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardbx_databases.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_acc_polardbx_db_%d", rand),
		dataSourcePolardbxDatabasesConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#":               CHECKSET,
			"databases.0.database_name": CHECKSET,
			"databases.0.encode":        CHECKSET,
			"databases.0.instance_id":   CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#": "0",
			"names.#":     "0",
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
				"instance_id": "${alibabacloudstack_polardbx_database.default.instance_id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id": "polardbxusrztw1cfake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":   "${alibabacloudstack_polardbx_database.default.instance_id}",
				"database_name": fmt.Sprintf("tf_acc_polardbx_db_%d", rand),
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":   "${alibabacloudstack_polardbx_database.default.instance_id}",
				"database_name": "tf_acc_polardbx_db_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"instance_id":    "${alibabacloudstack_polardbx_database.default.instance_id}",
				"database_names": []string{fmt.Sprintf("tf_acc_polardbx_db_%d", rand)},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"instance_id":    "${alibabacloudstack_polardbx_database.default.instance_id}",
				"database_names": []string{"tf_acc_polardbx_db_fake"},
			}),
		},
	)
}

func dataSourcePolardbxDatabasesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s

resource "alibabacloudstack_polardbx_database" "default" {
	instance_id   = "${local.polardbx_instance.id}"
	database_name = "${var.name}"
	encode        = "utf8mb4"
	mode          = "auto"
}
`, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
