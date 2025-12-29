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
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
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
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"destination": "0.0.0.0/16",

					"description": "testtf",

					"application_name": "ANY",

					"source_type": "net",

					"dest_port": "33/33",

					"acl_action": "log",

					"destination_type": "net",

					"source": "0.0.0.0/16",

					"dest_port_type": "port",

					"proto": "UDP",

					"new_order": "-1",

					"application_id": "0",

					"release": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination": "0.0.0.0/16",

						"description": "testtf",

						"application_name": "ANY",

						"source_type": "net",

						"dest_port": "33/33",

						"acl_action": "log",

						"destination_type": "net",

						"source": "0.0.0.0/16",

						"dest_port_type": "port",

						"proto": "UDP",

						"new_order": "-1",

						"application_id": "0",

						"release": "true",
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
					"destination":      "testIPterraform",
					"description":      "testtf111",
					"application_name": "HTTP",
					"source_type":      "net",
					"dest_port":        "55/33",
					"acl_action":       "accept",
					"destination_type": "group",
					"source":           "0.0.0.0/16",
					"dest_port_type":   "port",
					"proto":            "TCP",
					"application_id":   "7",
					"release":          true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination":             "testIPterraform",
						"description":             "testtf111",
						"application_name":        "HTTP",
						"source_type":             "net",
						"dest_port":               "55/33",
						"acl_action":              "accept",
						"destination_type":        "group",
						"source":                  "0.0.0.0/16",
						"dest_port_type":          "port",
						"proto":                   "TCP",
						"new_order":               "1",
						"application_id":          "7",
						"release":                 "true",
						"acl_uuid":                CHECKSET,
						"order":                   CHECKSET,
						"hit_times":               CHECKSET,
						"source_group_cidrs":      CHECKSET,
						"destination_group_cidrs": CHECKSET,
						"dest_port_group_ports":   CHECKSET,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCloudfirewallVpcControlPolicyCheckmap = map[string]string{
	"acl_uuid":                CHECKSET,
	"order":                   CHECKSET,
	"hit_times":               CHECKSET,
	"source_group_cidrs":      CHECKSET,
	"destination_group_cidrs": CHECKSET,
	"dest_port_group_ports":   CHECKSET,
}

func AlibabacloudTestAccCloudfirewallVpcControlPolicyBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
