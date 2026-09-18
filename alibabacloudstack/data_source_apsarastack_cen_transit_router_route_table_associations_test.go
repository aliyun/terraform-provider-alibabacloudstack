package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterTableAssociationsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_route_table_associations.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouterTableAssociationsDatasource%d", rand),
		dataSourceCenTransitRouterRouteTableAssociationsConfigDependence)

	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table_association.default.transit_router_route_table_id}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_table_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_route_table_id": "${alibabacloudstack_cen_transit_router_route_table_association.default.transit_router_route_table_id}",
			"ids":                           []string{"${alibabacloudstack_cen_transit_router_route_table_association.default.id}_fake"},
		}),
	}

	RouterTableAssociationsCheckInfo.dataSourceTestCheck(t, rand, IdsConf)
}

func dataSourceCenTransitRouterRouteTableAssociationsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
	transit_router_attachment_name = "${var.name}"
	transit_router_attachment_description = "${var.name}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	zone_mappings {
			vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
			 zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
		}
}

data "alibabacloudstack_cen_transit_router_route_tables" "default" {
	transit_router_route_id="${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_route_table_association" "default" {
	transit_router_route_table_id = "${data.alibabacloudstack_cen_transit_router_route_tables.default.transit_router_route_tables.0.id}"
	transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
}
`, name, VSwitchCommonTestCase)
}

var existRouterTableAssociationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_associations.#": "1",
		"ids.#":                               "1",
		"transit_router_route_associations.0.transit_router_attachment_id":  CHECKSET,
		"transit_router_route_associations.0.transit_router_route_table_id": CHECKSET,
		"transit_router_route_associations.0.resource_id":                   CHECKSET,
		"transit_router_route_associations.0.resource_type":                 CHECKSET,
		"transit_router_route_associations.0.status":                        CHECKSET,
	}
}

var fakeRouterTableAssociationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_associations.#": "0",
		"ids.#":                               "0",
	}
}

var RouterTableAssociationsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_route_table_associations.default",
	existMapFunc: existRouterTableAssociationsMapFunc,
	fakeMapFunc:  fakeRouterTableAssociationsMapFunc,
}
