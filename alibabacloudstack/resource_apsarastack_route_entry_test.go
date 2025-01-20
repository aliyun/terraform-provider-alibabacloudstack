package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudstackRouteEntryInstance(t *testing.T) {
	var v *vpc.RouteEntry
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "alibabacloudstack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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

func TestAccAlibabacloudstackRouteEntryInterface(t *testing.T) {
	var v *vpc.RouteEntry
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "alibabacloudstack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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

func TestAccAlibabacloudstackRouteEntryNatGateway(t *testing.T) {
	var v *vpc.RouteEntry
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "alibabacloudstack_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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

func TestAccAlibabacloudstackRouteEntryMulti(t *testing.T) {
	var v *vpc.RouteEntry
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "alibabacloudstack_route_entry.default.2"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
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
						"destination_cidrblock": "172.16.2.0/24",
						"name":                  fmt.Sprintf("tf-testAccRouteEntryConcurrence%d", rand),
					}),
				),
			},
		},
	})
}

func testAccCheckRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_route_entry" || rs.Type != "alibabacloudstack_route_entry" {
			continue
		}
		entry, err := vpcService.DescribeRouteEntry(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}

		if entry.RouteTableId != "" {
			return errmsgs.WrapError(errmsgs.Error("Route entry still exist"))
		}
	}

	//	testAccCheckRouterInterfaceDestroy(s)

	return nil
}

func testAccRouteEntryConfig_instance(rand int) string {
	return fmt.Sprintf(`

variable "name" {
   default = "tf-testAccRouteEntryConfigName%d"
}


%s

resource "alibabacloudstack_route_entry" "default" {
	route_table_id = "${alibabacloudstack_vpc_vpc.default.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alibabacloudstack_ecs_instance.default.id}"
	name = "${var.name}"
 }
	`, rand, ECSInstanceCommonTestCase)
}

func testAccRouteEntryConfig_interface(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccRouteEntryInterfaceConfig%d"
	}
%s

resource "alibabacloudstack_router_interface" "default" {
  opposite_region = "%s"
  router_type = "VRouter"
  router_id = "${alibabacloudstack_vpc_vpc.default.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "${var.name}"
}
resource "alibabacloudstack_route_entry" "default" {
  route_table_id = "${alibabacloudstack_vpc_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "RouterInterface"
  nexthop_id = "${alibabacloudstack_router_interface.default.id}"
  name = "${var.name}"
}
`, rand, SecurityGroupCommonTestCase, os.Getenv("ALIBABACLOUDSTACK_REGION"),)
}

func testAccRouteEntryConfig_natGateway(rand int) string {
	return fmt.Sprintf(
		`
data "alibabacloudstack_zones" "default" {
  available_resource_creation= "VSwitch"
}
variable "name" {
   default = "tf-testAccRouteEntryNatGatewayConfig%d"
}
resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
resource "alibabacloudstack_nat_gateway" "default" {
   vpc_id = "${alibabacloudstack_vpc.default.id}"
   specification = "Middle"
   name = "${var.name}"
}
resource "alibabacloudstack_route_entry" "default" {
  route_table_id = "${alibabacloudstack_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "NatGateway"
  nexthop_id = "${alibabacloudstack_nat_gateway.default.id}"
  name = "${var.name}"
}`, rand)
}

func testAccRouteEntryConfigMulti(rand int) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
   available_resource_creation= "VSwitch"
}
variable "name" {
   default = "tf-testAccRouteEntryConcurrence%d"
}
resource "alibabacloudstack_vpc" "default" {
   name = "${var.name}"
   cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
    name = "${var.name}"
    cidr_block = "10.1.1.0/24"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
}
resource "alibabacloudstack_route_entry" "default" {
   count = 3
   route_table_id = "${alibabacloudstack_vpc.default.route_table_id}"
   destination_cidrblock = "172.16.${count.index}.0/24"
   nexthop_type = "NetworkInterface"
   nexthop_id = "${alibabacloudstack_network_interface.default.id}"
   name = "${var.name}"
}
`, rand)
}

var testAccRouteEntryCheckMap = map[string]string{
	"route_table_id":        CHECKSET,
	"nexthop_id":            CHECKSET,
	"destination_cidrblock": "172.11.1.1/32",
}
