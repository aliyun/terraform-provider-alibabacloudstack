package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackEssScalingConfigurationsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_scaling_configurations.default"
	name := fmt.Sprintf("tf-essnotifications%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig)
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}_fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_configuration.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_configuration.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}",
			"ids": []string{"${alibabacloudstack_ess_scaling_configuration.default.id}"},
			"name_regex": "${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_configuration.default.scaling_group_id}_fake",
			"ids": []string{"${alibabacloudstack_ess_scaling_configuration.default.id}_fake"},
			"name_regex": "${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}",
		}),
	}

	var existEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			 "ids.#":                                       "1",
			 "names.#":                                     "1",
			 "configurations.#":                            "1",
			 "configurations.0.name":                       name,
			 "configurations.0.scaling_group_id":           CHECKSET,
			 "configurations.0.image_id":                   CHECKSET,
			 "configurations.0.instance_type":              CHECKSET,
			 "configurations.0.security_group_id":          CHECKSET,
			 "configurations.0.creation_time":              CHECKSET,
			 "configurations.0.system_disk_category":       CHECKSET,
			 "configurations.0.system_disk_size":           CHECKSET,
			 "configurations.0.internet_max_bandwidth_in":  CHECKSET,
			 "configurations.0.internet_max_bandwidth_out": CHECKSET,
		}
	}

	var fakeEssScalingconfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			 "configurations.#": "0",
			 "ids.#":            "0",
			 "names.#":          "0",
		}
	}

	var essScalingconfigurationsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEssScalingconfigurationsMapFunc,
		fakeMapFunc:  fakeEssScalingconfigurationsMapFunc,
	}

	essScalingconfigurationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackEssScalingconfigurationsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
%s
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
	scaling_configuration_name = var.name
	scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
	image_id = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type = "${local.default_instance_type_id}"
	security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
	force_delete = true
	active = true
	enable = true
	deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
}

`, name, ECSInstanceCommonTestCase)
}
