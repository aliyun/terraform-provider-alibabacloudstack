package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackVpcVswitchesDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
		}),
	}

	is_defaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
		}),
	}

	route_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}_fake"`,
		}),
	}

	vswitch_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vswitch_name": `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vswitch_name": `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vswitches.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vswitches.default.VpcId}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"zone_id": `"${alibabacloudstack_vpc_vswitches.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_vpc_vswitches.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,

			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}"`,
			"vswitch_name":   `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vswitches.default.VpcId}"`,
			"zone_id":        `"${alibabacloudstack_vpc_vswitches.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,

			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}_fake"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}_fake"`,
			"vswitch_name":   `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}_fake"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vswitches.default.VpcId}_fake"`,
			"zone_id":        `"${alibabacloudstack_vpc_vswitches.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackVpcVswitchesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, is_defaultConf, route_table_idConf, vswitch_idConf, vswitch_nameConf, vpc_idConf, zone_idConf, allConf)
}

var existAlibabacloudstackVpcVswitchesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vswitches.#":    "1",
		"vswitches.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcVswitchesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vswitches.#": "0",
	}
}

var AlibabacloudstackVpcVswitchesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_vswitches.default",
	existMapFunc: existAlibabacloudstackVpcVswitchesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcVswitchesDataMapFunc,
}

func testAccCheckAlibabacloudstackVpcVswitchesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcVswitches%d"
}






data "alibabacloudstack_vpc_vswitches" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
