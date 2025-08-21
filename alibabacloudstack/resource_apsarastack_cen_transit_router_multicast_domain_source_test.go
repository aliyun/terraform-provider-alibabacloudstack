package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainSource0(t *testing.T) {
	var v *CbnDescribeTransitRouterMulticastDomainSourceResponse

	resourceId := "alibabacloudstack_cen_transit_router_multicast_domain_source.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitMulticastDomainSourceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterMuliticastDomainSourceRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smulticast_domain_source%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitMulticastDomainSourceBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"group_ip_address":                   "239.192.0.2",
					"vswitch_id":                         "${alibabacloudstack_cen_transit_router_multicast_domain_association.default.vswitch_id}",
					"network_interface_id":               "${alibabacloudstack_network_interface.interface.id}",
					"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"group_ip_address":                   "239.192.0.2",
						"vswitch_id":                         CHECKSET,
						"network_interface_id":               CHECKSET,
						"transit_router_multicast_domain_id": CHECKSET,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccCenTransitMulticastDomainSourceCheckmap = map[string]string{

	"status": CHECKSET,

	"group_ip_address": CHECKSET,

	"vswitch_id":                         CHECKSET,
	"network_interface_id":               CHECKSET,
	"transit_router_multicast_domain_id": CHECKSET,
}

func AlibabacloudTestAccCenTransitMulticastDomainSourceBasicdependence(name string) string {
	return DataAlibabacloudstackVswitchZones + DataAlibabacloudstackImages + fmt.Sprintf(
		`
variable "name" {
	default = "%s"
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
`, name)
}
