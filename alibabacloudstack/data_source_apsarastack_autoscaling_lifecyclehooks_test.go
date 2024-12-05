package alibabacloudstack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackEssLifecycleHooksDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10, 1000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"name_regex":       `"${alibabacloudstack_ess_lifecycle_hook.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"name_regex":       `"${alibabacloudstack_ess_lifecycle_hook.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alibabacloudstack_ess_lifecycle_hook.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alibabacloudstack_ess_lifecycle_hook.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alibabacloudstack_ess_lifecycle_hook.default.id}"]`,
			"name_regex":       `"${alibabacloudstack_ess_lifecycle_hook.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"`,
			"ids":              `["${alibabacloudstack_ess_lifecycle_hook.default.id}_fake"]`,
			"name_regex":       `"${alibabacloudstack_ess_lifecycle_hook.default.name}"`,
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
		resourceId:   "data.alibabacloudstack_ess_lifecycle_hooks.default",
		existMapFunc: existEsslifecyclehooksMapFunc,
		fakeMapFunc:  fakeEsslifecyclehooksMapFunc,
	}

	essLifecyclehooksCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackEssLifecycleHooksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceLcHooks-%d"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
}
resource "alibabacloudstack_ess_lifecycle_hook" "default" {
  scaling_group_id      = "${alibabacloudstack_ess_scaling_group.default.id}"
  name                  = "${var.name}"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}

data "alibabacloudstack_ess_lifecycle_hooks" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
