package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectBgpgroup_basic0(t *testing.T) {
	var v *ExpressconnectBgpGroup
	resourceId := "alibabacloudstack_expressconnect_bgp_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectBgpgroupCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribebgpgroupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testaccexpressconnect-bgp-group%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectBgpgroupDependence0(rand))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_name": "${var.name}",
					"description":    "${var.name}",
					"local_asn":      "65534",
					"peer_asn":       "10",
					"router_id":      "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_name": name,
						"description":    name,
						"local_asn":      "65534",
						"peer_asn":       "10",
						"router_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_name": "${var.name}_test",
					"description":    "${var.name}_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bgp_group_name": fmt.Sprintf("%s_test", name),
						"description":    fmt.Sprintf("%s_test", name),
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

var AlibabacloudStackExpressconnectBgpgroupCheckMap = map[string]string{}

func AlibabacloudStackExpressconnectBgpgroupDependence0(vlanId int) func(string) string {
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
	description =                "TestAccAlibabacloudStackExpressconnectBgpgroup_basic0"
}

`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
