package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackAutoscalingScalingGroupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_autoscaling_scaling_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_autoscaling_scaling_groups.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_autoscaling_scaling_groups.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_autoscaling_scaling_groups.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_autoscaling_scaling_groups.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_autoscaling_scaling_groups.default.ResourceGroupId}_fake"`,
		}),
	}

	scaling_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_autoscaling_scaling_groups.default.id}"]`,
			"scaling_group_name": `"${alibabacloudstack_autoscaling_scaling_groups.default.ScalingGroupName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_autoscaling_scaling_groups.default.id}_fake"]`,
			"scaling_group_name": `"${alibabacloudstack_autoscaling_scaling_groups.default.ScalingGroupName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_autoscaling_scaling_groups.default.id}"]`,

			"resource_group_id":  `"${alibabacloudstack_autoscaling_scaling_groups.default.ResourceGroupId}"`,
			"scaling_group_name": `"${alibabacloudstack_autoscaling_scaling_groups.default.ScalingGroupName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_autoscaling_scaling_groups.default.id}_fake"]`,

			"resource_group_id":  `"${alibabacloudstack_autoscaling_scaling_groups.default.ResourceGroupId}_fake"`,
			"scaling_group_name": `"${alibabacloudstack_autoscaling_scaling_groups.default.ScalingGroupName}_fake"`}),
	}

	AlibabacloudstackAutoscalingScalingGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, resource_group_idConf, scaling_group_nameConf, allConf)
}

var existAlibabacloudstackAutoscalingScalingGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackAutoscalingScalingGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackAutoscalingScalingGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_autoscaling_scaling_groups.default",
	existMapFunc: existAlibabacloudstackAutoscalingScalingGroupsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackAutoscalingScalingGroupsMapFunc,
}

func testAccCheckAlibabacloudstackAutoscalingScalingGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackAutoscalingScalingGroups%d"
}






data "alibabacloudstack_autoscaling_scaling_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
