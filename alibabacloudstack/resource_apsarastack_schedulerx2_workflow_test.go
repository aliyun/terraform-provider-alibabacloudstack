package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSchedulerx2Workflow_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_schedulerx2_workflow.example"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &Schedulerx2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackSchedulerx2WorkflowDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":        "${alibabacloudstack_schedulerx2_app_group.example.group_id}",
					"name":            "${var.name}",
					"description":     "Initial description",
					"time_type":       "cron",
					"time_expression": "34 14 14 */1 * ?",
					"time_zone":       "PRC",
					"max_concurrency": "1",
					"enabled":         true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"description":     "Initial description",
						"time_type":       "cron",
						"time_expression": "34 14 14 */1 * ?",
						"max_concurrency": "1",
						"enabled":         "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"time_type":       "cron",
					"description":     "Updated description",
					"time_zone":       "Hongkong",
					"time_expression": "40 2 16 */1 * ?",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_type":       "cron",
						"description":     "Updated description",
						"time_zone":       "Hongkong",
						"time_expression": "40 2 16 */1 * ?",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}_updated",
					"description":     "Final description",
					"time_expression": REMOVEKEY,
					"time_type":       "api",
					"time_zone":       "GTM",
					"enabled":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name + "_updated",
						"description":     "Final description",
						"time_expression": REMOVEKEY,
						"time_type":       "api",
						"time_zone":       "GTM",
						"enabled":         "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"namespace", "time_zone"},
			},
		},
	})
}

func AlibabacloudStackSchedulerx2WorkflowDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}


resource "alibabacloudstack_schedulerx2_app_group" "example" {
  group_id         = "${var.name}.terra"
  app_name         = "${var.name}"
  description      = "${var.name}"
  max_jobs         = 30
  max_concurrency  = 10
  metrics_threshold {
    load5       = 10
    heap5_usage = 100
    disk_usage  = 100
  }
  monitor_config {
    send_channel = "mail,ding"
	alarm_type = "CustomContacts"
  }
}

`, name)
}
