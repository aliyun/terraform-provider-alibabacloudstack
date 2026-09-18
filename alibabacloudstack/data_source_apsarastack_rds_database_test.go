package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackRdsDatabasesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_rds_databases.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackRdsDatabases%d", rand),
		dataSourceRdsDatabasesConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#":    "1",
			"databases.0.id": CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"databases.#": "0",
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
				"ids": []string{"${alibabacloudstack_rds_databases.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_rds_databases.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_databases.default.id}"},
				"data_base_instance_id": "${alibabacloudstack_rds_databases.default.DataBaseInstanceId}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_databases.default.id}_fake"},
				"data_base_instance_id": "${alibabacloudstack_rds_databases.default.DataBaseInstanceId}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":            []string{"${alibabacloudstack_rds_databases.default.id}"},
				"data_base_name": "${alibabacloudstack_rds_databases.default.DataBaseName}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":            []string{"${alibabacloudstack_rds_databases.default.id}_fake"},
				"data_base_name": "${alibabacloudstack_rds_databases.default.DataBaseName}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":    []string{"${alibabacloudstack_rds_databases.default.id}"},
				"status": "${alibabacloudstack_rds_databases.default.Status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":    []string{"${alibabacloudstack_rds_databases.default.id}_fake"},
				"status": "${alibabacloudstack_rds_databases.default.Status}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_databases.default.id}"},
				"data_base_instance_id": "${alibabacloudstack_rds_databases.default.DataBaseInstanceId}",
				"data_base_name":        "${alibabacloudstack_rds_databases.default.DataBaseName}",
				"status":                "${alibabacloudstack_rds_databases.default.Status}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":                   []string{"${alibabacloudstack_rds_databases.default.id}_fake"},
				"data_base_instance_id": "${alibabacloudstack_rds_databases.default.DataBaseInstanceId}_fake",
				"data_base_name":        "${alibabacloudstack_rds_databases.default.DataBaseName}_fake",
				"status":                "${alibabacloudstack_rds_databases.default.Status}_fake",
			}),
		},
	)
}

func dataSourceRdsDatabasesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
