package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterVbrAttachmentsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_cen_transit_router_vbr_attachments.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouterVbrAttachmentsDatasource%d", rand),
		dataSourceCenTransitRouterVbrAttachmentsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}"},
			"description_regex": "${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}_fake",
		}),
	}

	RouterVbrAttachmentsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceCenTransitRouterVbrAttachmentsConfigDependence(name string) string {
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

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id = "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	enable_ipv6              = true
	local_ipv6_gateway_ip = "2408:4004:cc:400::1"
	peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
	peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
}
resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	transit_router_attachment_name = "${var.name}"
	transit_router_attachment_description = "${var.name}"
}
`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), getAccTestRandInt(1000, 2000))
}

var existRouterVbrAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transitrouterattachments.#": "1",
		"ids.#":                      "1",
		"transitrouterattachments.0.creation_time":                         CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_name":        CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_description": CHECKSET,
		"transitrouterattachments.0.resource_type":                         CHECKSET,
		"transitrouterattachments.0.status":                                CHECKSET,
		"transitrouterattachments.0.auto_publish_route_enabled":            CHECKSET,
		"transitrouterattachments.0.transit_router_id":                     CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_id":          CHECKSET,
	}
}

var fakeRouterVbrAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_tables.#": "0",
		"ids.#":                         "0",
	}
}

var RouterVbrAttachmentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_vbr_attachments.default",
	existMapFunc: existRouterVbrAttachmentsMapFunc,
	fakeMapFunc:  fakeRouterVbrAttachmentsMapFunc,
}
