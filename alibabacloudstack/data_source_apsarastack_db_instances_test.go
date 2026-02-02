package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDBInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_db_instances.default"
	name := fmt.Sprintf("tf-testacc-dbinstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDBInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_db_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_db_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_db_instance.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Creating",
		}),
	}

	engineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine": "MySQL",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"engine": "SQLServer",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_db_instance.default.instance_name}",
			"ids":        []string{"${alibabacloudstack_db_instance.default.id}"},
			"status":     "Running",
			"engine":     "MySQL",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
			"ids":        []string{"${alibabacloudstack_db_instance.default.id}_fake"},
			"status":     "Creating",
			"engine":     "SQLServer",
		}),
	}

	var existDBInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"instances.#":                         "1",
			"instances.0.id":                      CHECKSET,
			"instances.0.name":                    name,
			"instances.0.charge_type":             CHECKSET,
			"instances.0.db_type":                 "Primary",
			"instances.0.region_id":               CHECKSET,
			"instances.0.create_time":             CHECKSET,
			"instances.0.status":                  "Running",
			"instances.0.engine":                  "MySQL",
			"instances.0.engine_version":          CHECKSET,
			"instances.0.net_type":                CHECKSET,
			"instances.0.connection_mode":         CHECKSET,
			"instances.0.instance_type":           CHECKSET,
			"instances.0.availability_zone":       CHECKSET,
			"instances.0.readonly_instance_ids.#": "0",
			"instances.0.vpc_id":                  CHECKSET,
			"instances.0.vswitch_id":              CHECKSET,
			"instances.0.connection_string":       CHECKSET,
			"instances.0.port":                    CHECKSET,
			"instances.0.instance_storage":        CHECKSET,
		}
	}

	var fakeDBInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var dbInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDBInstancesMapFunc,
		fakeMapFunc:  fakeDBInstancesMapFunc,
	}
	dbInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, engineConf, allConf)
}

func dataSourceDBInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s
`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
