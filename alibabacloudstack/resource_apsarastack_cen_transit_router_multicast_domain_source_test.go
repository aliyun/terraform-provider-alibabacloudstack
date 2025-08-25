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
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
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
`, name, SecurityGroupCommonTestCase , DataAlibabacloudstackImages )
}
