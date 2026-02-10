package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlicloudCloudFirewallControlPolicyOrder_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cloud_firewall_control_policy_order.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallControlPolicyOrderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallControlPolicyOrderBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testYunDunProviders(),
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_uuid":  "${alibabacloudstack_cloud_firewall_control_policy.default.0.acl_uuid}",
					"direction": "${alibabacloudstack_cloud_firewall_control_policy.default.0.direction}",
					"order":     "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order": "2",
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

var AlicloudCloudFirewallControlPolicyOrderMap0 = map[string]string{
	"acl_uuid":  CHECKSET,
	"direction": CHECKSET,
	"order":     CHECKSET,
}

func AlicloudCloudFirewallControlPolicyOrderBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	
variable name {
	default = "%s"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
count            = 3
application_name = "ANY"
acl_action       = "accept"
description      = "${var.name}_${count.index}"
destination_type = "net"
destination      = "114.2.${count.index}.0/24"
direction        = "in"
proto            = "ANY"
source           = "192.1.${count.index}.0/24"
source_type      = "net"
dest_port  = "8080/8080"
dest_port_type   = "port"
release          = "true"
}

`, name)
}
