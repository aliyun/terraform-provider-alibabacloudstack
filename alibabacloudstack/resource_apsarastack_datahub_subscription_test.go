package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDatahubSubscription_basic(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alibabacloudstack_datahub_subscription.default"
	ra := resourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alibabacloudstack_datahub_project.basic.name}",
					"topic_name":   "${alibabacloudstack_datahub_topic.basic.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": "subscription for basic.",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "subscription for basic.",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"comment": REMOVEKEY,
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"comment": "subscription added by terraform",
			//		}),
			//	),
			//},
		},
	})
}
func TestAccAlibabacloudStackDatahubSubscription_multi(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alibabacloudstack_datahub_subscription.default.4"
	ra := resourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alibabacloudstack_datahub_project.basic.name}",
					"topic_name":   "${alibabacloudstack_datahub_topic.basic.name}",
					"count":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourceDatahubSubscriptionConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "project_name" {
	  default = "%s"
	}
	variable "topic_name" {
	  default = "tf_testacc_datahub_topic"
	}
	variable "record_type" {
	  default = "BLOB"
	}
	resource "alibabacloudstack_datahub_project" "basic" {
	  name = "${var.project_name}"
	  comment = "project for basic."
	}
	resource "alibabacloudstack_datahub_topic" "basic" {
	  project_name = "${alibabacloudstack_datahub_project.basic.name}"
	  name = "${var.topic_name}"
	  record_type = "${var.record_type}"
	  shard_count = 3
	  life_cycle = 7
	  comment = "topic for basic."
	}
	`, name)
}

var datahubSubscriptionBasicMap = map[string]string{
	"project_name": CHECKSET,
	"topic_name":   CHECKSET,
	//"comment":      "subscription added by terraform",
}
