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
					"app_name":     "${var.name}",
					"content":      "step1",
					"data_id":      "${var.name}",
					"group":        "DEFAULT_GROUP",
					"desc":         "${var.name}",
					"type":         "text",
					"tags":         "aaaaaa,bbbbbbb",
					"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name": name,
						"content":  "step1",
						"data_id":  name,
						"group":    "DEFAULT_GROUP",
						"desc":     name,
						"type":     "text",
						"tags":     "aaaaaa,bbbbbbb",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"beta_app_name": "${var.name}_beta",
					"beta_ips":      "192.168.1.1",
					"beta_content":  "{step2}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"beta_app_name": fmt.Sprintf("%s_beta", name),
						"beta_ips":      "192.168.1.1",
						"beta_content":  "{step2}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":          "${var.name}_update",
					"desc":              "${var.name}_update",
					"tags":              "aaaaaa",
					"beta_ips":          REMOVEKEY,
					"content":           "{\\\"aaa\\\": \\\"bbb\\\"}",
					"type":              "json",
					"encrypt_algorithm": "AES_128",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":          fmt.Sprintf("%s_update", name),
						"desc":              fmt.Sprintf("%s_update", name),
						"tags":              "aaaaaa",
						"beta_ips":          REMOVEKEY,
						"beta_content":      REMOVEKEY,
						"beta_app_name":     REMOVEKEY,
						"content":           "{\"aaa\": \"bbb\"}",
						"type":              "json",
						"encrypt_algorithm": "AES_128",
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
					"data_id": "${var.name}_forceNew",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"data_id": fmt.Sprintf("%s_forceNew", name),
					}),
				),
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

resource "alibabacloudstack_edas_namespace" "default" {
	description = "${var.name}"
	namespace_name = "${var.name}"
	namespace_logical_id = "${var.logical_id}"
}
`, name, defaultRegionToTest, name)
}
