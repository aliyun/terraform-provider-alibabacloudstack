package alibabacloudstack

import (
	"fmt"

	"strings"
	"testing"
)

func TestAccAlibabacloudStackEssScalingconfigurationsDataSource(t *testing.T) {
	rand := getAccTestRandInt(0, 500)
	// scalingGroupIdConf := dataSourceTestAccConfig{
	// 	existConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	// 		"scaling_group_id": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}"`,
	// 	}),
	// 	fakeConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	// 		"scaling_group_id": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}_fake"`,
	// 	}),
	// }

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}_fake"`,
		}),
	}

	//idsConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	//		"ids": `["${alibabacloudstack_ess_scaling_configuration.default.id}"]`,
	//	}),
	//	fakeConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	//		"ids": `["${alibabacloudstack_ess_scaling_configuration.default.id}_fake"]`,
	//	}),
	// }`"${alibabacloudstack_apigateway_api_groups.default.ApiGroupName}"`
	// `["${alibabacloudstack_ess_scaling_configuration.default.id}_fake}"]`
	// allConf := dataSourceTestAccConfig{
	// 	existConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	// 		// "scaling_group_id": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}"`,
	// 		// "ids":              `["${alibabacloudstack_ess_scaling_configuration.default.id}"]`,
	// 		"name_regex": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}"`,
	// 	}),
	// 	fakeConfig: testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand, map[string]string{
	// 		// "scaling_group_id": "${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}",
	// 		// "ids":              `["${alibabacloudstack_ess_scaling_configuration.default.id}_fake}"]`,
	// 		"name_regex": `"${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}"`,
	// 	}),
	// }

	var existEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"test": NOSET,
			// "ids.#":                                       "1",
			// "names.#":                                     "1",
			// "configurations.#":                            "1",
			// "configurations.0.name":                       fmt.Sprintf("tf-testAccDataSourceEssScalingRules-%d", rand),
			// "configurations.0.scaling_group_id":           CHECKSET,
			// "configurations.0.image_id":                   CHECKSET,
			// "configurations.0.instance_type":              CHECKSET,
			// "configurations.0.security_group_id":          CHECKSET,
			// "configurations.0.creation_time":              CHECKSET,
			// "configurations.0.system_disk_category":       CHECKSET,
			// "configurations.0.system_disk_size":           CHECKSET,
			// "configurations.0.internet_max_bandwidth_in":  CHECKSET,
			// "configurations.0.internet_max_bandwidth_out": CHECKSET,
			// "configurations.0.data_disks.#":               "0",
		}
	}

	var fakeEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"test": NOSET,
			// "configurations.#": "0",
			// "ids.#":            "0",
			// "names.#":          "0",
		}
	}

	var essScalingconfigurationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_ess_scaling_configurations.default",
		existMapFunc: existEssScalingconfigurationsMapFunc,
		fakeMapFunc:  fakeEssScalingconfigurationsMapFunc,
	}

	essScalingconfigurationsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s
variable "name" {
	default = "tf-testAccDataSourceEssScalingRules-%d"
}
resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
	strategy            = "Availability"
	domain              = "Default"
	granularity         = "Host"
	deployment_set_name = "example_value"
	description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
	scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
	image_id = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type = "ecs.e4.small"
	security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
	force_delete = true
	active = true
	enable = true
	zone_id = data.alibabacloudstack_zones.default.zones.0.id
	deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
}

data "alibabacloudstack_ess_scaling_configurations" "default"{
  %s
}
`, ECSInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}

// "scaling_group_id":   "${alibabacloudstack_ess_scaling_group.default.id}",
// "image_id":           "${data.alibabacloudstack_images.default.images.0.id}",
// "instance_type":      "ecs.n4.large",
// "security_group_ids": []string{"${alibabacloudstack_ecs_securitygroup.default.id}"},
// "zone_id":            "${data.alibabacloudstack_zones.default.zones.0.id}",
// "deployment_set_id":  "${alibabacloudstack_ecs_deployment_set.default.id}",
// "force_delete":       "true",
