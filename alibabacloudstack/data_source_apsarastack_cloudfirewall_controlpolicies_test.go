package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCloudfirewallControlpoliciesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_cloud_firewall_control_policies.default"
	testDataSourceAttr := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCloudFirewallControlPoliciesMapFunc,
		fakeMapFunc:  fakeCloudFirewallControlPoliciesMapFunc,
		Providers:    testYunDunProviders(),
	}
	testDataSourceAttr.dataSourceTestCheck(t, rand,

		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesDataSourceConfig(rand, map[string]string{
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesDataSourceConfig(rand, map[string]string{
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}_fake"`,
			}),
		},

		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesDataSourceConfig(rand, map[string]string{
				"acl_action": `"${alibabacloudstack_cloud_firewall_control_policy.default.acl_action}"`,
				// "destination": `"${alibabacloudstack_cloud_firewall_control_policy.default.destination}"`,
				"proto": `"${alibabacloudstack_cloud_firewall_control_policy.default.proto}"`,
				// "source":      `"${alibabacloudstack_cloud_firewall_control_policy.default.source}"`,
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesDataSourceConfig(rand, map[string]string{
				"acl_action": `"${alibabacloudstack_cloud_firewall_control_policy.default.acl_action}"`,
				// "destination": `"${alibabacloudstack_cloud_firewall_control_policy.default.destination}"`,
				"proto": `"${alibabacloudstack_cloud_firewall_control_policy.default.proto}"`,
				// "source":      `"${alibabacloudstack_cloud_firewall_control_policy.default.source}"`,
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}_fake"`,
			}),
		},
	)
}

func TestAccAlibabacloudStackCenTransitCloudFirewallControlPoliciesv2DataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_cloud_firewall_control_policies.default"
	testDataSourceAttr := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCloudFirewallControlPoliciesMapFunc,
		fakeMapFunc:  fakeCloudFirewallControlPoliciesMapFunc,
		Providers:    testYunDunProviders(),
	}
	testDataSourceAttr.dataSourceTestCheck(t, rand,

		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"description": `"${alibabacloudstack_cloud_firewall_control_policy.default.description}_fake"`,
			}),
		},

		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"source": `"${alibabacloudstack_cloud_firewall_control_policy.default.source}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"source": `"${alibabacloudstack_cloud_firewall_control_policy.default.source}_fake"`,
			}),
		},

		dataSourceTestAccConfig{
			existConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"destination": `"${alibabacloudstack_cloud_firewall_control_policy.default.destination}"`,
			}),
			fakeConfig: testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand, map[string]string{
				"destination": `"${alibabacloudstack_cloud_firewall_control_policy.default.destination}_fake"`,
			}),
		},
	)
}

func testAccCheckAlibabacloudStacCloudFirewallControlPoliciesv2DataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "description" {	
	default = "tf-testAccCloudFirewallControlPolicies-%d"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
    group_type = "ip"
    group_name = "cloudfw_test_123"
    address_list = ["100.100.100.100/30"]
    description = "test address book"
}

resource "alibabacloudstack_cloudfw_address_book" "port" {
    group_type = "port"
    group_name =  "cloudfw_test_port_123"
    address_list = ["8888","9999"]
    description = "test port book"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
application_name = "ANY"
acl_action       = "accept"
description      = "test-update"
destination_type = "net"
destination      = "114.2.3.0/24"
direction        = "in"
proto            = "ANY"
source           = "192.1.1.0/24"
source_type      = "net"
dest_port  = "8080/8080"
dest_port_type   = "port"
release          = "true"
}

data "alibabacloudstack_cloud_firewall_control_policies" "default" {	
	direction = alibabacloudstack_cloud_firewall_control_policy.default.direction
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

func testAccCheckAlibabacloudStacCloudFirewallControlPoliciesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "description" {	
	default = "tf-testAccCloudFirewallControlPolicies-%d"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
    group_type = "ip"
    group_name = "cloudfw_test_123"
    address_list = ["100.100.100.100/30"]
    description = "test address book"
}

resource "alibabacloudstack_cloudfw_address_book" "port" {
    group_type = "port"
    group_name =  "cloudfw_test_port_123"
    address_list = ["8888","9999"]
    description = "test port book"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
application_name = "ANY"
acl_action       = "accept"
description      = "test-update"
destination_type = "group"
destination      = alibabacloudstack_cloudfw_address_book.default.group_name
direction        = "in"
proto            = "ANY"
source           = alibabacloudstack_cloudfw_address_book.default.group_name
source_type      = "group"
dest_port_group  = alibabacloudstack_cloudfw_address_book.port.group_name
dest_port_type   = "group"
release          = "true"
}

data "alibabacloudstack_cloud_firewall_control_policies" "default" {	
	direction = alibabacloudstack_cloud_firewall_control_policy.default.direction
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

var existCloudFirewallControlPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                       "1",
		"policies.#":                  "1",
		"policies.0.description":      CHECKSET,
		"policies.0.application_name": CHECKSET,
		"policies.0.acl_action":       CHECKSET,
		"policies.0.destination_type": CHECKSET,
		"policies.0.destination":      CHECKSET,
		"policies.0.direction":        CHECKSET,
		"policies.0.proto":            CHECKSET,
		"policies.0.source":           CHECKSET,
		"policies.0.source_type":      CHECKSET,
	}
}

var fakeCloudFirewallControlPoliciesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
	}
}
