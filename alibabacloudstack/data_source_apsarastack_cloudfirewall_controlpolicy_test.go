package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackCloudfirewallControlPoliciesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudfirewallControlPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudfirewall_control_policies.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudfirewallControlPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudfirewall_control_policies.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudfirewallControlPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudfirewall_control_policies.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudfirewallControlPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudfirewall_control_policies.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackCloudfirewallControlPoliciesCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackCloudfirewallControlPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#":    "1",
		"policies.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCloudfirewallControlPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policies.#": "0",
	}
}

var AlibabacloudstackCloudfirewallControlPoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cloudfirewall_control_policies.default",
	existMapFunc: existAlibabacloudstackCloudfirewallControlPoliciesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCloudfirewallControlPoliciesMapFunc,
}

func testAccCheckAlibabacloudstackCloudfirewallControlPoliciesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCloudfirewallControlPolicies%d"
}






data "alibabacloudstack_cloudfirewall_control_policies" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
