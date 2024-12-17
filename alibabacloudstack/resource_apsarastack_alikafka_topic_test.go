package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAlikafkaTopic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAlikafkaTopicCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoAlikafkaGettopiclistRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sali_kafkatopic%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAlikafkaTopicBasicdependence)
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

					"instance_id": "alikafka_post-cn-9lb33b12l004",

					"region_id": "cn-hangzhou",

					"topic": "rdktest",

					"remark": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_id": "alikafka_post-cn-9lb33b12l004",

						"region_id": "cn-hangzhou",

						"topic": "rdktest",

						"remark": "test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"region_id": "cn-hangzhou",

					"remark": "update",

					"instance_id": "alikafka_post-cn-9lb33b12l004",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-hangzhou",

						"remark": "update",

						"instance_id": "alikafka_post-cn-9lb33b12l004",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccAlikafkaTopicCheckmap = map[string]string{

	"instance_id": CHECKSET,

	"remark": CHECKSET,

	"partition_num": CHECKSET,

	"region_id": CHECKSET,

	"topic": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccAlikafkaTopicBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
