package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksFolder_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_folder.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksFolderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksFolder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_folder%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksFolderBasicDependence0)
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
					"project_id":    "${alibabacloudstack_data_works_project.default.id}",
					"business_name": "${alibabacloudstack_data_works_business.default.name}",
					"engine_type":   "General",
					"folder_path":   "folder_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_path": "folder_test",
					}),
				),
			},
			// FIXME: popapi can not move folder
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"folder_path": "folder_test_update",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"folder_path": "folder_test_update",
			//					}),
			//				),
			//			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// AlibabacloudStackDataWorksFolderMap0 is the expected state map for the DataWorks folder resource
var AlibabacloudStackDataWorksFolderMap0 = map[string]string{
	"folder_id":  CHECKSET,
	"project_id": CHECKSET,
}

// AlibabacloudStackDataWorksFolderBasicDependence0 returns the basic dependencies for the DataWorks folder resource test
func AlibabacloudStackDataWorksFolderBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
	}

resource "alibabacloudstack_data_works_business" "default" {
	project_id=  "${alibabacloudstack_data_works_project.default.id}"
	name=        "${var.name}"
	description= "${var.name}_desc"
}
	
`, name)
}
