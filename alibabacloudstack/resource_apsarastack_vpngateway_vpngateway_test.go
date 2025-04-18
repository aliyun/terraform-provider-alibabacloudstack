package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpngatewayVpngateway0(t *testing.T) {
	var v vpc.DescribeVpnGatewayResponse

	resourceId := "alibabacloudstack_vpngateway_vpngateway.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpngatewayVpngatewayCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeVpnGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpn_gatewayvpn_gateway%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpngatewayVpngatewayBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test_vpn",

					"vpn_gateway_name": "test_vpn",

					"bandwidth": "10",

					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",

					"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",

					"enable_ssl": "true",

					"instance_charge_type": "PostPaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test_vpn",

						"vpn_gateway_name": "test_vpn",

						"bandwidth": "10",

						"vswitch_id": CHECKSET,

						"vpc_id": CHECKSET,

						"enable_ssl": "true",

						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"description": "tes_vpn_new",

					"vpn_gateway_name": "tes_vpn_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "tes_vpn_new",

						"vpn_gateway_name": "tes_vpn_new",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpngatewayVpngatewayCheckmap = map[string]string{}

func AlibabacloudTestAccVpngatewayVpngatewayBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

`, name, VSwitchCommonTestCase)
}
