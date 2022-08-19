package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackVpcNetworkAclsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_network_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_network_acl.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_network_acl.default.network_acl_name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${apsarastack_network_acl.default.network_acl_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${apsarastack_network_acl.default.id}"]`,
			"status": `"${apsarastack_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${apsarastack_network_acl.default.id}"]`,
			"status": `"Modifying"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${apsarastack_network_acl.default.id}"]`,
			"name_regex":       `"${apsarastack_network_acl.default.network_acl_name}"`,
			"network_acl_name": `"${apsarastack_network_acl.default.network_acl_name}"`,
			"status":           `"${apsarastack_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckApsaraStackNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${apsarastack_network_acl.default.id}"]`,
			"name_regex":       `"${apsarastack_network_acl.default.network_acl_name}_fake"`,
			"network_acl_name": `"${apsarastack_network_acl.default.network_acl_name}_fake"`,
			"status":           `"Modifying"`,
		}),
	}
	var existApsaraStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
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
	var fakeApsaraStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var apsarastackNetworkAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_network_acls.default",
		existMapFunc: existApsaraStackNetworkAclsDataSourceNameMapFunc,
		fakeMapFunc:  fakeApsaraStackNetworkAclsDataSourceNameMapFunc,
	}
	apsarastackNetworkAclsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckApsaraStackNetworkAclsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNetworkAcl-%d"
}

resource "apsarastack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "apsarastack_network_acl" "default" {
	description = "${var.name}"
	network_acl_name = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}

data "apsarastack_network_acls" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
