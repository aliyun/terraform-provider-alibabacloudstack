package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackVpcRouteTablesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_route_tables.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_route_tables.default.ResourceGroupId}_fake"`,
		}),
	}

	route_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"route_table_id": `"${alibabacloudstack_vpc_route_tables.default.RouteTableId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"route_table_id": `"${alibabacloudstack_vpc_route_tables.default.RouteTableId}_fake"`,
		}),
	}

	route_table_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"route_table_name": `"${alibabacloudstack_vpc_route_tables.default.RouteTableName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"route_table_name": `"${alibabacloudstack_vpc_route_tables.default.RouteTableName}_fake"`,
		}),
	}

	router_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":       `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"router_id": `"${alibabacloudstack_vpc_route_tables.default.RouterId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":       `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"router_id": `"${alibabacloudstack_vpc_route_tables.default.RouterId}_fake"`,
		}),
	}

	router_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"router_type": `"${alibabacloudstack_vpc_route_tables.default.RouterType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"router_type": `"${alibabacloudstack_vpc_route_tables.default.RouterType}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_route_tables.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_route_tables.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_route_tables.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_tables.default.id}"]`,

			"resource_group_id": `"${alibabacloudstack_vpc_route_tables.default.ResourceGroupId}"`,
			"route_table_id":    `"${alibabacloudstack_vpc_route_tables.default.RouteTableId}"`,
			"route_table_name":  `"${alibabacloudstack_vpc_route_tables.default.RouteTableName}"`,
			"router_id":         `"${alibabacloudstack_vpc_route_tables.default.RouterId}"`,
			"router_type":       `"${alibabacloudstack_vpc_route_tables.default.RouterType}"`,
			"vpc_id":            `"${alibabacloudstack_vpc_route_tables.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_tables.default.id}_fake"]`,

			"resource_group_id": `"${alibabacloudstack_vpc_route_tables.default.ResourceGroupId}_fake"`,
			"route_table_id":    `"${alibabacloudstack_vpc_route_tables.default.RouteTableId}_fake"`,
			"route_table_name":  `"${alibabacloudstack_vpc_route_tables.default.RouteTableName}_fake"`,
			"router_id":         `"${alibabacloudstack_vpc_route_tables.default.RouterId}_fake"`,
			"router_type":       `"${alibabacloudstack_vpc_route_tables.default.RouterType}_fake"`,
			"vpc_id":            `"${alibabacloudstack_vpc_route_tables.default.VpcId}_fake"`}),
	}

	AlibabacloudstackVpcRouteTablesCheckInfo.dataSourceTestCheck(t, rand, idsConf, resource_group_idConf, route_table_idConf, route_table_nameConf, router_idConf, router_typeConf, vpc_idConf, allConf)
}

var existAlibabacloudstackVpcRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"tables.#":    "1",
		"tables.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"tables.#": "0",
	}
}

var AlibabacloudstackVpcRouteTablesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_route_tables.default",
	existMapFunc: existAlibabacloudstackVpcRouteTablesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcRouteTablesMapFunc,
}

func testAccCheckAlibabacloudstackVpcRouteTablesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcRouteTables%d"
}






data "alibabacloudstack_vpc_route_tables" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
