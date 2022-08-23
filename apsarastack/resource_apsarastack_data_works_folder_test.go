package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDataWorksFolder_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_data_works_folder.default"
	ra := resourceAttrInit(resourceId, ApsaraStackDataWorksFolderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDataWorksFolder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksfolder%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDataWorksFolderBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":  "10023",
					"folder_path": "业务流程/test/folderUserDefined/testcxt",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":  "10023",
						"folder_path": "业务流程/test/folderUserDefined/testcxt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"folder_path": "业务流程/test/folderUserDefined/testcxt2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"folder_path": "业务流程/test/folderUserDefined/testcxt2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"project_identifier"},
			},
		},
	})
}

var ApsaraStackDataWorksFolderMap0 = map[string]string{
	"folder_id":          CHECKSET,
	"folder_path":        "",
	"project_identifier": NOSET,
	"project_id":         "10023",
}

func ApsaraStackDataWorksFolderBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
