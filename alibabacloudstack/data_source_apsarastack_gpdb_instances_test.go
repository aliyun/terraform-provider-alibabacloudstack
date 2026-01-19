package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackGpdbInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_gpdb_instances.default"
	name := fmt.Sprintf("tf-testacc-gpdbinstance%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGpdbInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_gpdb_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_gpdb_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_gpdb_instance.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_gpdb_instance.default.id}"},
			"name_regex": "${alibabacloudstack_gpdb_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_gpdb_instance.default.id}_fake"},
			"name_regex": "${alibabacloudstack_gpdb_instance.default.description}_fake",
		}),
	}

	var existGpdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"instances.#":                              "1",
			"instances.0.id":                           CHECKSET,
			"instances.0.description":                  name,
			"instances.0.region_id":                    CHECKSET,
			"instances.0.availability_zone":            CHECKSET,
			"instances.0.creation_time":                CHECKSET,
			"instances.0.status":                       CHECKSET,
			"instances.0.engine":                       "gpdb",
			"instances.0.engine_version":               CHECKSET,
			"instances.0.instance_class":               CHECKSET,
			"instances.0.instance_group_count":         CHECKSET,
			"instances.0.instance_network_type":        "VPC",
			"instances.0.charge_type":                  CHECKSET,
		}
	}

	var fakeGpdbInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var gpdbInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existGpdbInstancesMapFunc,
		fakeMapFunc:  fakeGpdbInstancesMapFunc,
	}
	gpdbInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceGpdbInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_gpdb_instance_types" "default" {
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "10.1.0.0/16"
  name              = var.name
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_gpdb_instance" "default" {
  vswitch_id                  = alibabacloudstack_vswitch.default.id
  engine                      = "gpdb"
  engine_version              = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version
  instance_class              = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id
  db_instance_mode            = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.db_instance_mode
  db_instance_storage_type    = "local_ssd"
  description                 = var.name
  seg_node_num                = "2"
  network_type                = "VPC"
  cpu_type                    = "Intel"
}
`, name, DataZoneCommonTestCase)
}
