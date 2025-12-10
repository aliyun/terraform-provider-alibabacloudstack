package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func buildBasicSchedulerx2AppGroup(name string) string {
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
}
`, name)
}

func TestAccAlibabacloudStackSchedulerx2Job_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_schedulerx2_job.example"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &Schedulerx2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-test%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildBasicSchedulerx2AppGroup)

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
					"execute_mode":                "Standalone",
					"priority":                    "1",
					"name":                        "${var.name}",
					"group_id":                    "${alibabacloudstack_schedulerx2_app_group.example.group_id}",
					"description":                 "ddddddd",
					"job_type":                    "java",
					"parameters":                  "testargs=1",
					"max_attempt":                 "3",
					"attempt_interval":            "30",
					"max_concurrency":             "2",
					"time_type":                   "1",
					"time_expression":             "8 59 15 */1 * ?",
					"content":                     "{\"className\":\"Create\"}",
					"monitor_timeout_enable":      "true",
					"monitor_timeout_kill_enable": "true",
					"monitor_fail_enable":         "true",
					"monitor_miss_worker_enable":  "true",
					"monitor_timeout":             "3600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name,
						"execute_mode":                "Standalone",
						"priority":                    "1",
						"description":                 "ddddddd",
						"job_type":                    "java",
						"parameters":                  "testargs=1",
						"max_attempt":                 "3",
						"attempt_interval":            "30",
						"max_concurrency":             "2",
						"time_type":                   "1",
						"time_expression":             "8 59 15 */1 * ?",
						"content":                     "{\"className\":\"Create\"}",
						"monitor_timeout_enable":      "true",
						"monitor_timeout_kill_enable": "true",
						"monitor_fail_enable":         "true",
						"monitor_miss_worker_enable":  "true",
						"monitor_timeout":             "3600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"execute_mode":                "Standalone",
					"priority":                    "5",
					"name":                        "${var.name}_updated",
					"group_id":                    "${alibabacloudstack_schedulerx2_app_group.example.group_id}",
					"description":                 "updated description",
					"job_type":                    "java",
					"parameters":                  "testargs=1 updated",
					"max_attempt":                 "5",
					"attempt_interval":            "60",
					"max_concurrency":             "3",
					"time_type":                   "1",
					"time_expression":             "0 0 12 */1 * ?",
					"content":                     "{\"className\":\"Updated\"}",
					"monitor_timeout_enable":      "false",
					"monitor_timeout_kill_enable": "false",
					"monitor_fail_enable":         "false",
					"monitor_miss_worker_enable":  "false",
					"monitor_timeout":             "7200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                        name + "_updated",
						"execute_mode":                "Standalone",
						"priority":                    "5",
						"description":                 "updated description",
						"job_type":                    "java",
						"parameters":                  "testargs=1 updated",
						"max_attempt":                 "5",
						"attempt_interval":            "60",
						"max_concurrency":             "3",
						"time_type":                   "1",
						"time_expression":             "0 0 12 */1 * ?",
						"content":                     "{\"className\":\"Updated\"}",
						"monitor_timeout_enable":      "false",
						"monitor_timeout_kill_enable": "false",
						"monitor_fail_enable":         "false",
						"monitor_miss_worker_enable":  "false",
						"monitor_timeout":             "7200",
					}),
				),
			},
		},
	})
}
