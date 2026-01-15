package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEssScalingGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_scaling_groups.default"
	name := fmt.Sprintf("tf-essgroup%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackEssScalinggroupsDataSourceConfig)
	
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ess_scaling_group.default.scaling_group_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ess_scaling_group.default.scaling_group_name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_group.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_group.default.id}"},
			"name_regex": "^${alibabacloudstack_ess_scaling_group.default.scaling_group_name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scaling_group.default.id}_fake"},
			"name_regex": "${alibabacloudstack_ess_scaling_group.default.scaling_group_name}_fake",
		}),
	}

	var existEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                     "1",
			"ids.#":                        "1",
			"names.#":                      "1",
			"groups.0.id":                  CHECKSET,
			"groups.0.name":                name,
			"groups.0.region_id":           CHECKSET,
			"groups.0.min_size":            "0",
			"groups.0.max_size":            "2",
			"groups.0.cooldown_time":       "20",
			"groups.0.removal_policies.#":  "2",
			"groups.0.removal_policies.0":  "OldestInstance",
			"groups.0.removal_policies.1":  "NewestInstance",
			"groups.0.load_balancer_ids.#": "0",
			"groups.0.db_instance_ids.#":   "0",
			"groups.0.vswitch_ids.#":       "1",
			"groups.0.total_capacity":      CHECKSET,
			"groups.0.active_capacity":     CHECKSET,
			"groups.0.pending_capacity":    CHECKSET,
			"groups.0.removing_capacity":   CHECKSET,
			"groups.0.creation_time":       CHECKSET,
			"groups.0.lifecycle_state":     CHECKSET,
		}
	}

	var fakeEssScalingGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essScalingGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEssScalingGroupsMapFunc,
		fakeMapFunc:  fakeEssScalingGroupsMapFunc,
	}

	essScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackEssScalinggroupsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	scaling_group_name = "${var.name}"
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

`, ECSInstanceCommonTestCase, name)
}
