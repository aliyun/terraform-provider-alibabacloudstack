package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectBgpNetwork_basic0(t *testing.T) {
	var v *VpcDescribebgpnetworksResponse
	resourceId := "alibabacloudstack_expressconnect_bgp_network.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectBgpNetworkCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribebgpnetworksRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testaccexpressconnect-bgp-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectBgpNetworkDependence0(rand))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dst_cidr_block": "1.1.1.1",
					"router_id":      "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dst_cidr_block": "1.1.1.1",
						"router_id":      CHECKSET,
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

var AlibabacloudStackExpressconnectBgpNetworkCheckMap = map[string]string{
	"status": CHECKSET,
}

func AlibabacloudStackExpressconnectBgpNetworkDependence0(vlanId int) func(string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
    physical_connection_id = "%s"
    vlan_id = %d
    local_gateway_ip = "10.0.0.1"
    peer_gateway_ip = "10.0.0.2"
    peering_subnet_mask = "255.255.255.252"
    virtual_border_router_name = "${var.name}"
    description = "TestAccAlibabacloudStackExpressconnectBgpNetwork_basic0"
}

`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
