package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

	name := fmt.Sprintf("tf-testAccRouteEntryConfigName%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRouteEntryInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alibabacloudstack_vpc_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alibabacloudstack_ecs_instance.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "Instance",
						"name":         name,
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

	name := fmt.Sprintf("tf-testAccRouteEntryInterfaceConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRouteEntryInterfaceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alibabacloudstack_vpc_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "RouterInterface",
					"nexthop_id":            "${alibabacloudstack_router_interface.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "RouterInterface",
						"name":         name,
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

	name := fmt.Sprintf("tf-testAccRouteEntryNatGatewayConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRouteEntryNatGatewayConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alibabacloudstack_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alibabacloudstack_nat_gateway.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "NatGateway",
						"name":         name,
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

	name := fmt.Sprintf("tf-testAccRouteEntryConcurrence%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRouteEntryMultiConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                 "3",
					"route_table_id":        "${alibabacloudstack_vpc.default.route_table_id}",
					"destination_cidrblock": "172.16.${count.index}.0/24",
					"nexthop_type":          "NetworkInterface",
					"nexthop_id":            "${alibabacloudstack_network_interface.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":          "NetworkInterface",
						"destination_cidrblock": "172.16.2.0/24",
						"name":                  name,
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

	return nil
}

func resourceRouteEntryInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, ECSInstanceCommonTestCase)
}

func resourceRouteEntryInterfaceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_router_interface" "default" {
	opposite_region = "%s"
	router_type     = "VRouter"
	router_id       = "${alibabacloudstack_vpc_vpc.default.router_id}"
	role            = "InitiatingSide"
	specification   = "Large.2"
	name            = "${var.name}"
	description     = "${var.name}"
}
`, name, SecurityGroupCommonTestCase, os.Getenv("ALIBABACLOUDSTACK_REGION"))
}

func resourceRouteEntryNatGatewayConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id              = "${alibabacloudstack_vpc.default.id}"
	cidr_block          = "10.1.1.0/24"
	availability_zone   = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name                = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id        = "${alibabacloudstack_vpc.default.id}"
	specification = "Middle"
	name          = "${var.name}"
}
`, name)
}

func resourceRouteEntryMultiConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "default" {
	name              = "${var.name}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	vpc_id            = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_network_interface" "default" {
	name            = "${var.name}"
	vswitch_id      = "${alibabacloudstack_vswitch.default.id}"
	security_groups = ["${alibabacloudstack_security_group.default.id}"]
}
`, name)
}

var testAccRouteEntryCheckMap = map[string]string{
	"route_table_id":        CHECKSET,
	"nexthop_id":            CHECKSET,
	"destination_cidrblock": "172.11.1.1/32",
}
