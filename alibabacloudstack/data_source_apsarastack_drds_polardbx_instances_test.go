package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardbx_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sPolardbxInstancesDataSource-%d", defaultRegionToTest, rand),
		dataSourcePolardbxInstancesDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_polardbx_instance.default.description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardbx_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardbx_instance.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardbx_instance.default.description}",
			"ids":        []string{"${alibabacloudstack_polardbx_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardbx_instance.default.description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_polardbx_instance.default.id}-fakeTestAcccc"},
		}),
	}

	var existPolardbxInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"ids.0":                               CHECKSET,
			"polardbx_instances.#":                "1",
			"polardbx_instances.0.description":    fmt.Sprintf("tf-testAcc%sPolardbxInstancesDataSource-%d", defaultRegionToTest, rand),
			"polardbx_instances.0.create_time":    CHECKSET,
			"polardbx_instances.0.storage":        CHECKSET,
			"polardbx_instances.0.cpu_type":       CHECKSET,
			"polardbx_instances.0.cn_node_class":  CHECKSET,
			"polardbx_instances.0.cn_node_count":  CHECKSET,
			"polardbx_instances.0.dn_node_class":  CHECKSET,
			"polardbx_instances.0.dn_node_count":  CHECKSET,
			"polardbx_instances.0.engine_version": CHECKSET,
			"polardbx_instances.0.resource_type":  CHECKSET,
		}
	}

	var fakePolardbxInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "0",
			"polardbx_instances.#": "0",
		}
	}

	var PolardbxInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbxInstancesMapFunc,
		fakeMapFunc:  fakePolardbxInstancesMapFunc,
	}

	PolardbxInstancesCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourcePolardbxInstancesDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_polardbx_instance" "default" {
  	description = "testtf1111"
	series = "enterprise"
	topology_type = "1azone"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

 `, name, VSwitchCommonTestCase)
}
