package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckRouterInterfaceConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		vpcService := VpcService{client}

		response, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		if response.Status != string(Active) {
			return fmt.Errorf("Error connection router interface id %s is not Active.", response.RouterInterfaceId)
		}

		return nil
	}
}

func testAccCheckRouterInterfaceConnectionDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "apsarastack_router_interface_connection" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		vpcService := VpcService{client}

		ri, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if ri.Status == string(Active) {
			return WrapError(Error("Interface connection %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccApsaraStackRouterInterfaceConnectionBasic(t *testing.T) {
	resourceId := "apsarastack_router_interface_connection.foo"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceConnectionCheckMap)
	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterInterfaceConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConnectionConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceConnectionExists(resourceId),
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRouterInterfaceConnectionConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
provider "apsarastack" {
  region = "${var.region}"
}
variable "region" {
  default = "cn-qingdao-env66-d01"
}
variable "name" {
  default = "tf-test%d"
}
resource "apsarastack_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}
resource "apsarastack_vpc" "bar" {
  provider = "apsarastack"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}
resource "apsarastack_router_interface" "initiate" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${apsarastack_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"

}
resource "apsarastack_router_interface" "opposite" {
  provider = "apsarastack"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${apsarastack_vpc.bar.router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}-opposite"
}

resource "apsarastack_router_interface_connection" "foo" {
  interface_id = "${apsarastack_router_interface.initiate.id}"
  opposite_interface_id = "${apsarastack_router_interface.opposite.id}"
  depends_on = ["apsarastack_router_interface_connection.bar"]
  opposite_interface_owner_id = "1262302482727553"
  opposite_router_id = apsarastack_vpc.foo.router_id
  opposite_router_type = "VRouter"
}

resource "apsarastack_router_interface_connection" "bar" {
  provider = "apsarastack"
  interface_id = "${apsarastack_router_interface.opposite.id}"
  opposite_interface_id = "${apsarastack_router_interface.initiate.id}"
  opposite_interface_owner_id =  "1262302482727553"
  opposite_router_id = apsarastack_vpc.bar.router_id
  opposite_router_type = "VRouter"
}
`, rand)
}

var testAccRouterInterfaceConnectionCheckMap = map[string]string{
	"interface_id":                CHECKSET,
	"opposite_interface_id":       CHECKSET,
	"opposite_router_type":        "VRouter",
	"opposite_router_id":          CHECKSET,
	"opposite_interface_owner_id": CHECKSET,
}
