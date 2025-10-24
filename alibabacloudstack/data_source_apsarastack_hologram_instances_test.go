package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackHologramInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_hologram_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccHologramInstances-%d", rand),
		dataSourceHologramInstancesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_hologram_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_hologram_instance.default.instance_name}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_hologram_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_hologram_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_hologram_instance.default.instance_name}",
			"ids":        []string{"${alibabacloudstack_hologram_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_hologram_instance.default.instance_name}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_hologram_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existHologramInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"instances.#":                 "1",
			"instances.0.instance_name":   fmt.Sprintf("tf-testAccHologramInstances-%d", rand),
			"instances.0.creation_time":   CHECKSET,
			"instances.0.compute_type":    CHECKSET,
			"instances.0.cpu":             CHECKSET,
			"instances.0.node":            CHECKSET,
			"instances.0.cluster":         CHECKSET,
			"instances.0.instance_status": CHECKSET,
			"instances.0.version":         CHECKSET,
		}
	}

	var fakeHologramInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var HologramInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existHologramInstancesMapFunc,
		fakeMapFunc:  fakeHologramInstancesMapFunc,
	}

	HologramInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceHologramInstancesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_hologram_instance" "default" {
  	compute_type = "Standard"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	cpu = "intel"
	node = "2"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	instance_name = "${var.name}"
	cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
}

`, name, VSwitchCommonTestCase)
}
