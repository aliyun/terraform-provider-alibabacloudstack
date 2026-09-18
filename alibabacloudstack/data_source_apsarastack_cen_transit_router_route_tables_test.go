package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterTablesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_route_tables.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouterTablesDatasource%d", rand),
		dataSourceCenTransitRouterRouteTablesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_route_table.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_route_table.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"name_regex":          "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}",
			"description_regex":   "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}",
			"ids":                 []string{"${alibabacloudstack_cen_transit_router_route_table.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_id": "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_id}",
			"name_regex":          "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_name}",
			"ids":                 []string{"${alibabacloudstack_cen_transit_router_route_table.default.id}"},
			"description_regex":   "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_description}_fake",
		}),
	}

	RouterTablesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceCenTransitRouterRouteTablesConfigDependence(name string) string {
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
`, name)
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
