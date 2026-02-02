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

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_dbinstance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake_instance_id"},
		}),
	}
	
	dbInstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}_fake",
		}),
	}

	networkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"network_type":   string(Vpc),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"network_type":   string(Classic),
		}),
	}

	engineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"engine":         "${alibabacloudstack_polardb_dbinstance.default.engine}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"engine":         "fake_engine_id",
		}),
	}

	engineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"engine_version": "${alibabacloudstack_polardb_dbinstance.default.engine_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
			"engine_version": "fake_engineversion_id",
		}),
	}

	instanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}",
			"db_instance_class": "${alibabacloudstack_polardb_dbinstance.default.db_instance_class}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"db_instance_id":    "${alibabacloudstack_polardb_dbinstance.default.id}",
			"db_instance_class": "fake_instance_class",
		}),
	}

	vswtichIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${alibabacloudstack_polardb_dbinstance.default.vswitch_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vswitch_id": "${alibabacloudstack_polardb_dbinstance.default.vswitch_id}_fake",
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}_fake",
		}),
	}

	var existPolardbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"db_instances.#":                         "1",
			"db_instances.0.id":                      CHECKSET,
			"db_instances.0.db_instance_id":          CHECKSET,
			"db_instances.0.db_instance_description": name,
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
	polardbInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, dbInstanceIdConf, networkTypeConf, engineConf, engineVersionConf, instanceClassConf, vswtichIdConf, vpcIdConf)
}

func dataSourcePolardbInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s
data "alibabacloudstack_polardb_instance_types" "intel" {
	cpu_type = "intel"
	sorted_by = "CPU"
}

data "alibabacloudstack_polardb_instance_types" "anyone" {
	sorted_by = "CPU"
}

locals {
	instance_type = length(data.alibabacloudstack_polardb_instance_types.intel.instance_types) > 0 ? data.alibabacloudstack_polardb_instance_types.intel.instance_types.0 : data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  engine                    = local.instance_type.engine
  engine_version            = local.instance_type.engine_version
  instance_name             = var.name
  db_instance_storage_type  = local.instance_type.storage_type
  db_instance_storage       = local.instance_type.storage_min
  db_instance_class         = local.instance_type.id
  zone_id                   = data.alibabacloudstack_zones.default.zones.0.id
  vswitch_id                = alibabacloudstack_vpc_vswitch.default.id
  cpu_type                  = local.instance_type.cpu_type
}
`, name, VSwitchCommonTestCase)
}
