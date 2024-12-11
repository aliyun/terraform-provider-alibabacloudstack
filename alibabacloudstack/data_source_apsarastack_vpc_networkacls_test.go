package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackVpcNetworkAclsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_network_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_network_acl.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_network_acl.default.network_acl_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_network_acl.default.network_acl_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alibabacloudstack_network_acl.default.id}"]`,
			"status": `"${alibabacloudstack_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alibabacloudstack_network_acl.default.id}"]`,
			"status": `"Modifying"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${alibabacloudstack_network_acl.default.id}"]`,
			"name_regex":       `"${alibabacloudstack_network_acl.default.network_acl_name}"`,
			"network_acl_name": `"${alibabacloudstack_network_acl.default.network_acl_name}"`,
			"status":           `"${alibabacloudstack_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${alibabacloudstack_network_acl.default.id}"]`,
			"name_regex":       `"${alibabacloudstack_network_acl.default.network_acl_name}_fake"`,
			"network_acl_name": `"${alibabacloudstack_network_acl.default.network_acl_name}_fake"`,
			"status":           `"Modifying"`,
		}),
	}
	var existAlibabacloudStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"acls.#":                       "1",
			"acls.0.description":           fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.egress_acl_entries.#":  "1",
			"acls.0.ingress_acl_entries.#": "1",
			"acls.0.network_acl_name":      fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.vpc_id":                CHECKSET,
			"acls.0.status":                "Available",
		}
	}
	var fakeAlibabacloudStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackNetworkAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_network_acls.default",
		existMapFunc: existAlibabacloudStackNetworkAclsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackNetworkAclsDataSourceNameMapFunc,
	}
	alibabacloudstackNetworkAclsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlibabacloudStackNetworkAclsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNetworkAcl-%d"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_network_acl" "default" {
	description = "${var.name}"
	network_acl_name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

data "alibabacloudstack_network_acls" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
