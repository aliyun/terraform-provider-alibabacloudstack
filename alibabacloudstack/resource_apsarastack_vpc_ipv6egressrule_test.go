package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcIpv6Egressrule0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_vpc_ipv6egressrule.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcIpv6EgressruleCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeipv6EgressonlyrulesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6_egress_rule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcIpv6EgressruleBasicdependence)
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

					"description": "test",

					"instance_id": "ipv6-bp11fycipeipu0ae95nz1",

					"ipv6_gateway_id": "ipv6gw-bp1kc0t2ndkccmnbnwjsu",

					"ipv6_egress_rule_name": "rdk-test",

					"instance_type": "Ipv6Address",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"instance_id": "ipv6-bp11fycipeipu0ae95nz1",

						"ipv6_gateway_id": "ipv6gw-bp1kc0t2ndkccmnbnwjsu",

						"ipv6_egress_rule_name": "rdk-test",

						"instance_type": "Ipv6Address",
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

var AlibabacloudTestAccVpcIpv6EgressruleCheckmap = map[string]string{

	"ipv6_egress_rule_id": CHECKSET,

	"status": CHECKSET,

	"description": CHECKSET,

	"instance_id": CHECKSET,

	"ipv6_gateway_id": CHECKSET,

	"ipv6_egress_rule_name": CHECKSET,

	"instance_type": CHECKSET,
}

func AlibabacloudTestAccVpcIpv6EgressruleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
