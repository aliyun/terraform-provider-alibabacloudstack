package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterTablesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_route_table.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_route_table.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_route_table.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_route_table.default.id}"]`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}_fake"`,
		}),
	}

	RouterTablesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStacRouterTablesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouterTablesDatasource%d"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_route_table" "default" {
	transit_router_route_table_description = "${var.name}"
	transit_router_route_table_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}"
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouterTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_tables.#": "1",
		"ids.#":                         "1",
		"transit_router_route_tables.0.transit_router_route_table_name":        CHECKSET,
		"transit_router_route_tables.0.transit_router_route_table_description": CHECKSET,
		"transit_router_route_tables.0.create_time":                            CHECKSET,
		"transit_router_route_tables.0.transit_router_route_table_type":        CHECKSET,
		"transit_router_route_tables.0.status":                                 CHECKSET,
	}
}

var fakeRouterTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_tables.#": "0",
		"ids.#":                         "0",
	}
}

var RouterTablesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_route_tables.default",
	existMapFunc: existRouterTablesMapFunc,
	fakeMapFunc:  fakeRouterTablesMapFunc,
}
