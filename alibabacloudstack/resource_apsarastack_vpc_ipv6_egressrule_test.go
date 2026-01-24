package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVPCIpv6EgressRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_vpc_ipv6_egress_rule.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackVPCIpv6EgressRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeVpcIpv6EgressRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6egressrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackVPCIpv6EgressRuleBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_egress_rule_name": "${var.name}",
					"ipv6_gateway_id":       "${data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id}",
					"instance_id":           "${data.alibabacloudstack_vpc_ipv6_addresses.default.ids.0}",
					"instance_type":         "Ipv6Address",
					"description":           "${var.name}",
					"depends_on":            []string{"alibabacloudstack_vpc_ipv6_internet_bandwidth.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_egress_rule_name": name,
						"ipv6_gateway_id":       CHECKSET,
						"instance_id":           CHECKSET,
						"instance_type":         "Ipv6Address",
						"description":           name,
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

var AlibabacloudStackVPCIpv6EgressRuleMap0 = map[string]string{
	"instance_type": CHECKSET,
	"status":        CHECKSET,
}

func AlibabacloudStackVPCIpv6EgressRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 

variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = alibabacloudstack_ecs_instance.default.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "default" {
  vpc_id            = alibabacloudstack_vpc_vpc.default.id
  ipv6_gateway_name = var.name
  description       = var.name
  spec              = "Medium"
}

resource "alibabacloudstack_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
  depends_on = ["alibabacloudstack_vpc_ipv6_gateway.default"]
}

`, name, ECSInstanceCommonTestCase)
}
