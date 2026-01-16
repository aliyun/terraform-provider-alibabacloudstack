package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEssScalingRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_scaling_rules.default"
	name := fmt.Sprintf("tf-testacc-essrule%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_rule.default.id}_fake"},
		}),
	}

	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}",
			"type":             "SimpleScalingRule",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake",
			"type":             "TargetTrackingScalingRule",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_ess_scaling_rule.default.id}"},
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}",
			"type":             "SimpleScalingRule",
			"name_regex":       "${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_ess_scaling_rule.default.id}"},
			"scaling_group_id": "${alibabacloudstack_ess_scaling_rule.default.scaling_group_id}_fake",
			"type":             "SimpleScalingRule",
			"name_regex":       "${alibabacloudstack_ess_scaling_rule.default.scaling_rule_name}",
		}),
	}

	var existEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			 "rules.#":                  "1",
			 "ids.#":                    "1",
			 "names.#":                  "1",
			 "rules.0.id":               CHECKSET,
			 "rules.0.scaling_group_id": CHECKSET,
			 "rules.0.name":             CHECKSET,
			 "rules.0.type":             CHECKSET,
			 "rules.0.cooldown":         CHECKSET,
			 "rules.0.adjustment_type":  "TotalCapacity",
			 "rules.0.adjustment_value": "1",
			 "rules.0.scaling_rule_ari": CHECKSET,
		}
	}

	var fakeEssRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			 "rules.#": "0",
			 "ids.#":   "0",
			 "names.#": "0",
		}
	}

	var EssScalingrulesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEssRecordsMapFunc,
		fakeMapFunc:  fakeEssRecordsMapFunc,
	}

	EssScalingrulesRecordsCheckInfo.dataSourceTestCheck(t, -1, idsConf, scalingGroupIdConf, typeConf, nameRegexConf, allConf)
}

func testAccCheckAlibabacloudStackEssScalingrulesDataSourceConfig(name string) string {
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
resource "alibabacloudstack_ess_scaling_rule" "default" {
	scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
	adjustment_type = "TotalCapacity"
	adjustment_value = "1"
	cooldown = 0
}
`, name, ECSInstanceCommonTestCase)
}
