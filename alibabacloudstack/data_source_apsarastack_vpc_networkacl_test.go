package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackVpcNetworkAclsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_network_acls.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_network_acls.default.id}_fake"]`,
		}),
	}

	network_acl_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_network_acls.default.id}"]`,
			"network_acl_id": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_network_acls.default.id}_fake"]`,
			"network_acl_id": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclId}_fake"`,
		}),
	}

	network_acl_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_vpc_network_acls.default.id}"]`,
			"network_acl_name": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_vpc_network_acls.default.id}_fake"]`,
			"network_acl_name": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclName}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_network_acls.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_network_acls.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_network_acls.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_network_acls.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_network_acls.default.id}"]`,

			"network_acl_id":   `"${alibabacloudstack_vpc_network_acls.default.NetworkAclId}"`,
			"network_acl_name": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclName}"`,
			"vpc_id":           `"${alibabacloudstack_vpc_network_acls.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_network_acls.default.id}_fake"]`,

			"network_acl_id":   `"${alibabacloudstack_vpc_network_acls.default.NetworkAclId}_fake"`,
			"network_acl_name": `"${alibabacloudstack_vpc_network_acls.default.NetworkAclName}_fake"`,
			"vpc_id":           `"${alibabacloudstack_vpc_network_acls.default.VpcId}_fake"`}),
	}

	AlibabacloudstackVpcNetworkAclsCheckInfo.dataSourceTestCheck(t, rand, idsConf, network_acl_idConf, network_acl_nameConf, vpc_idConf, allConf)
}

var existAlibabacloudstackVpcNetworkAclsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"acls.#":    "1",
		"acls.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcNetworkAclsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"acls.#": "0",
	}
}

var AlibabacloudstackVpcNetworkAclsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_network_acls.default",
	existMapFunc: existAlibabacloudstackVpcNetworkAclsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcNetworkAclsMapFunc,
}

func testAccCheckAlibabacloudstackVpcNetworkAclsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcNetworkAcls%d"
}






data "alibabacloudstack_vpc_network_acls" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
