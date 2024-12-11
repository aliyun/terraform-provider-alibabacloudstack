package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackMongoDBInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_mongodb_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_mongodb_instance.default.name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_mongodb_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_mongodb_instance.default.id}_fake"]`,
		}),
	}

	//	tagsConf := dataSourceTestAccConfig{
	//		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_mongodb_instance.default.name}"`,
	//			"tags":       `{Created = "TF"}`,
	//		}),
	//		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
	//			"name_regex": `"${alibabacloudstack_mongodb_instance.default.name}"`,
	//			"tags":       `{Created = "TF1"}`,
	//		}),
	//	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":    `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"instance_type": `"replicate"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":    `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"instance_type": `"sharding"`,
		}),
	}
	instanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"instance_class": `"dds.mongo.mid"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"instance_class": `"test.rds.mid"`,
		}),
	}
	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"availability_zone": `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"availability_zone": `"test_zone"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_mongodb_instance.default.name}"`,
			"ids":               `["${alibabacloudstack_mongodb_instance.default.id}"]`,
			"availability_zone": `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
			"instance_type":     `"replicate"`,
			"instance_class":    `"dds.mongo.mid"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_mongodb_instance.default.name}_fake"`,
			"ids":               `["${alibabacloudstack_mongodb_instance.default.id}"]`,
			"availability_zone": `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
			"instance_type":     `"replicate"`,
			"instance_class":    `"dds.mongo.mid"`,
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
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
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_mongodb_instances.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, instanceTypeConf, instanceClassConf, availabilityZoneConf, allConf)
}

func testAccCheckAlibabacloudStackMongoDBDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "name" {
  default = "tf-testAccMongoDBInstance_datasource_%d"
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
  vswitch_id          = alibabacloudstack_vswitch.default.id
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = "${var.name}"
}
data "alibabacloudstack_mongodb_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
