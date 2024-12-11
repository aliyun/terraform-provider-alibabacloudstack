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
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCloudfwDescribecontrolpolicyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloud_firewallcontrol_policy%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCloudfirewallControlpolicyBasicdependence)
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

					"destination": "0.0.0.0/0",

					"description": "test",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "80",

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

						"dest_port": "80",

						"acl_action": "accept",

						"destination_type": "net",

						"direction": "in",

						"source": "0.0.0.0/0",

						"dest_port_type": "port",

						"proto": "ANY",

						"release": "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"destination": "192.1.1.0/24",

					"description": "test-update",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "8080",

					"acl_action": "accept",

					"destination_type": "net",

					"source": "114.2.3.0/24",

					"dest_port_type": "port",

					"proto": "ANY",

					"direction": "in",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"destination": "192.1.1.0/24",

						"description": "test-update",

						"application_name": "ANY",

						"source_type": "net",

						"dest_port": "8080",

						"acl_action": "accept",

						"destination_type": "net",

						"source": "114.2.3.0/24",

						"dest_port_type": "port",

						"proto": "ANY",

						"direction": "in",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudfirewallControlpolicyCheckmap = map[string]string{

	"destination": CHECKSET,

	"description": CHECKSET,

	"source_type": CHECKSET,

	"dest_port": CHECKSET,

	"destination_type": CHECKSET,

	"direction": CHECKSET,

	"source": CHECKSET,

	"dest_port_type": CHECKSET,

	"proto": CHECKSET,

	"dest_port_group_ports": CHECKSET,

	"destination_group_cidrs": CHECKSET,

	"repeat_days": CHECKSET,

	"release": CHECKSET,

	"source_group_cidrs": CHECKSET,

	"dest_port_group": CHECKSET,

	"order": CHECKSET,

	"application_name": CHECKSET,

	"application_name_list": CHECKSET,

	"acl_action": CHECKSET,

	"acl_uuid": CHECKSET,
}

func AlibabacloudTestAccCloudfirewallControlpolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
