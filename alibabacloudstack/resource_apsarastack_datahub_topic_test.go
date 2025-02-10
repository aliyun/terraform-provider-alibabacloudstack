package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDatahubTopic0(t *testing.T) {
	var v *GetTopicResult

	resourceId := "alibabacloudstack_datahub_topic.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubTopicCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGettopicRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc_datahub_topic%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDatahubTopicBasicdependence)
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

					"comment": "test",

					"record_type": "BLOB",

					"project_name": "${alibabacloudstack_datahub_project.default.name}",

					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "test",

						"record_type": "BLOB",

						"project_name": name,

						"name": name,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDatahubTopicCheckmap = map[string]string{

	"comment": CHECKSET,

	"project_name": CHECKSET,

	"create_time": CHECKSET,

	"life_cycle": CHECKSET,

	"shard_count": CHECKSET,

	"name": CHECKSET,

	"record_type": CHECKSET,

}

func AlibabacloudTestAccDatahubTopicBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = var.name
}


`, name)
}
