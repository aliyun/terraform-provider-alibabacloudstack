package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_express_connect_virtual_border_router", &resource.Sweeper{
		Name: "alibabacloudstack_express_connect_virtual_border_router",
		F: func(region string) error {
			return SweeperFunc(
				sweeperEntry{fn: testSweepVirtualBorderRouters, precise: []string{"tf-testacc-ecvbr"}, broad: broadPrefixes, fatal: true, waitSec: 10},
			)
		},
	})
}

// testSweepVirtualBorderRouters deletes VBRs created by acceptance tests. A VBR still
// attached to a CEN refuses to delete, so in the connect_peer cleanup chain this runs as
// the entry after testSweepCenInstances, which detaches and deletes the CEN first.
func testSweepVirtualBorderRouters(prefixes []string) error {
	rawClient, err := sharedClientForRegion("")
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	// DescribeVirtualBorderRouters on 17E honours no filter parameter at all and rejects
	// PageSize above 50, so every page is pulled and matched on name locally.
	var allVbrs []VirtualBorderRouter
	pageNumber := 1
	for {
		request := client.NewCommonRequest("GET", "Vpc", "2016-04-28", "DescribeVirtualBorderRouters", "")
		request.QueryParams["PageNumber"] = fmt.Sprintf("%d", pageNumber)
		request.QueryParams["PageSize"] = "50"
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			return fmt.Errorf("error listing VBRs: %s", err)
		}
		var response struct {
			VirtualBorderRouterSet struct {
				VirtualBorderRouterType []VirtualBorderRouter `json:"VirtualBorderRouterType"`
			} `json:"VirtualBorderRouterSet"`
		}
		if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &response); err != nil {
			return fmt.Errorf("error parsing VBR list: %s", err)
		}
		allVbrs = append(allVbrs, response.VirtualBorderRouterSet.VirtualBorderRouterType...)
		if len(response.VirtualBorderRouterSet.VirtualBorderRouterType) < 50 {
			break
		}
		pageNumber++
	}

	leftover := []string{}
	for _, vbr := range allVbrs {
		matched := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(vbr.Name, prefix) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		log.Printf("[INFO] Sweeping VBR: %s (%s)", vbr.Name, vbr.VbrId)
		request := client.NewCommonRequest("POST", "Vpc", "2016-04-28", "DeleteVirtualBorderRouter", "")
		request.QueryParams["VbrId"] = vbr.VbrId
		if _, err := client.ProcessCommonRequest(request); err != nil {
			log.Printf("[WARN] Error deleting VBR %s: %s", vbr.VbrId, err)
			leftover = append(leftover, vbr.VbrId)
		}
	}

	if len(leftover) > 0 {
		return fmt.Errorf("VBR sweep left %d resource(s) behind: %s", len(leftover), strings.Join(leftover, ", "))
	}
	return nil
}

func TestAccAlicloudExpressConnectVirtualBorderRouter_basic0(t *testing.T) {
	//	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alibabacloudstack_express_connect_virtual_border_router.default"
	ra := resourceAttrInit(resourceId, AlicloudExpressConnectVirtualBorderRouterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeExpressConnectVirtualBorderRouter")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1, 2999)
	name := fmt.Sprintf("tf-testacc-ecvbr%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudExpressConnectVirtualBorderRouterBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckPhysicalConnection(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id":     "${alibabacloudstack_expressconnect_physicalconnection.default.id}",
					"vlan_id":                    fmt.Sprint(rand),
					"local_gateway_ip":           "10.0.0.1",
					"peer_gateway_ip":            "10.0.0.2",
					"peering_subnet_mask":        "255.255.255.252",
					"virtual_border_router_name": name,
					"description":                name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": name,
						"description":                name,
						"physical_connection_id":     CHECKSET,
						"vlan_id":                    fmt.Sprint(rand),
						"local_gateway_ip":           "10.0.0.1",
						"peer_gateway_ip":            "10.0.0.2",
						"peering_subnet_mask":        "255.255.255.252",
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
					"virtual_border_router_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"detect_multiplier": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"detect_multiplier": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_rx_interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_rx_interval": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"min_tx_interval": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"min_tx_interval": "300",
					}),
				),
			},
			// The terminate and recover functions are available, but due to environment issues, the status remains in recovery and the query status times out. Commented out.
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "terminated",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "terminated",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"status": "active",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"status": "active",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ipv6": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_subnet_mask": "255.255.255.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_subnet_mask": "255.255.255.0",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_gateway_ip": "10.0.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_gateway_ip": "10.0.0.4",
					}),
				),
			},
			// Currently, the product does not support ipv6
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"enable_ipv6": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"enable_ipv6": "true",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"local_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"local_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"peer_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"peer_ipv6_gateway_ip": "2001:4004:3c4d:0015:0000:0000:0000:1a2b",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"peering_ipv6_subnet_mask": "2408:4004:cc:400::/56",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id": fmt.Sprint(rand + 1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": fmt.Sprint(rand + 1),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"virtual_border_router_name": name + "-update1",
					"status":                     "active",
					"circuit_code":               name + "-update1",
					"description":                name + "-update1",
					"detect_multiplier":          "10",
					"enable_ipv6":                "false",
					"min_rx_interval":            "300",
					"local_gateway_ip":           "192.168.0.11",
					"min_tx_interval":            "300",
					"peer_gateway_ip":            "192.168.0.12",
					"peering_subnet_mask":        "255.255.255.0",
					"vlan_id":                    fmt.Sprint(rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"virtual_border_router_name": name + "-update1",
						"status":                     "active",
						"circuit_code":               name + "-update1",
						"description":                name + "-update1",
						"detect_multiplier":          "10",
						"enable_ipv6":                "false",
						"min_rx_interval":            "300",
						"local_gateway_ip":           "192.168.0.11",
						"min_tx_interval":            "300",
						"peer_gateway_ip":            "192.168.0.12",
						"peering_subnet_mask":        "255.255.255.0",
						"vlan_id":                    fmt.Sprint(rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				//				ImportStateVerifyIgnore: []string{"vbr_owner_id", "bandwidth"},
			},
		},
	})
}

var AlicloudExpressConnectVirtualBorderRouterMap0 = map[string]string{
	"enable_ipv6":  CHECKSET,
	"vbr_owner_id": NOSET,
	"bandwidth":    NOSET,
	"status":       "active",
}

func AlicloudExpressConnectVirtualBorderRouterBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
%s
`, name, ExpressconnectPhysicalConnectionsCommonTestCase)
}
