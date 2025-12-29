package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudfirewallVpcControlPolicy_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cloudfw_vpc_control_policy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCloudfirewallVpcControlPolicyCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCloudfwVpcControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc_vpc_control_policy%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudfirewallVpcControlPolicyBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testYunDunProviders(),

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "0.0.0.0/16",

					"description": "${var.name}",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "33/33",

					"acl_action": "log",

					"destination_type": "net",

					"source": "0.0.0.0/16",

					"dest_port_type": "port",

					"proto": "UDP",

					"application_id": "0",

					"release": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "0.0.0.0/16",

						"description": name,

						"application_name": "ANY",

						"source_type": "net",

						"dest_port": "33/33",

						"acl_action": "log",

						"destination_type": "net",

						"source": "0.0.0.0/16",

						"dest_port_type": "port",

						"proto": "UDP",

						"application_id": "0",

						"release": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"new_order"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination":      "${alibabacloudstack_cloudfw_address_book.default.group_name}",
					"description":      "${var.name}_updated",
					"application_name": "HTTP",
					"source_type":      "net",
					"dest_port":        "55/55",
					"acl_action":       "accept",
					"destination_type": "group",
					"source":           "0.0.0.0/16",
					"proto":            "TCP",
					"application_id":   "7",
					"release":          true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":      name,
						"description":      name + "_updated",
						"application_name": "HTTP",
						"source_type":      "net",
						"dest_port":        "55/55",
						"acl_action":       "accept",
						"destination_type": "group",
						"source":           "0.0.0.0/16",
						"proto":            "TCP",
						"application_id":   "7",
						"release":          "true",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudfirewallVpcControlPolicyCheckmap = map[string]string{
	"acl_uuid":  CHECKSET,
	"order":     CHECKSET,
	"hit_times": CHECKSET,
}

func AlibabacloudTestAccCloudfirewallVpcControlPolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


resource "alibabacloudstack_cloudfw_address_book" "default" {
    group_type = "ip"
    group_name = var.name
    address_list = ["100.100.100.100/30"]
    description = "test address book"
}

`, name)
}
