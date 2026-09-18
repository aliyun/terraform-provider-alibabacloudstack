package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainAssociationsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_multicast_domain_associations.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMulticastDomainAssociationsDatasource%d", rand),
		dataSourceCenMulticastDomainAssociationsConfigDependence)

	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"ids": []string{"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"ids": []string{"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.id}_fake"},
		}),
	}

	vswitchConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"vswitch_id_regex": "${alibabacloudstack_vpc_vswitch.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"vswitch_id_regex": "${alibabacloudstack_vpc_vswitch.default.id}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"vswitch_id_regex": "${alibabacloudstack_vpc_vswitch.default.id}",
			"ids":              []string{"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
			"ids":              []string{"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.id}"},
			"vswitch_id_regex": "${alibabacloudstack_vpc_vswitch.default.id}_fake",
		}),
	}

	MulticastDomainAssociationsCheckInfo.dataSourceTestCheck(t, rand, IdsConf, vswitchConf, allConf)
}

func dataSourceCenMulticastDomainAssociationsConfigDependence(name string) string {
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

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
	transit_router_multicast_domain_description = "${var.name}"
	transit_router_multicast_domain_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
resource "alibabacloudstack_cen_transit_router_multicast_domain_association" "default" {
	transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}"
	transit_router_attachment_id = "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
`, name, VSwitchCommonTestCase)
}

var existMulticastDomainAssociationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_associations.#": "1",
		"ids.#": "1",
		"transit_router_multicast_associations.0.transit_router_attachment_id":       CHECKSET,
		"transit_router_multicast_associations.0.transit_router_multicast_domain_id": CHECKSET,
		"transit_router_multicast_associations.0.vswitch_id":                         CHECKSET,
		"transit_router_multicast_associations.0.status":                             CHECKSET,
	}
}

var fakeMulticastDomainAssociationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_associations.#": "0",
		"ids.#": "0",
	}
}

var MulticastDomainAssociationsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_multicast_domain_associations.default",
	existMapFunc: existMulticastDomainAssociationsMapFunc,
	fakeMapFunc:  fakeMulticastDomainAssociationsMapFunc,
}
