package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackVpcVSwitchesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
		}),
	}

	is_defaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
		}),
	}

	route_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}_fake"`,
		}),
	}

	vswitch_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vswitch_name": `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vswitch_name": `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vswitches.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_vswitches.default.VpcId}_fake"`,
		}),
	}

	zone_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_vpc_vswitches.default.id}"]`,
			"zone_id": `"${alibabacloudstack_vpc_vswitches.default.ZoneId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids":     `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,
			"zone_id": `"${alibabacloudstack_vpc_vswitches.default.ZoneId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}"]`,

			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}"`,
			"vswitch_name":   `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vswitches.default.VpcId}"`,
			"zone_id":        `"${alibabacloudstack_vpc_vswitches.default.ZoneId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_vswitches.default.id}_fake"]`,

			"route_table_id": `"${alibabacloudstack_vpc_vswitches.default.RouteTableId}_fake"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_vswitches.default.VSwitchId}_fake"`,
			"vswitch_name":   `"${alibabacloudstack_vpc_vswitches.default.VSwitchName}_fake"`,
			"vpc_id":         `"${alibabacloudstack_vpc_vswitches.default.VpcId}_fake"`,
			"zone_id":        `"${alibabacloudstack_vpc_vswitches.default.ZoneId}_fake"`}),
	}

	AlibabacloudstackVpcVSwitchesCheckInfo.dataSourceTestCheck(t, rand, idsConf, is_defaultConf, route_table_idConf, vswitch_idConf, vswitch_nameConf, vpc_idConf, zone_idConf, allConf)
}

var existAlibabacloudstackVpcVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"switches.#":    "1",
		"switches.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"switches.#": "0",
	}
}

var AlibabacloudstackVpcVSwitchesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_vswitches.default",
	existMapFunc: existAlibabacloudstackVpcVSwitchesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcVSwitchesMapFunc,
}

func testAccCheckAlibabacloudstackVpcVSwitchesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcVSwitches%d"
}






data "alibabacloudstack_vpc_vswitches" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
