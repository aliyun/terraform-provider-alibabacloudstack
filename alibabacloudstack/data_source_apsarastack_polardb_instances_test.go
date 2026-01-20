package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_dbinstances.default"
	name := fmt.Sprintf("tf-testacc-polardbinstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbInstancesConfigDependence)

	dbInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}_fake",
		}),
	}

	dbInstanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_class": "${alibabacloudstack_polardb_dbinstance.default.db_instance_class}",
			"status":            "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_class": "rds.mysql.t1.fake",
			"status":            "Running",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}",
			"db_instance_class": "${alibabacloudstack_polardb_dbinstance.default.db_instance_class}",
			"status":            "Running",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}_fake",
			"db_instance_class": "${alibabacloudstack_polardb_dbinstance.default.db_instance_class}_fake",
			"status":            "Running",
		}),
	}

	var existPolardbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"db_instances.#":                                  "1",
			"db_instances.0.id":                               CHECKSET,
			"db_instances.0.db_instance_id":                   CHECKSET,
			"db_instances.0.db_instance_description":          name,
		}
	}

	var fakePolardbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"db_instances.#": "0",
		}
	}

	var polardbInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbInstancesMapFunc,
		fakeMapFunc:  fakePolardbInstancesMapFunc,
	}
	polardbInstancesCheckInfo.dataSourceTestCheck(t, rand, dbInstanceIdConf, dbInstanceClassConf, allConf)
}

func dataSourcePolardbInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s
data "alibabacloudstack_polardb_instance_types" "anyone" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  engine                    = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.engine
  engine_version            = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.engine_version
  instance_name             = var.name
  db_instance_storage_type  = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.storage_type
  db_instance_storage       = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.storage_min
  db_instance_class         = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.id
  zone_id                   = data.alibabacloudstack_zones.default.zones.0.id
  vswitch_id                = alibabacloudstack_vpc_vswitch.default.id
  cpu_type                  = data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0.cpu_type
}
`, name, VSwitchCommonTestCase)
}
