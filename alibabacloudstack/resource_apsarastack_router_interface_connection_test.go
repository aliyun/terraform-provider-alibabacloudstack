package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

func TestAccAlibabacloudStackRouterInterfaceConnection_initiating(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alibabacloudstack_router_interface_connection.initiating"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceConnectionCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRouterInterfaceConnection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccRouterInterfaceConnectionConfigInitiating)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterInterfaceConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"interface_id":                "${alibabacloudstack_router_interface.initiating.id}",
					"opposite_interface_id":       "${alibabacloudstack_router_interface.opposite.id}",
					"opposite_router_id":          "${alibabacloudstack_vpc.default.1.router_id}",
					"opposite_router_type":        "VRouter",
					"depends_on":                  []string{"alibabacloudstack_router_interface_connection.opposite"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"interface_id":                CHECKSET,
						"opposite_interface_id":       CHECKSET,
						"opposite_router_type":        "VRouter",
						"opposite_router_id":          CHECKSET,
						"opposite_interface_owner_id": CHECKSET,
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

func testAccCheckRouterInterfaceConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_router_interface_connection" {
			continue
		}

		object, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}

		if object.Status != string(Inactive) {
			return errmsgs.WrapError(errmsgs.Error("Router interface connection still exists and is not inactive"))
		}
	}

	return nil
}

func testAccRouterInterfaceConnectionConfigOpposite(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_account" "current" {}

resource "alibabacloudstack_vpc" "default" {
  count      = 2
  name       = "${var.name}_${count.index}"
  cidr_block = element(["172.16.0.0/12", "192.168.0.0/16"], count.index)
}

resource "alibabacloudstack_router_interface" "initiating" {
  opposite_region = data.alibabacloudstack_account.current.region
  router_type     = "VRouter"
  router_id       = alibabacloudstack_vpc.default.1.router_id
  role            = "InitiatingSide"
  specification   = "Large.2"
  name            = "${var.name}_initiating"
}

resource "alibabacloudstack_router_interface" "opposite" {
  opposite_region = data.alibabacloudstack_account.current.region
  router_type     = "VRouter"
  router_id       = alibabacloudstack_vpc.default.0.router_id
  role            = "AcceptingSide"
  specification   = "Large.1"
  name            = "${var.name}_opposite"
}

`, name, DataZoneCommonTestCase)
}

func testAccRouterInterfaceConnectionConfigInitiating(name string) string {
	return fmt.Sprintf(`
	%s
	resource alibabacloudstack_router_interface_connection "opposite"{
		interface_id                  = alibabacloudstack_router_interface.opposite.id
		opposite_interface_id         = alibabacloudstack_router_interface.initiating.id
		opposite_router_id            = alibabacloudstack_vpc.default.0.router_id
		opposite_router_type          = "VRouter"
	}
`, testAccRouterInterfaceConnectionConfigOpposite(name))
}

var testAccRouterInterfaceConnectionCheckMap = map[string]string{
	"interface_id":                CHECKSET,
	"opposite_interface_id":       CHECKSET,
	"opposite_router_type":        "VRouter",
	"opposite_router_id":          CHECKSET,
	"opposite_interface_owner_id": CHECKSET,
}
