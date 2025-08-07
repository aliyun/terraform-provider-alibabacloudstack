package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterEntriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_route_entry.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_route_entry.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}"`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_route_entry.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_route_entry.default.id}"]`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}_fake"`,
		}),
	}

	RouterEntriesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStacRouterEntriesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouterEntriesDatasource%d"
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

resource "alibabacloudstack_cen_transit_router_route_entry" "default" {
	transit_router_route_entry_description = "${var.name}"

	transit_router_route_entry_destination_cidr_block = "10.10.10.1/32"

	transit_router_route_entry_name = "${var.name}"
	transit_router_route_entry_next_hop_type = "BlackHole"
	transit_router_route_table_id = "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
}

data "alibabacloudstack_cen_transit_router_route_entries" "default" {
	transit_router_route_table_id="${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}"
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existRouterEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_entries.#": "1",
		"ids.#":                          "1",
		"transit_router_route_entries.0.transit_router_route_entry_name":                   CHECKSET,
		"transit_router_route_entries.0.transit_router_route_entry_description":            CHECKSET,
		"transit_router_route_entries.0.create_time":                                       CHECKSET,
		"transit_router_route_entries.0.transit_router_route_entry_type":                   CHECKSET,
		"transit_router_route_entries.0.transit_router_route_entry_status":                 CHECKSET,
		"transit_router_route_entries.0.transit_router_route_entry_next_hop_type":          CHECKSET,
		"transit_router_route_entries.0.transit_router_route_entry_destination_cidr_block": CHECKSET,
		"transit_router_route_entries.0.operational_mode":                                  CHECKSET,
	}
}

var fakeRouterEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_entries.#": "0",
		"ids.#":                          "0",
	}
}

var RouterEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_route_entries.default",
	existMapFunc: existRouterEntriesMapFunc,
	fakeMapFunc:  fakeRouterEntriesMapFunc,
}
