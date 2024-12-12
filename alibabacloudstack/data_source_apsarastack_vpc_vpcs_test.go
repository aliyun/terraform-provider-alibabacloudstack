package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackVpcVpcsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
		}),
	}

	dhcp_options_set_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
			"dhcp_options_set_id": `"${alibabacloudstack_vpc_vpcs.default.DhcpOptionsSetId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
			"dhcp_options_set_id": `"${alibabacloudstack_vpc_vpcs.default.DhcpOptionsSetId}_fake"`,
		}),
	}

	is_defaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_vpcs.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_vpcs.default.ResourceGroupId}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vpcs.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vpcs.default.VpcId}_fake"`,
		}),
	}

	vpc_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_vpc_vpcs.default.id}"]`,
			"vpc_name": `"${alibabacloudstack_vpc_vpcs.default.VpcName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,
			"vpc_name": `"${alibabacloudstack_vpc_vpcs.default.VpcName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}"]`,

			"dhcp_options_set_id": `"${alibabacloudstack_vpc_vpcs.default.DhcpOptionsSetId}"`,
			"resource_group_id":   `"${alibabacloudstack_vpc_vpcs.default.ResourceGroupId}"`,
			"vpc_id":              `"${alibabacloudstack_vpc_vpcs.default.VpcId}"`,
			"vpc_name":            `"${alibabacloudstack_vpc_vpcs.default.VpcName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpcs.default.id}_fake"]`,

			"dhcp_options_set_id": `"${alibabacloudstack_vpc_vpcs.default.DhcpOptionsSetId}_fake"`,
			"resource_group_id":   `"${alibabacloudstack_vpc_vpcs.default.ResourceGroupId}_fake"`,
			"vpc_id":              `"${alibabacloudstack_vpc_vpcs.default.VpcId}_fake"`,
			"vpc_name":            `"${alibabacloudstack_vpc_vpcs.default.VpcName}_fake"`}),
	}

	AlibabacloudstackVpcVpcsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, dhcp_options_set_idConf, is_defaultConf, resource_group_idConf, vpc_idConf, vpc_nameConf, allConf)
}

var existAlibabacloudstackVpcVpcsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vpcs.#":    "1",
		"vpcs.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcVpcsDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vpcs.#": "0",
	}
}

var AlibabacloudstackVpcVpcsDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_vpcs.default",
	existMapFunc: existAlibabacloudstackVpcVpcsDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcVpcsDataMapFunc,
}

func testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcVpcs%d"
}






data "alibabacloudstack_vpc_vpcs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
