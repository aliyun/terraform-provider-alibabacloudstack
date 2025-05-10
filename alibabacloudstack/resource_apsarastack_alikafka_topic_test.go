package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_alikafka_topic", &resource.Sweeper{
		Name: "alibabacloudstack_alikafka_topic",
		F:    testSweepAlikafkaTopic,
	})
}

func testSweepAlikafkaTopic(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting alibabacloudstack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = region

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	var instanceIds []string
	for _, v := range instanceListResp.InstanceList.InstanceVO {
		instanceIds = append(instanceIds, v.InstanceId)
	}

	for _, instanceId := range instanceIds {

		// Control the topic list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafka.CreateGetTopicListRequest()
		request.InstanceId = instanceId
		request.RegionId = defaultRegionToTest
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka topics on instance (%s): %s", instanceId, err)
			continue
		}

		topicListResp, _ := raw.(*TopicListResponse)
		topics := topicListResp.TopicList

		for _, v := range topics {
			name := v.Topic
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka topic: %s ", name)
				continue
			}
			log.Printf("[INFO] delete alikafka topic: %s ", name)

			// Control the topic delete rate
			time.Sleep(time.Duration(400) * time.Millisecond)

			request := alikafka.CreateDeleteTopicRequest()
			request.InstanceId = instanceId
			request.Topic = v.Topic
			request.RegionId = defaultRegionToTest

			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.DeleteTopic(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka topic (%s): %s", name, err)
			}
		}
	}

	return nil
}

func TestAccAlibabacloudStackAlikafkaTopic_basic(t *testing.T) {

	var v *AliKafkaTopic
	resourceId := "alibabacloudstack_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, alikafkaTopicBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000,20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaTopicConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"instance_id":   "${alibabacloudstack_alikafka_instance.default.id}",
					"instance_id":   "${local.alikafka_instnace_id}",
					"topic":         "${var.name}",
					"local_topic":   "true",
					"compact_topic": "false",
					"partition_num": "12",
					"remark":        "alibabacloudstack_alikafka_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand),
						"local_topic":   "true",
						"compact_topic": "false",
						"partition_num": "12",
						"remark":        "alibabacloudstack_alikafka_topic_remark",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"compact_topic", "partition_num", "remark"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"partition_num": "24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "24",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic":  "tf-testacc-alibabacloudstack_alikafka_default_topic_change",
					"remark": "modified remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":  "tf-testacc-alibabacloudstack_alikafka_default_topic_change",
						"remark": "modified remark"}),
				),
			},

			// alibabacloudstack_alikafka_instance only support create post pay instance.
			// Post pay instance does not support create local or compact topic, so skip the following two test case temporarily.
			//{
			//	SkipFunc: shouldSkipLocalAndCompact("${alibabacloudstack_alikafka_instance.default.id}"),
			//	Config: testAccConfig(map[string]interface{}{
			//		"local_topic": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"local_topic": "true",
			//		}),
			//	),
			//},

			//{
			//	SkipFunc: shouldSkipLocalAndCompact("${alibabacloudstack_alikafka_instance.default.id}"),
			//	Config: testAccConfig(map[string]interface{}{
			//		"compact_topic": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"compact_topic": "true",
			//		}),
			//	),
			//},

			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "acceptance test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "acceptance test",
			// 		}),
			// 	),
			// },

			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "acceptance test",
			// 			"Updated": "TF",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "3",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "acceptance test",
			// 			"tags.Updated": "TF",
			// 		}),
			// 	),
			// },
		},
	})

}

/*
func TestAccAlibabacloudStackAlikafkaTopic_multi(t *testing.T) {

	var v *alikafka.InstanceDo
	resourceId := "alibabacloudstack_alikafka_topic.default.4"
	ra := resourceAttrInit(resourceId, alikafkaTopicBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaTopicConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":         "5",
					"instance_id":   "${alibabacloudstack_alikafka_instance.default.id}",
					"topic":         "${var.name}-${count.index}",
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "6",
					"remark":        "alibabacloudstack_alikafka_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v-4", rand),
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "6",
						"remark":        "alibabacloudstack_alikafka_topic_remark",
					}),
				),
			},
		},
	})

}
*/
func resourceAlikafkaTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

data "alibabacloudstack_alikafka_instances" "default" {
	ids = ["cluster-private-paas-default"]
}

resource "alibabacloudstack_alikafka_instance" "default" {
	count = length(data.alibabacloudstack_alikafka_instances.default.instances) > 0 ? 0 : 1
	name = "${var.name}"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	sasl = true
	plaintext = true
	spec = "Broker4C16G"
	
	provisioner "local-exec" {
		//防止broker未就绪导致的失败
		command = "sleep 300"
	}
}

locals{
alikafka_instnace_id = length(data.alibabacloudstack_alikafka_instances.default.instances) > 0 ? data.alibabacloudstack_alikafka_instances.default.instances.0.id : alibabacloudstack_alikafka_instance.default.0.id
}

`, name, DataZoneCommonTestCase)
}

var alikafkaTopicBasicMap = map[string]string{
	"topic":         "${var.name}",
	"local_topic":   "true",
	"compact_topic": "false",
	"partition_num": "12",
	"remark":        "alibabacloudstack_alikafka_topic_remark",
}
