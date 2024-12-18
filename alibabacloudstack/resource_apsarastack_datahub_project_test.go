package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDatahubProject0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_datahub_project.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubProjectCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGetkafkagroupRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdata_hubproject%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDatahubProjectBasicdependence)
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

					"comment": "rdk_test_comment_1",

					"project_name": "rdk_test_project_name_358",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "rdk_test_comment_1",

						"project_name": "rdk_test_project_name_358",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"comment": "update_test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "update_test",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDatahubProjectCheckmap = map[string]string{

	"comment": CHECKSET,

	"project_name": CHECKSET,

	"create_time": CHECKSET,

	"update_time": CHECKSET,

	"creator": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccDatahubProjectBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
