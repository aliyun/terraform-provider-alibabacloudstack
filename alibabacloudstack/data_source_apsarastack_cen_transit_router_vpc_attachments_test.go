package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterVpcAttachmentsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_vpc_attachments.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouterVpcAttachmentsDatasource%d", rand),
		dataSourceCenTransitRouterVpcAttachmentsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vpc_attachment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vpc_attachment.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}",
			"description_regex": "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_description}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vpc_attachment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
			"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
			"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
			"name_regex":        "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_name}",
			"ids":               []string{"${alibabacloudstack_cen_transit_router_vpc_attachment.default.id}"},
			"description_regex": "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_description}_fake",
		}),
	}

	RouterVpcAttachmentsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceCenTransitRouterVpcAttachmentsConfigDependence(name string) string {
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
`, name, VSwitchCommonTestCase)
}

var existRouterVpcAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transitrouterattachments.#": "1",
		"ids.#":                      "1",
		"transitrouterattachments.0.creation_time":                         CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_name":        CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_description": CHECKSET,
		"transitrouterattachments.0.resource_type":                         CHECKSET,
		"transitrouterattachments.0.status":                                CHECKSET,
		"transitrouterattachments.0.auto_publish_route_enabled":            CHECKSET,
		"transitrouterattachments.0.charge_type":                           CHECKSET,
		"transitrouterattachments.0.transit_router_id":                     CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_id":          CHECKSET,
		"transitrouterattachments.0.zone_mappings.#":                       "1",
	}
}

var fakeRouterVpcAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_tables.#": "0",
		"ids.#":                         "0",
	}
}

var RouterVpcAttachmentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_vpc_attachments.default",
	existMapFunc: existRouterVpcAttachmentsMapFunc,
	fakeMapFunc:  fakeRouterVpcAttachmentsMapFunc,
}
