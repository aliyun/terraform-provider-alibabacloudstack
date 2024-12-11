package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOosTemplate0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_oos_template.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccOosTemplateCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoOosGettemplateRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soostemplate%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccOosTemplateBasicdependence)
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

					"content": "{\"FormatVersion\": \"OOS-2019-06-01\", \"Description\": \"test\", \"Parameters\": {\"Status\": {\"Type\": \"String\", \"Description\": \"test\"}}",

					"template_name": "rdk-test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"content": "{\"FormatVersion\": \"OOS-2019-06-01\", \"Description\": \"test\", \"Parameters\": {\"Status\": {\"Type\": \"String\", \"Description\": \"test\"}}",

						"template_name": "rdk-test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"content": "{\"FormatVersion\": \"OOS-2023-11-29\", \"Description\": \"test\", \"Parameters\": {\"Status\": {\"Type\": \"String\", \"Description\": \"test\"}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"content": "{\"FormatVersion\": \"OOS-2023-11-29\", \"Description\": \"test\", \"Parameters\": {\"Status\": {\"Type\": \"String\", \"Description\": \"test\"}}",
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

var AlibabacloudTestAccOosTemplateCheckmap = map[string]string{

	"description": CHECKSET,

	"template_format": CHECKSET,

	"updated_date": CHECKSET,

	"template_version": CHECKSET,

	"updated_by": CHECKSET,

	"has_trigger": CHECKSET,

	"template_name": CHECKSET,

	"tags": CHECKSET,

	"template_id": CHECKSET,

	"created_by": CHECKSET,

	"create_time": CHECKSET,

	"content": CHECKSET,

	"share_type": CHECKSET,
}

func AlibabacloudTestAccOosTemplateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
