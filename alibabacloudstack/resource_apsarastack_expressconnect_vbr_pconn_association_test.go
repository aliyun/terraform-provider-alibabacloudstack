package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectVbrPconn_basic0(t *testing.T) {
	var v *VbrAssociatedPhysicalConnection

	resourceId := "alibabacloudstack_expressconnect_vbr_pconn_association.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectVbrPconnCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeVbrpconnassociationRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testaccexpressconnect-vbrpconn%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectVbrPconnDependence0(rand))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id":   getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD"),
					"vbr_id":                   "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
					"vlan_id":                  fmt.Sprintf("%d", rand+3),
					"local_gateway_ip":         "10.100.0.1",
					"peer_gateway_ip":          "10.100.0.2",
					"peering_subnet_mask":      "255.255.255.0",
					"enable_ipv6":              "true",
					"local_ipv6_gateway_ip":    "2408:4004:cc:500::1",
					"peer_ipv6_gateway_ip":     "2408:4004:cc:500::2",
					"peering_ipv6_subnet_mask": "2408:4004:cc:500::/56",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_id":   getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD"),
						"local_gateway_ip":         "10.100.0.1",
						"peer_gateway_ip":          "10.100.0.2",
						"peering_subnet_mask":      "255.255.255.0",
						"local_ipv6_gateway_ip":    "2408:4004:cc:500::1",
						"peer_ipv6_gateway_ip":     "2408:4004:cc:500::2",
						"peering_ipv6_subnet_mask": "2408:4004:cc:500::/56",
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

var AlibabacloudStackExpressconnectVbrPconnCheckMap = map[string]string{
	"vbr_id":  CHECKSET,
	"vlan_id": CHECKSET,
}

func AlibabacloudStackExpressconnectVbrPconnDependence0(vlanId int) func(string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id =     "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	enable_ipv6              = true
	local_ipv6_gateway_ip = "2408:4004:cc:400::1"
	peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
	peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
}

`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
