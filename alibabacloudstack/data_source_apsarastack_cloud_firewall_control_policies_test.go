package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccCheckAlibabacloudStackCloudFirewallControlPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	var existAlibabacloudStackCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"policies.#":                  "1",
			"policies.0.description":      fmt.Sprintf("tf-testAccCloudFirewallControlPolicies-%d", rand),
			"policies.0.application_name": "ANY",
			"policies.0.acl_action":       "accept",
			"policies.0.destination_type": "net",
			"policies.0.destination":      "100.1.1.0/24",
			"policies.0.direction":        "out",
			"policies.0.proto":            "ANY",
			"policies.0.source":           "1.2.3.0/24",
			"policies.0.source_type":      "net",
		}
	}
	var fakeAlibabacloudStackCloudFirewallControlPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackEventBridgeEventBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_cloud_firewall_control_policies.default",
		existMapFunc: existAlibabacloudStackCloudFirewallControlPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackCloudFirewallControlPoliciesDataSourceNameMapFunc,
	}
	alibabacloudstackEventBridgeEventBusesCheckInfo.dataSourceTestCheck(t, rand)
}
func testAccCheckAlibabacloudStackCloudFirewallControlPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "description" {	
	default = "tf-testAccCloudFirewallControlPolicies-%d"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
	application_name =  "ANY"
	acl_action       =  "accept"
	description      =  var.description
	destination_type =  "net"
	destination      =  "100.1.1.0/24"
	direction        =  "out"
	proto            =  "ANY"
	source           =  "1.2.3.0/24"
	source_type      =  "net"
}

data "alibabacloudstack_cloud_firewall_control_policies" "default" {	
	direction = alibabacloudstack_cloud_firewall_control_policy.default.direction
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
