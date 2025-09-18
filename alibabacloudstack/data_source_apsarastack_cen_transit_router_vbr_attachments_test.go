package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterVbrAttachmentsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}"`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_name}"`,
			"ids":               `["${alibabacloudstack_cen_transit_router_vbr_attachment.default.id}"]`,
			"description_regex": `"${alibabacloudstack_cen_transit_router_vbr_attachment.default.transit_router_attachment_description}_fake"`,
		}),
	}

	RouterVbrAttachmentsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStacRouterVbrAttachmentsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouterVbrAttachmentsDatasource%d"
}

%s

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
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

data "alibabacloudstack_cen_transit_router_vbr_attachments" "default" {
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	%s
}`, rand, VSwitchCommonTestCase, strings.Join(pairs, "\n  "))
	return config
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
		"transitrouterattachments.0.charge_type":                           CHECKSET,
		"transitrouterattachments.0.transit_router_id":                     CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_id":          CHECKSET,
		"transitrouterattachments.0.zone_mappings.#":                       "1",
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
