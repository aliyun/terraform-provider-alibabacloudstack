package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_router_interface", &resource.Sweeper{
		Name: "alibabacloudstack_router_interface",
		F:    testSweepRouterInterfaces,
	})
}

func testSweepRouterInterfaces(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var ris []vpc.RouterInterfaceType
	req := vpc.CreateDescribeRouterInterfacesRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeRouterInterfaces(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Router Interfaces: %s", err)
		}
		resp, _ := raw.(*vpc.DescribeRouterInterfacesResponse)
		if resp == nil || len(resp.RouterInterfaceSet.RouterInterfaceType) < 1 {
			break
		}
		ris = append(ris, resp.RouterInterfaceSet.RouterInterfaceType...)

		if len(resp.RouterInterfaceSet.RouterInterfaceType) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}
	service := VpcService{client}
	for _, v := range ris {
		name := v.Name
		id := v.RouterInterfaceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a RI name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcInstanceId, ""); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Router Interface: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Router Interface: %s (%s)", name, id)
		req := vpc.CreateDeleteRouterInterfaceRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.RouterInterfaceId = id
		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Router Interface (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func testAccCheckRouterInterfaceExists(n string, ri *vpc.RouterInterfaceType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		vpcService := VpcService{client}

		response, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		ri = &response
		return nil
	}
}

func testAccCheckRouterInterfaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_router_interface" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		vpcService := VpcService{client}

		ri, err := vpcService.DescribeRouterInterface(rs.Primary.ID, client.RegionId)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.RouterInterfaceId == rs.Primary.ID {
			return errmsgs.WrapError(errmsgs.Error("Interface %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccAlibabacloudStackRouterInterface_vrouter(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alibabacloudstack_router_interface.default"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccRouterInterfaceConfigBasic)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"opposite_region": "${data.alibabacloudstack_account.current.region}",
					"router_type":     "VRouter",
					"router_id":       "${alibabacloudstack_vpc_vpc.default.router_id}",
					"role":            "AcceptingSide",
					"name":            "${var.name}",
					"description":     "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"router_type": "VRouter",
						"role":        "AcceptingSide",
						"name":        name,
						"description": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role":          "InitiatingSide",
					"specification": "Large.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"role":          "InitiatingSide",
						"specification": "Large.2",
					}),
				),
			},
//			{
//				Config: testAccConfig(map[string]interface{}{
//					"health_check_source_ip": "172.16.1.100",
//					"health_check_target_ip": "172.16.2.100",
//				}),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheck(map[string]string{
//						"health_check_source_ip": "172.16.1.100",
//						"health_check_target_ip": "172.16.2.100",
//					}),
//				),
//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"specification": "Large.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"specification": "Large.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"description": name + "_desc",
					}),
				),
			},
		},
	})

}
func TestAccAlibabacloudStackRouterInterface_vbr(t *testing.T) {
	var v vpc.RouterInterfaceType
	resourceId := "alibabacloudstack_router_interface.default"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceCheckMap)

	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testAccRouterInterfaceConfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccRouterInterfaceConfigVbr)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPhysicalConnection(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"opposite_region": "${data.alibabacloudstack_account.current.region}",
					"router_type":     "VBR",
					"router_id":       "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
					"role":            "InitiatingSide",
					"specification":   "Large.2",
					"name":            "${var.name}",
					"description":     "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceExists(resourceId, &v),
					testAccCheck(map[string]string{
						"router_type": "VBR",
						"role":        "InitiatingSide",
						"name":        name,
						"description": name,
					}),
				),
			},
		},
	})
}
func testAccRouterInterfaceConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

data "alibabacloudstack_account" "current"{
}
`, name, VpcCommonTestCase)
}

func testAccRouterInterfaceConfigVbr(name string) string {
	return testAccRouterInterfaceConfigBasic(name) + fmt.Sprintf(`
%s
resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = alibabacloudstack_expressconnect_physicalconnection.default.id
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}`, ExpressconnectPhysicalConnectionsCommonTestCase, getAccTestRandInt(1, 2999))
}

var testAccRouterInterfaceCheckMap = map[string]string{
	"opposite_region":        CHECKSET,
	"router_id":              CHECKSET,
	"role":                   "AcceptingSide",
	"description":            "",
	"health_check_source_ip": "",
	"health_check_target_ip": "",
}
