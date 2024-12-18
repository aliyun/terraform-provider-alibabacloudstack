package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOosExecution0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_oos_execution.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccOosExecutionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OosService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoOosListexecutionsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%soosexecution%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccOosExecutionBasicdependence)
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

					"parameters": "{\"Status\":\"Running\"}",

					"template_name": "MyTemplate",

					"mode": "Automatic",

					"template_version": "v1",

					"description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"parameters": "{\"Status\":\"Running\"}",

						"template_name": "MyTemplate",

						"mode": "Automatic",

						"template_version": "v1",

						"description": "test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"parameters": "{\"Status\":\"Update\"}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"parameters": "{\"Status\":\"Update\"}",
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

var AlibabacloudTestAccOosExecutionCheckmap = map[string]string{

	"parent_execution_id": CHECKSET,

	"category": CHECKSET,

	"description": CHECKSET,

	"template_version": CHECKSET,

	"start_date": CHECKSET,

	"update_date": CHECKSET,

	"template_name": CHECKSET,

	"executed_by": CHECKSET,

	"tags": CHECKSET,

	"template_id": CHECKSET,

	"status": CHECKSET,

	"parameters": CHECKSET,

	"is_parent": CHECKSET,

	"create_time": CHECKSET,

	"mode": CHECKSET,

	"end_date": CHECKSET,

	"status_message": CHECKSET,

	"outputs": CHECKSET,

	"ram_role": CHECKSET,

	"counters": CHECKSET,

	"execution_id": CHECKSET,
}

func AlibabacloudTestAccOosExecutionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
