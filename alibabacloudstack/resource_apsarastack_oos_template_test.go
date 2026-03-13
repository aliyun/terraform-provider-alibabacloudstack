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
	name := fmt.Sprintf("tf-oostemplate%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccOosTemplateBasicdependence)
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

					"content": TfRawString(`	jsonencode({
    FormatVersion = "OOS-2019-06-01"
    Description   = "Update Describe instances of given status"
    
    Parameters = {
      Status = {
        Type        = "String"
        Description = "(Required) The status of the Ecs instance."
      }
    }

    Tasks = [
      {
        Name   = "foo"
        Action = "ACS::ExecuteApi"
        
        Properties = {
          API      = "DescribeInstances"
          Service  = "Ecs"
          Parameters = {
            Status = "{{ Status }}"
          }
        }
      }
    ]
  })`),
					"template_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":       `{"Description":"Update Describe instances of given status","FormatVersion":"OOS-2019-06-01","Parameters":{"Status":{"Description":"(Required) The status of the Ecs instance.","Type":"String"}},"Tasks":[{"Action":"ACS::ExecuteApi","Name":"foo","Properties":{"API":"DescribeInstances","Parameters":{"Status":"{{ Status }}"},"Service":"Ecs"}}]}`,
						"template_name": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"content": TfRawString(`	jsonencode({
  FormatVersion = "OOS-2019-06-01"
  Description   = "Update Describe instances of given status"
  
  Parameters = {
    Status = {
      Type        = "String"
      Description = "(Required) The status of the Ecs instance."
    }
  }

  Tasks = [
    {
      Name   = "bar"
      Action = "ACS::ExecuteApi"
      
      Properties = {
        API      = "DescribeInstances"
        Service  = "Ecs"
        Parameters = {
          Status = "{{ Status }}"
        }
      }
    }
  ]
})`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":       `{"Description":"Update Describe instances of given status","FormatVersion":"OOS-2019-06-01","Parameters":{"Status":{"Description":"(Required) The status of the Ecs instance.","Type":"String"}},"Tasks":[{"Action":"ACS::ExecuteApi","Name":"bar","Properties":{"API":"DescribeInstances","Parameters":{"Status":"{{ Status }}"},"Service":"Ecs"}}]}`,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"auto_delete_executions"},
			},
		},
	})
}

var AlibabacloudTestAccOosTemplateCheckmap = map[string]string{
	"description":      CHECKSET,
	"template_format":  CHECKSET,
	"updated_date":     CHECKSET,
	"template_version": CHECKSET,
	"updated_by":       CHECKSET,
	"has_trigger":      CHECKSET,
	"template_name":    CHECKSET,
	"template_id":      CHECKSET,
	"created_by":       CHECKSET,
	"content":          CHECKSET,
	"share_type":       CHECKSET,
}

func AlibabacloudTestAccOosTemplateBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
