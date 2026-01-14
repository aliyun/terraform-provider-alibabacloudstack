package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackVpcVpcsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
		}),
	}

	vswtich_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_vpc_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_vpc_vswitch.default.id}_fake"`,
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
	name_regexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"name_regex": `"${alibabacloudstack_vpc_vpc.default.vpc_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_vpc_vpc.default.vpc_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"status": `"Pending"`,
		}),
	}

	isDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"is_default": `false`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"is_default": `true`,
		}),
	}

	cidr_blockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"cidr_block": `"${alibabacloudstack_vpc_vpc.default.cidr_block}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"cidr_block": `"198.1.0.0/16"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"tags": `{"common_test": "terraform"}`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"tags": `{"common_test": "terraform_fake"}`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitch.default.id}"`,
			"vpc_name":   `"${alibabacloudstack_vpc_vpc.default.vpc_name}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVpcsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vpc.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitch.default.id}_fake"`,
			"vpc_name":   `"${alibabacloudstack_vpc_vpc.default.vpc_name}_fake"`}),
	}

	AlibabacloudstackVpcVpcsDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, vswtich_idConf, vpc_nameConf, name_regexConf, statusConf, isDefaultConf, cidr_blockConf, tagsConf, allConf)
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
