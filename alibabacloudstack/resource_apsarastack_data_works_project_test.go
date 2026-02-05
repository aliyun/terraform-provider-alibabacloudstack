package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksProject_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_project.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksProjectMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksProject")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksProjectBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":           "${var.name}",
					"description":    "${var.name}_desc",
					"task_auth_type": "PROJECT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"description":    name + "_desc",
						"task_auth_type": "PROJECT",
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
					"allow_download": true,
					"status":         "FORBIDDEN",
					"description":    "${var.name}_desc update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_download": "true",
						"status":         "FORBIDDEN",
						"description":    name + "_desc update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_download": false,
					"status":         "AVAILABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_download": "false",
						"status":         "AVAILABLE",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDataWorksProjectMap0 = map[string]string{
	"identifier": CHECKSET,
}

func AlibabacloudStackDataWorksProjectBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
