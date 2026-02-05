package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksBusiness_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_business.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksBusinessMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksBusiness")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_business%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksBusinessBasicDependence0)
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
					"project_id":  "${alibabacloudstack_data_works_project.default.id}",
					"name":        "${var.name}",
					"description": "${var.name}_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"description": name + "_desc",
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
					"name":        "${var.name}_update",
					"description": "${var.name}_desc update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name + "_update",
						"description": name + "_desc update",
					}),
				),
			},
		},
	})
}

// AlibabacloudStackDataWorksBusinessMap0 is the expected state map for the DataWorks Business resource
var AlibabacloudStackDataWorksBusinessMap0 = map[string]string{
	"project_id": CHECKSET,
}

// AlibabacloudStackDataWorksBusinessBasicDependence0 returns the basic dependencies for the DataWorks Business resource test
func AlibabacloudStackDataWorksBusinessBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
	}

`, name)
}
