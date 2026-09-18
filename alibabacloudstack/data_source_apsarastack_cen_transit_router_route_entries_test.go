package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterEntriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_route_entries.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouterEntriesDatasource%d", rand),
		dataSourceCenTransitRouterRouteEntriesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"name_regex":                    "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"name_regex":                    "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_entry.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_entry.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"description_regex":             "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"description_regex":             "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"name_regex":                    "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}",
			"description_regex":             "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_entry.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
			"name_regex":                    "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_name}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_entry.default.id}"},
			"description_regex":             "${alibabacloudstack_cen_transit_router_route_entry.default.transit_router_route_entry_description}_fake",
		}),
	}

	RouterEntriesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceCenTransitRouterRouteEntriesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
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
