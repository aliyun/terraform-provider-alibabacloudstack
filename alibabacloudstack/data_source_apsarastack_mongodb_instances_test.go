package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMongoDBInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_mongodb_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMongoDBInstance_datasource_%d", rand),
		dataSourceMongodbInstancesConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                CHECKSET,
			"instances.0.name":           fmt.Sprintf("tf-testAccMongoDBInstance_datasource_%d", rand),
			"instances.0.instance_class": "dds.mongo.mid",
			"instances.0.engine":         "MongoDB",
			"instances.0.engine_version": "3.4",
			"instances.0.charge_type":    string(PostPaid),
			"instances.0.storage":        "10",
			"instances.0.instance_type":  "replicate",
			"instances.0.id":             CHECKSET,
			"instances.0.creation_time":  CHECKSET,
			"instances.0.region_id":      CHECKSET,
			"instances.0.status":         CHECKSET,
			"instances.0.network_type":   CHECKSET,
			"instances.0.lock_mode":      CHECKSET,
			"ids.#":                      "1",
			"ids.0":                      CHECKSET,
			"names.#":                    "1",
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
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
				"name_regex": "${alibabacloudstack_mongodb_instance.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_mongodb_instance.default.name}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_mongodb_instance.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_mongodb_instance.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex":    "${alibabacloudstack_mongodb_instance.default.name}",
				"instance_type": "replicate",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex":    "${alibabacloudstack_mongodb_instance.default.name}",
				"instance_type": "sharding",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex":     "${alibabacloudstack_mongodb_instance.default.name}",
				"instance_class": "dds.mongo.mid",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex":     "${alibabacloudstack_mongodb_instance.default.name}",
				"instance_class": "test.rds.mid",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex":        "${alibabacloudstack_mongodb_instance.default.name}",
				"availability_zone": "${data.alibabacloudstack_zones.default.zones.0.id}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex":        "${alibabacloudstack_mongodb_instance.default.name}",
				"availability_zone": "test_zone",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex":        "${alibabacloudstack_mongodb_instance.default.name}",
				"ids":               []string{"${alibabacloudstack_mongodb_instance.default.id}"},
				"availability_zone": "${data.alibabacloudstack_zones.default.zones.0.id}",
				"instance_type":     "replicate",
				"instance_class":    "dds.mongo.mid",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex":        "${alibabacloudstack_mongodb_instance.default.name}_fake",
				"ids":               []string{"${alibabacloudstack_mongodb_instance.default.id}"},
				"availability_zone": "${data.alibabacloudstack_zones.default.zones.0.id}",
				"instance_type":     "replicate",
				"instance_class":    "dds.mongo.mid",
			}),
		},
	)
}

func dataSourceMongodbInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "MongoDB"
}

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = "${alibabacloudstack_vpc.default.id}"
	cidr_block        = "172.16.0.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id           = alibabacloudstack_vswitch.default.id
	engine_version       = "3.4"
	db_instance_class    = "dds.mongo.mid"
	db_instance_storage  = "10"
	name                 = "${var.name}"
	storage_engine       = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor   = "3"
}
`, name)
}
