package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCloudfirewallControlpolicy0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cloudfirewall_controlpolicy.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCloudfirewallControlpolicyCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testYundunProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCloudfwDescribecontrolpolicyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloud_firewallcontrol_policy%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudfirewallControlpolicyBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
			testAccPreYunCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testYunDunProviders(),

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"destination": "${alibabacloudstack_cloudfw_address_book.default.group_name}",

					"description": "test-update",

					"application_name": "ANY",

					"source_type": "group",

					"dest_port_group": "${alibabacloudstack_cloudfw_address_book.port.group_name}",

					"acl_action": "accept",

					"destination_type": "group",

					"source": "${alibabacloudstack_cloudfw_address_book.default.group_name}",

					"dest_port_type": "group",

					"proto": "ANY",

					"direction": "in",
					"release":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination": CHECKSET,

						"description": "test-update",

						"application_name": "ANY",

						"source_type": "group",

						"dest_port": CHECKSET,

						"acl_action": "accept",

						"destination_type": "group",

						"source":          CHECKSET,
						"dest_port_group": CHECKSET,

						"dest_port_type": "group",

						"proto": "ANY",

						"direction": "in",
						"release":   "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"destination": "0.0.0.0/0",

					"description": "test",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "80/80",

					"acl_action": "accept",

					"destination_type": "net",

					"direction": "in",

					"source": "0.0.0.0/0",

					"dest_port_type": "port",

					"proto": "ANY",

					"release": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination": "0.0.0.0/0",

						"description": "test",

						"application_name": "ANY",

						"source_type": "net",

						"dest_port": "80/80",

						"acl_action": "accept",

						"destination_type": "net",

						"direction": "in",

						"source": "0.0.0.0/0",

						"dest_port_type": "port",

						"proto": "ANY",

						"release":         "true",
						"dest_port_group": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"destination": "192.1.1.0/24",

					"description": "test-update",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "8080/8080",

					"acl_action": "accept",

					"destination_type": "net",

					"source": "114.2.3.0/24",

					"dest_port_type": "port",

					"proto": "ANY",

					"direction": "in",
					"release":   "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination": "192.1.1.0/24",

						"description": "test-update",

						"application_name": "ANY",

						"source_type": "net",

						"dest_port": "8080/8080",

						"acl_action": "accept",

						"destination_type": "net",

						"source": "114.2.3.0/24",

						"dest_port_type": "port",

						"proto": "ANY",

						"direction": "in",
						"release":   "false",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudfirewallControlpolicyCheckmap = map[string]string{
	"destination":      CHECKSET,
	"description":      CHECKSET,
	"source_type":      CHECKSET,
	"dest_port":        CHECKSET,
	"destination_type": CHECKSET,
	"direction":        CHECKSET,
	"source":           CHECKSET,
	"dest_port_type":   CHECKSET,
	"proto":            CHECKSET,
	"application_name": CHECKSET,
	"acl_action":       CHECKSET,
	"acl_uuid":         CHECKSET,
	"release":          CHECKSET,
}

func AlibabacloudTestAccCloudfirewallControlpolicyBasicdependence(name string) string {
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

resource "alibabacloudstack_cloudfw_address_book" "port" {
    group_type = "port"
    group_name =  "${var.name}port"
    address_list = ["8888","9999"]
    description = "test port book"
}



`, name)
}
