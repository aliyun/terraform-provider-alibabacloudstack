package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainMembersDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	AttachmentConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain_member.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain_member.default.id}_fake"]`,
		}),
	}

	SwitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_member.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_member.default.vswitch_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"vswitch_id":                   `"${alibabacloudstack_cen_transit_router_multicast_domain_member.default.vswitch_id}"`,
			"ids":                          `["${alibabacloudstack_cen_transit_router_multicast_domain_member.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"ids":                          `["${alibabacloudstack_cen_transit_router_multicast_domain_member.default.id}"]`,
			"vswitch_id":                   `"${alibabacloudstack_cen_transit_router_multicast_domain_member.default.vswitch_id}_fake"`,
		}),
	}

	MulticastDomainMembersCheckInfo.dataSourceTestCheck(t, rand, AttachmentConf, IdsConf, SwitchConf, allConf)
}

func testAccCheckAlibabacloudStacMulticastDomainMembersDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "datasource_multicast_domain_member%d"
	}

	%s

	%s

	data "alibabacloudstack_instance_types" "all" {
	  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
	  sorted_by         = "CPU"
	  eni_amount        = 2
	}

	resource "alibabacloudstack_ecs_instance" "default" {
	  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	  instance_type        = data.alibabacloudstack_instance_types.all.instance_types.0.id
	  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	  system_disk_size     = 20
	  system_disk_name     = "test_sys_disk"
	  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
	  instance_name        = "${var.name}_ecs"
	  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
	  is_outdated          = false
	  lifecycle {
	    ignore_changes = [
	      instance_type
	    ]
	  }
	}

	  resource "alibabacloudstack_network_interface" "interface" {
		name            = var.name
		vswitch_id      = alibabacloudstack_vpc_vswitch.default.id
		security_groups = [alibabacloudstack_ecs_securitygroup.default.id]

	  }
	  
	  resource "alibabacloudstack_network_interface_attachment" "attachment" {
		instance_id          = alibabacloudstack_ecs_instance.default.id
		network_interface_id = alibabacloudstack_network_interface.interface.id

	  }

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
				 zone_id = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
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

	resource "alibabacloudstack_cen_transit_router_multicast_domain_member" "default" {
		group_ip_address = "239.192.0.2"
		resource_type = "VPC"
		transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}"
		network_interface_id = "${alibabacloudstack_network_interface.interface.id}"
		vswitch_id = "${alibabacloudstack_cen_transit_router_multicast_domain_association.default.vswitch_id}"
	}


data "alibabacloudstack_cen_transit_router_multicast_domain_members" "default" {
	transit_router_multicast_domain_id="${alibabacloudstack_cen_transit_router_multicast_domain_member.default.transit_router_multicast_domain_id}"
	%s
}`, rand, SecurityGroupCommonTestCase, DataAlibabacloudstackImages, strings.Join(pairs, "\n  "))
	return config
}

var existMulticastDomainMembersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_groups.#": "1",
		"ids.#":                             "1",
		"transit_router_multicast_groups.0.group_ip_address":                   CHECKSET,
		"transit_router_multicast_groups.0.network_interface_id":               CHECKSET,
		"transit_router_multicast_groups.0.status":                             CHECKSET,
		"transit_router_multicast_groups.0.transit_router_multicast_domain_id": CHECKSET,
		"transit_router_multicast_groups.0.transit_router_attachment_id":       CHECKSET,
		"transit_router_multicast_groups.0.vswitch_id":                         CHECKSET,
		"transit_router_multicast_groups.0.resource_type":                      CHECKSET,
		"transit_router_multicast_groups.0.member_type":                        CHECKSET,
		"transit_router_multicast_groups.0.resource_id":                        CHECKSET,
		"transit_router_multicast_groups.0.group_source":                       CHECKSET,
		"transit_router_multicast_groups.0.group_member":                       CHECKSET,
	}
}

var fakeMulticastDomainMembersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_groups.#": "0",
		"ids.#":                             "0",
	}
}

var MulticastDomainMembersCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_multicast_domain_members.default",
	existMapFunc: existMulticastDomainMembersMapFunc,
	fakeMapFunc:  fakeMulticastDomainMembersMapFunc,
}
