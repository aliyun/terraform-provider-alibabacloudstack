package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpngatewayVpngateway0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpngateway_vpngateway.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpngatewayVpngatewayCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpngatewayRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpn_gatewayvpn_gateway%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpngatewayVpngatewayBasicdependence)
	resource.Test(t, resource.TestCase{
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

					"spec": "10",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVswitch_1.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test_vpn",

						"vpn_gateway_name": "test_vpn",

						"spec": "10",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVswitch_1.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
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

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test_vpn",

					"vpn_gateway_name": "test_vpn",

					"spec": "10",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVswitch_1.VSwitchId)}}",

					"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test_vpn",

						"vpn_gateway_name": "test_vpn",

						"spec": "10",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::2.0.0.2.pre::defaultVswitch_1.VSwitchId)}}",

						"vpc_id": "${{ref(resource, VPC::VPC::4.0.0.26.pre::defaultVpc.VpcId)}}",

						"payment_type": "PayAsYouGo",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpngatewayVpngatewayCheckmap = map[string]string{

	"ipsec_vpn": CHECKSET,

	"ssl_vpn": CHECKSET,

	"description": CHECKSET,

	"end_time": CHECKSET,

	"business_status": CHECKSET,

	"vpn_instance_id": CHECKSET,

	"internet_ip": CHECKSET,

	"payment_type": CHECKSET,

	"ssl_max_connections": CHECKSET,

	"status": CHECKSET,

	"vpn_gateway_name": CHECKSET,

	"create_time": CHECKSET,

	"vswitch_id": CHECKSET,

	"vpc_id": CHECKSET,

	"spec": CHECKSET,
}

func AlibabacloudTestAccVpngatewayVpngatewayBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
