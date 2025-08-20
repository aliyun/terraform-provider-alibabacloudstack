package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenRouteMapsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_route_map.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_route_map.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_route_map.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_route_map.default.description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_route_map.default.description}"`,
			"ids":               `["${alibabacloudstack_cen_route_map.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_cen_route_map.default.id}"]`,
			"description_regex": `"${alibabacloudstack_cen_route_map.default.description}_fake"`,
		}),
	}

	RouteMapsCheckInfo.dataSourceTestCheck(t, rand, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStacRouteMapsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouteMapsDatasource%d"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}


resource "alibabacloudstack_cen_route_map" "default" {

	cen_id = "${alibabacloudstack_cen_instance.default.cen_id}"
	transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
	priority = "3"
	transmit_direction = "RegionIn"
	map_result = "Deny"
	description = "tf_testAccCenRouteMap"
}

data "alibabacloudstack_cen_route_maps" "default" {
	cen_id="${alibabacloudstack_cen_instance.default.cen_id}"
	transit_router_route_table_id="${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouteMapsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"route_maps.#":                               "1",
		"ids.#":                                      "1",
		"route_maps.0.route_map_id":                  CHECKSET,
		"route_maps.0.cen_id":                        CHECKSET,
		"route_maps.0.map_result":                    CHECKSET,
		"route_maps.0.transmit_direction":            CHECKSET,
		"route_maps.0.status":                        CHECKSET,
		"route_maps.0.transit_router_route_table_id": CHECKSET,
		"route_maps.0.description":                   CHECKSET,
	}
}

var fakeRouteMapsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"route_maps.#": "0",
		"ids.#":        "0",
	}
}

var RouteMapsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_route_maps.default",
	existMapFunc: existRouteMapsMapFunc,
	fakeMapFunc:  fakeRouteMapsMapFunc,
}
