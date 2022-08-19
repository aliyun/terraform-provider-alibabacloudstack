package apsarastack

import (
	"fmt"
	"testing"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsarastackRouteEntryInstance(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "apsarastack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_instance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "Instance",
						"name":         fmt.Sprintf("tf-testAccRouteEntryConfigName%d", rand),
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

func TestAccApsarastackRouteEntryInterface(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "apsarastack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_interface(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "RouterInterface",
						"name":         fmt.Sprintf("tf-testAccRouteEntryInterfaceConfig%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccApsarastackRouteEntryNatGateway(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "apsarastack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_natGateway(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "NatGateway",
						"name":         fmt.Sprintf("tf-testAccRouteEntryNatGatewayConfig%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccApsarastackRouteEntryMulti(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "apsarastack_route_entry.default.2"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":          "NetworkInterface",
						"destination_cidrblock": "172.16.4.0/24",
						"name":                  fmt.Sprintf("tf-testAccRouteEntryConcurrence%d", rand),
					}),
				),
			},
		},
	})
}

func testAccCheckRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_route_entry" || rs.Type != "apsarastack_route_entry" {
			continue
		}
		entry, err := vpcService.DescribeRouteEntry(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if entry.RouteTableId != "" {
			return WrapError(Error("Route entry still exist"))
		}
	}

	//	testAccCheckRouterInterfaceDestroy(s)

	return nil
}

func testAccRouteEntryConfig_instance(rand int) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
   default = "tf-testAccRouteEntryConfigName%d"
}
resource "apsarastack_vpc" "default" {
   name = "${var.name}"
   cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
   vpc_id = "${apsarastack_vpc.default.id}"
   cidr_block = "10.1.1.0/24"
   availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
   name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
   name = "${var.name}"
   description = "default"
   vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
   type = "ingress"
   ip_protocol = "tcp"
   nic_type = "intranet"
   policy = "accept"
   port_range = "22/22"
   priority = 1
   security_group_id = "${apsarastack_security_group.default.id}"
   cidr_ip = "0.0.0.0/0"
}
resource "apsarastack_instance" "default" {
   security_groups = ["${apsarastack_security_group.default.id}"]
   vswitch_id = "${apsarastack_vswitch.default.id}"
   instance_type = "${local.instance_type_id}"
   internet_max_bandwidth_out = 5
   system_disk_category = "cloud_efficiency"
   image_id = "${data.apsarastack_images.default.images.0.id}"
   instance_name = "${var.name}"
}
resource "apsarastack_route_entry" "default" {
   route_table_id = "${apsarastack_vpc.default.route_table_id}"
   destination_cidrblock = "172.11.1.1/32"
   nexthop_type = "Instance"
   nexthop_id = "${apsarastack_instance.default.id}"
   name = "${var.name}"
}
`, DataApsarastackVswitchZones, DataApsarastackInstanceTypes, DataApsarastackImages, rand)
}

func testAccRouteEntryConfig_interface(rand int) string {
	return fmt.Sprintf(
		`
data "apsarastack_zones" "default" {
  available_resource_creation= "VSwitch"
}
variable "name" {
   default = "tf-testAccRouteEntryInterfaceConfig%d"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  name = "${var.name}"
  description = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${apsarastack_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "apsarastack_router_interface" "default" {
  opposite_region = "cn-qingdao-env66-d01"
  router_type = "VRouter"
  router_id = "${apsarastack_vpc.default.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "${var.name}"
}
resource "apsarastack_route_entry" "default" {
  route_table_id = "${apsarastack_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "RouterInterface"
  nexthop_id = "${apsarastack_router_interface.default.id}"
  name = "${var.name}"
}
`, rand)
}

func testAccRouteEntryConfig_natGateway(rand int) string {
	return fmt.Sprintf(
		`
data "apsarastack_zones" "default" {
  available_resource_creation= "VSwitch"
}
variable "name" {
   default = "tf-testAccRouteEntryNatGatewayConfig%d"
}
resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
  vpc_id = "${apsarastack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "apsarastack_nat_gateway" "default" {
   vpc_id = "${apsarastack_vpc.default.id}"
   specification = "Middle"
   name = "${var.name}"
}
resource "apsarastack_route_entry" "default" {
  route_table_id = "${apsarastack_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "NatGateway"
  nexthop_id = "${apsarastack_nat_gateway.default.id}"
  name = "${var.name}"
}`, rand)
}

func testAccRouteEntryConfigMulti(rand int) string {
	return fmt.Sprintf(`
data "apsarastack_zones" "default" {
   available_resource_creation= "VSwitch"
}
variable "name" {
   default = "tf-testAccRouteEntryConcurrence%d"
}
resource "apsarastack_vpc" "default" {
   name = "${var.name}"
   cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "10.1.1.0/24"
    availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
    vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
}
resource "apsarastack_route_entry" "default" {
   count = 3
   route_table_id = "${apsarastack_vpc.default.route_table_id}"
   destination_cidrblock = "172.16.${count.index}.0/24"
   nexthop_type = "NetworkInterface"
   nexthop_id = "${apsarastack_network_interface.default.id}"
   name = "${var.name}"
}
`, rand)
}

var testAccRouteEntryCheckMap = map[string]string{
	"route_table_id":        CHECKSET,
	"nexthop_id":            CHECKSET,
	"destination_cidrblock": "172.11.1.1/32",
}
