package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
)

func TestAccAlibabacloudStackDatahubProject0(t *testing.T) {
	var v *datahub.GetProjectResult

	resourceId := "alibabacloudstack_datahub_project.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubProjectCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGetkafkagroupRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)

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

					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "rdk_test_comment_1",

						"name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			// 			{
			// 				Config: testAccConfig(map[string]interface{}{
			//
			// 					"comment": "update_test",
			// 				}),
			// 				Check: resource.ComposeTestCheckFunc(
			// 					testAccCheck(map[string]string{
			//
			// 						"comment": "update_test",
			// 					}),
			// 				),
			// 			},
		},
	})
}

var AlibabacloudTestAccDatahubProjectCheckmap = map[string]string{

	"comment": CHECKSET,

	"name": CHECKSET,

	"create_time": CHECKSET,
}

func AlibabacloudTestAccDatahubProjectBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
