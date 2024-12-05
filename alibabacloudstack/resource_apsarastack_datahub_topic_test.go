package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDatahubTopic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_datahub_topic.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDatahubTopicCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDatahubGettopicRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdata_hubtopic%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDatahubTopicBasicdependence)
	resource.Test(t, resource.TestCase{
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

					"project_name": "${{ref(resource, DataHub::Project::4.0.0::XviFXs.resourceAttribute.ProjectName)}}",

					"expand_mode": "true",

					"topic_name": "rdk_test_topic_name_355",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"comment": "test",

						"record_type": "BLOB",

						"project_name": "${{ref(resource, DataHub::Project::4.0.0::XviFXs.resourceAttribute.ProjectName)}}",

						"expand_mode": "true",

						"topic_name": "rdk_test_topic_name_355",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDatahubTopicCheckmap = map[string]string{

	"comment": CHECKSET,

	"enable_schema_registry": CHECKSET,

	"project_name": CHECKSET,

	"create_time": CHECKSET,

	"expand_mode": CHECKSET,

	"lifecycle": CHECKSET,

	"creator": CHECKSET,

	"shard_count": CHECKSET,

	"topic_name": CHECKSET,

	"total_count": CHECKSET,

	"storage": CHECKSET,

	"record_type": CHECKSET,

	"update_time": CHECKSET,

	"record_schema": CHECKSET,
}

func AlibabacloudTestAccDatahubTopicBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}




resource "_data_hub_project" "XviFXs" {


    


    
    "comment" : "test",











    


    
    "project_name" : "dependency_test",










}


`, name)
}
