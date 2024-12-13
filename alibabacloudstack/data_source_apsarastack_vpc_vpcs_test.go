package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackVpcVpcsDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_vpc.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_vpc.default.resource_group_id}_fake"`,
		}),
	}

	vswtich_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id":    `"${alibabacloudstack_vpc_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id":    `"${alibabacloudstack_vpc_vswitch.default.id}_fake"`,
		}),
	}

	vpc_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"vpc_name": `"${alibabacloudstack_vpc_vpc.default.vpc_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"vpc_name": `"${alibabacloudstack_vpc_vpc.default.vpc_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}"]`,

			"resource_group_id":   `"${alibabacloudstack_vpc_vpc.default.resource_group_id}"`,
			"vpc_name":            `"${alibabacloudstack_vpc_vpc.default.vpc_name}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,

			"resource_group_id":   `"${alibabacloudstack_vpc_vpc.default.resource_group_id}_fake"`,
			"vpc_name":            `"${alibabacloudstack_vpc_vpc.default.vpc_name}_fake"`}),
	}

	AlibabacloudstackVpcVpcsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, vswtich_idConf, resource_group_idConf, vpc_nameConf, allConf)
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
%s

data "alibabacloudstack_vpc_vpcs" "default" {
%s
}
`, rand, VSwitchCommonTestCase, strings.Join(pairs, "\n   "))
	return config
}
