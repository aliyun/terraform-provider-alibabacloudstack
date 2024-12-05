package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRosTemplate0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ros_template.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRosTemplateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoRosGettemplateRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srostemplate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRosTemplateBasicdependence)
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

					"description": "模板的描述",

					"template_name": "MyTemplateTest12",

					"template_body": "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "模板的描述",

						"template_name": "MyTemplateTest12",

						"template_body": "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "模板的描述123",

					"template_name": "TemplateName123",

					"resource_type": "template",

					"template_body": "{   \"ROSTemplateFormatVersion\": \"2015-09-01\",   \"Transform\": \"Aliyun::Terraform-v1.0\",   \"Workspace\": {     \"main.tf\": \"variable  \\\"name\\\" {  default = \\\"auto_provisioning_group\\\"}\"   },  \"Outputs\": {} }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "模板的描述123",

						"template_name": "TemplateName123",

						"resource_type": "template",

						"template_body": "{   \"ROSTemplateFormatVersion\": \"2015-09-01\",   \"Transform\": \"Aliyun::Terraform-v1.0\",   \"Workspace\": {     \"main.tf\": \"variable  \\\"name\\\" {  default = \"auto_provisioning_group\"}\"   },   \"Outputs\": {} }",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRosTemplateCheckmap = map[string]string{

	"description": CHECKSET,

	"resource_type": CHECKSET,

	"template_body": CHECKSET,

	"template_name": CHECKSET,

	"tags": CHECKSET,

	"template_id": CHECKSET,

	"stack_id": CHECKSET,
}

func AlibabacloudTestAccRosTemplateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
