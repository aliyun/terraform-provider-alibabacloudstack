package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAcmConfiguration0(t *testing.T) {
	var v *AcmDescribeconfigurationResponse

	resourceId := "alibabacloudstack_acm_configuration.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAcmConfigurationCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AcmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAcmDescribeconfigurationRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tftest%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAcmConfigurationBasicdependence)
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

					"app_name": "${var.name}",

					"content": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",

					"data_id": "${var.name}",

					"group": "DEFAULT_GROUP",

					"desc": "${var.name}",

					"type": "text",

					"beta_ips": "192.168.1.1,192.168.1.2",

					"tags": "aaaaaa,bbbbbbb",

					// "namespace_id": "${alibabacloudstack_edas_namespace.default.id}",
					"namespace_id": "ef8c9ec9-8b54-4be5-931b-7d8e5a21ca45",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name": name,

						"content": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",

						"data_id": name,

						"group": "DEFAULT_GROUP",

						"desc": name,

						"type": "text",

						"beta_ips": "192.168.1.1,192.168.1.2",

						"tags": "aaaaaa,bbbbbbb",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"app_name": "${var.name}_update",

					"desc": "${var.name}_update",

					"beta_ips": "192.168.1.1",

					"tags": "aaaaaa",

					"content": "{\\\"aaa\\\": \\\"bbb\\\"}",

					"type": "json",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"app_name": fmt.Sprintf("%s_update", name),

						"desc": fmt.Sprintf("%s_update", name),

						"beta_ips": "192.168.1.1",

						"tags": "aaaaaa",

						"content": "{\\\"aaa\\\": \\\"bbb\\\"}",

						"type": "json",
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

var AlibabacloudTestAccAcmConfigurationCheckmap = map[string]string{}

func AlibabacloudTestAccAcmConfigurationBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "logical_id" {
  default = "%s:%s"
}

// resource "alibabacloudstack_edas_namespace" "default" {
//   	description = "${var.name}"
// 	namespace_name = "${var.name}"
// 	namespace_logical_id = "${var.logical_id}"
// }
`, name, defaultRegionToTest, name)
}
