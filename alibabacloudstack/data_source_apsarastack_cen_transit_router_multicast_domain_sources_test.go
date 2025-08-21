package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainSourcesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	AttachmentConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain_source.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_multicast_domain_source.default.id}_fake"]`,
		}),
	}

	SwitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_source.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_source.default.vswitch_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"vswitch_id":                   `"${alibabacloudstack_cen_transit_router_multicast_domain_source.default.vswitch_id}"`,
			"ids":                          `["${alibabacloudstack_cen_transit_router_multicast_domain_source.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand, map[string]string{
			"transit_router_attachment_id": `"${alibabacloudstack_cen_transit_router_multicast_domain_association.default.transit_router_attachment_id}"`,
			"ids":                          `["${alibabacloudstack_cen_transit_router_multicast_domain_source.default.id}"]`,
			"vswitch_id":                   `"${alibabacloudstack_cen_transit_router_multicast_domain_source.default.vswitch_id}_fake"`,
		}),
	}

	MulticastDomainSourcesCheckInfo.dataSourceTestCheck(t, rand, AttachmentConf, IdsConf, SwitchConf, allConf)
}

func testAccCheckAlibabacloudStacMulticastDomainSourcesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := DataAlibabacloudstackVswitchZones + DataAlibabacloudstackImages + fmt.Sprintf(`
	variable "name" {
		default = "datasource_multicast_domain_source%d"
	}
	resource "alibabacloudstack_vpc" "vpc" {
		name       = var.name
		cidr_block = "192.168.0.0/24"
	
	  }
	  
	  resource "alibabacloudstack_vswitch" "vswitch" {
		name              = var.name
		cidr_block        = "192.168.0.0/24"
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		vpc_id            = alibabacloudstack_vpc.vpc.id
	
	  }
	  
	  resource "alibabacloudstack_security_group" "group" {
		name   = var.name
		vpc_id = alibabacloudstack_vpc.vpc.id
	
	  }
	
	data "alibabacloudstack_instance_types" "instance_type" {
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		eni_amount        = 2
		sorted_by         = "Memory"
	  }
	  
	  resource "alibabacloudstack_instance" "instance" {
		availability_zone = data.alibabacloudstack_zones.default.zones[0].id
		security_groups   = [alibabacloudstack_security_group.group.id]
		instance_type              = data.alibabacloudstack_instance_types.instance_type.instance_types[0].id
		system_disk_category       = "cloud_efficiency"
		image_id                   = data.alibabacloudstack_images.default.images[0].id
		instance_name              = var.name
		vswitch_id                 = alibabacloudstack_vswitch.vswitch.id
	
	  
	  }
	
	  resource "alibabacloudstack_network_interface" "interface" {
		name            = var.name
		vswitch_id      = alibabacloudstack_vswitch.vswitch.id
		security_groups = [alibabacloudstack_security_group.group.id]
	
	  }
	  
	  resource "alibabacloudstack_network_interface_attachment" "attachment" {
		instance_id          = alibabacloudstack_instance.instance.id
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
		vpc_id = "${alibabacloudstack_vswitch.vswitch.vpc_id}"
		transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
		zone_mappings {
				vswitch_id = "${alibabacloudstack_vswitch.vswitch.id}"
				 zone_id = "${alibabacloudstack_vswitch.vswitch.availability_zone}"
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
		vswitch_id = "${alibabacloudstack_vswitch.vswitch.id}"
	
	}

	resource "alibabacloudstack_cen_transit_router_multicast_domain_source" "default" {
		group_ip_address = "239.192.0.2"
		transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}"
		network_interface_id = "${alibabacloudstack_network_interface.interface.id}"
		vswitch_id = "${alibabacloudstack_cen_transit_router_multicast_domain_association.default.vswitch_id}"
	}


data "alibabacloudstack_cen_transit_router_multicast_domain_sources" "default" {
	transit_router_multicast_domain_id="${alibabacloudstack_cen_transit_router_multicast_domain_source.default.transit_router_multicast_domain_id}"
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existMulticastDomainSourcesMapFunc = func(rand int) map[string]string {
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
		"transit_router_multicast_groups.0.source_type":                        CHECKSET,
		"transit_router_multicast_groups.0.resource_id":                        CHECKSET,
		"transit_router_multicast_groups.0.group_source":                       CHECKSET,
		"transit_router_multicast_groups.0.group_member":                       CHECKSET,
	}
}

var fakeMulticastDomainSourcesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_multicast_groups.#": "0",
		"ids.#":                             "0",
	}
}

var MulticastDomainSourcesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_multicast_domain_sources.default",
	existMapFunc: existMulticastDomainSourcesMapFunc,
	fakeMapFunc:  fakeMulticastDomainSourcesMapFunc,
}
