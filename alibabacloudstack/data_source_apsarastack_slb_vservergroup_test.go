package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackSlbVServerGroupsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_v_server_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_v_server_groups.default.id}_fake"]`,
		}),
	}

	load_balancer_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_v_server_groups.default.id}"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_v_server_groups.default.LoadBalancerId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_v_server_groups.default.id}_fake"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_v_server_groups.default.LoadBalancerId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_v_server_groups.default.id}"]`,

			"load_balancer_id": `"${alibabacloudstack_slb_v_server_groups.default.LoadBalancerId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_v_server_groups.default.id}_fake"]`,

			"load_balancer_id": `"${alibabacloudstack_slb_v_server_groups.default.LoadBalancerId}_fake"`}),
	}

	AlibabacloudstackSlbVServerGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, load_balancer_idConf, allConf)
}

var existAlibabacloudstackSlbVServerGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbVServerGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackSlbVServerGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_v_server_groups.default",
	existMapFunc: existAlibabacloudstackSlbVServerGroupsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbVServerGroupsMapFunc,
}

func testAccCheckAlibabacloudstackSlbVServerGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbVServerGroups%d"
}






data "alibabacloudstack_slb_v_server_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
