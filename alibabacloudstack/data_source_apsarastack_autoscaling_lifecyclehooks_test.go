package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackEssLifecycleHooksDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_lifecycle_hooks.default"
	name := fmt.Sprintf("tf-testacc-esslifehookk%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig)
	
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}",
			"name_regex":       "${alibabacloudstack_ess_lifecycle_hook.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}",
			"name_regex":       "${alibabacloudstack_ess_lifecycle_hook.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_ess_lifecycle_hook.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_ess_lifecycle_hook.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}",
			"ids":              []string{"${alibabacloudstack_ess_lifecycle_hook.default.id}"},
			"name_regex":       "${alibabacloudstack_ess_lifecycle_hook.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}",
			"ids":              []string{"${alibabacloudstack_ess_lifecycle_hook.default.id}_fake"},
			"name_regex":       "${alibabacloudstack_ess_lifecycle_hook.default.name}",
		}),
	}

	var existEsslifecyclehooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"hooks.#":                       "1",
			"hooks.0.name":                  fmt.Sprintf("tf-testAccDataSourceLcHooks-%d", rand),
			"hooks.0.scaling_group_id":      CHECKSET,
			"hooks.0.default_result":        CHECKSET,
			"hooks.0.heartbeat_timeout":     "400",
			"hooks.0.lifecycle_transition":  "SCALE_OUT",
			"hooks.0.notification_arn":      CHECKSET,
			"hooks.0.notification_metadata": "helloworld",
		}
	}

	var fakeEsslifecyclehooksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"hooks.#": "0",
			"ids.#":   "0",
		}
	}

	var essLifecyclehooksCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEsslifecyclehooksMapFunc,
		fakeMapFunc:  fakeEsslifecyclehooksMapFunc,
	}

	essLifecyclehooksCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}
resource "alibabacloudstack_ess_lifecycle_hook" "default" {
  scaling_group_id      = "${alibabacloudstack_ess_scaling_group.default.id}"
  name                  = "${var.name}"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}

`, ECSInstanceCommonTestCase, name)
}
