package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNatgatewayForwardentry0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_natgateway_forwardentry.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNatgatewayForwardentryCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribeforwardtableentriesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snat_gatewayforward_entry%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNatgatewayForwardentryBasicdependence)
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

					"external_port": "12",

					"external_ip": "121.199.28.61",

					"ip_protocol": "tcp",

					"internal_port": "14",

					"internal_ip": "10.20.0.179",

					"forward_table_id": "ftb-bp15v0451hrgb2mg3vrz0",

					"forward_entry_name": "rdk_test_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"external_port": "12",

						"external_ip": "121.199.28.61",

						"ip_protocol": "tcp",

						"internal_port": "14",

						"internal_ip": "10.20.0.179",

						"forward_table_id": "ftb-bp15v0451hrgb2mg3vrz0",

						"forward_entry_name": "rdk_test_name",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"forward_entry_name": "rdk_test_name_after",

					"external_port": "80",

					"internal_port": "8080",

					"ip_protocol": "udp",

					"forward_table_id": "ftb-bp15v0451hrgb2mg3vrz0",

					"internal_ip": "10.20.0.189",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"forward_entry_name": "rdk_test_name_after",

						"external_port": "80",

						"internal_port": "8080",

						"ip_protocol": "udp",

						"forward_table_id": "ftb-bp15v0451hrgb2mg3vrz0",

						"internal_ip": "10.20.0.189",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNatgatewayForwardentryCheckmap = map[string]string{

	"status": CHECKSET,

	"external_port": CHECKSET,

	"forward_table_id": CHECKSET,

	"external_ip": CHECKSET,

	"forward_entry_id": CHECKSET,

	"ip_protocol": CHECKSET,

	"internal_port": CHECKSET,

	"forward_entry_name": CHECKSET,

	"internal_ip": CHECKSET,
}

func AlibabacloudTestAccNatgatewayForwardentryBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
