package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackLindormInstanceDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_lindorm_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-lindorm-%d", rand),
		dataSourceLindormInstanceDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_lindorm_instance.default.instance_alias}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_lindorm_instance.default.instance_alias}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_lindorm_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_lindorm_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_lindorm_instance.default.instance_alias}",
			"ids":        []string{"${alibabacloudstack_lindorm_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_lindorm_instance.default.instance_alias}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_lindorm_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existLindormInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"instances.#":                  "1",
			"instances.0.cpu_brand":        CHECKSET,
			"instances.0.create_time":      CHECKSET,
			"instances.0.instance_id":      CHECKSET,
			"instances.0.instance_storage": CHECKSET,
			"instances.0.engine_type":      CHECKSET,
			"instances.0.ascm_create_user": CHECKSET,
			"instances.0.instance_alias":   CHECKSET,
			"instances.0.network_type":     CHECKSET,
			"instances.0.service_type":     CHECKSET,
		}
	}

	var fakeLindormInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var LindormInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existLindormInstanceMapFunc,
		fakeMapFunc:  fakeLindormInstanceMapFunc,
	}

	LindormInstanceCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceLindormInstanceDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_lindorm_instance" "default" {
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	instance_alias = "${var.name}"
	cpu_brand = "Intel"
	disk_category = "HDD"
	local_disk_size = "6T"
	engine_type = "lindorm"
	instance_type = "lindorm.g1.8c32g"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	lindorm_num = 2
}

 `, name, VSwitchCommonTestCase)
}
