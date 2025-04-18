package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcIpv6Internetbandwidth0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_ipv6internetbandwidth.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcIpv6InternetbandwidthCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeipv6AddressesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6_internet_bandwidth%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcIpv6InternetbandwidthBasicdependence)
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

					"ipv6_address_id": "ipv6-bp11fycipeipu0ae95nz1",

					"ipv6_gateway_id": "ipv6gw-bp1kc0t2ndkccmnbnwjsu",

					"internet_charge_type": "PayByBandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ipv6_address_id": "ipv6-bp11fycipeipu0ae95nz1",

						"ipv6_gateway_id": "ipv6gw-bp1kc0t2ndkccmnbnwjsu",

						"internet_charge_type": "PayByBandwidth",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcIpv6InternetbandwidthCheckmap = map[string]string{

	"status": CHECKSET,

	"bandwidth": CHECKSET,

	"ipv6_address_id": CHECKSET,

	"ipv6_gateway_id": CHECKSET,

	"ipv6_internet_bandwidth_id": CHECKSET,

	"internet_charge_type": CHECKSET,
}

func AlibabacloudTestAccVpcIpv6InternetbandwidthBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
