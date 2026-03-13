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
	name := fmt.Sprintf("tf-oosexecution%d", rand)

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
					"parameters":       TfRawString(`jsonencode({Status="Running"})`),
					"template_name":    "${alibabacloudstack_oos_template.default.template_name}",
					"mode":             "Automatic",
					"template_version": "v1",
					"description":      "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters":       `{"Status":"Running"}`,
						"template_name":    name,
						"mode":             "Automatic",
						"template_version": "v1",
						"description":      "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": TfRawString(`jsonencode({Status="Update"})`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters": `{"Status":"Update"}`,
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

var AlibabacloudTestAccOosExecutionCheckmap = map[string]string{
	"category":         CHECKSET,
	"description":      CHECKSET,
	"template_version": CHECKSET,
	"start_date":       CHECKSET,
	"update_date":      CHECKSET,
	"template_name":    CHECKSET,
	"executed_by":      CHECKSET,
	"template_id":      CHECKSET,
	"status":           CHECKSET,
	"parameters":       CHECKSET,
	"is_parent":        CHECKSET,
	"mode":             CHECKSET,
	"end_date":         CHECKSET,
	"outputs":          CHECKSET,
}

func AlibabacloudTestAccOosExecutionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_oos_template" "default" {
  content= <<EOF
  {
	"FormatVersion": "OOS-2019-06-01",
	"Description": "Update Describe instances of given status",
	"Parameters":{
	  "Status":{
		"Type": "String",
		"Description": "(Required) The status of the Ecs instance."
	  }
	},
	"Tasks": [
	  {
		"Properties" :{
		  "Parameters":{
			"Status": "{{ Status }}"
		  },
		  "API": "DescribeInstances",
		  "Service": "Ecs"
		},
		"Name": "foo",
		"Action": "ACS::ExecuteApi"
	  }]
  }
  EOF
  template_name = var.name
  version_name = "test"
}

`, name)
}
