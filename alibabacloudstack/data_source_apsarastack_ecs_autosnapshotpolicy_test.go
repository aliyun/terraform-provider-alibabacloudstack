package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackEcsAutoSnapshotPoliciesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}_fake"]`,
		}),
	}

	auto_snapshot_policy_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}"]`,
			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_auto_snapshot_policies.default.AutoSnapshotPolicyId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids":                     `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}_fake"]`,
			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_auto_snapshot_policies.default.AutoSnapshotPolicyId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}"]`,

			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_auto_snapshot_policies.default.AutoSnapshotPolicyId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_auto_snapshot_policies.default.id}_fake"]`,

			"auto_snapshot_policy_id": `"${alibabacloudstack_ecs_auto_snapshot_policies.default.AutoSnapshotPolicyId}_fake"`}),
	}

	AlibabacloudstackEcsAutoSnapshotPoliciesCheckInfo.dataSourceTestCheck(t, rand, idsConf, auto_snapshot_policy_idConf, allConf)
}

var existAlibabacloudstackEcsAutoSnapshotPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#":    "1",
		"policies.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsAutoSnapshotPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#": "0",
	}
}

var AlibabacloudstackEcsAutoSnapshotPoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_auto_snapshot_policies.default",
	existMapFunc: existAlibabacloudstackEcsAutoSnapshotPoliciesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsAutoSnapshotPoliciesMapFunc,
}

func testAccCheckAlibabacloudstackEcsAutoSnapshotPoliciesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsAutoSnapshotPolicies%d"
}






data "alibabacloudstack_ecs_auto_snapshot_policies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
