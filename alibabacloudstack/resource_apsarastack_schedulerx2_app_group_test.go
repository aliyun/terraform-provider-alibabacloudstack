package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func AlibabacloudStackSchedulerx2AppGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlibabacloudStackSchedulerx2AppGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_schedulerx2_app_group.example"
	ra := resourceAttrInit(resourceId, map[string]string{
		"app_key": CHECKSET,
	})
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &Schedulerx2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccschedulerx2AppGroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackSchedulerx2AppGroupDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_id":        "${var.name}.terra",
					"app_name":        "${var.name}",
					"description":     "${var.name}",
					"max_jobs":        "30",
					"max_concurrency": "10",
					"monitor_config": []map[string]interface{}{
						{
							"send_channel": "mail,ding",
							"alarm_type":   "CustomContacts",
						},
					},
					"contacts": []map[string]interface{}{
						{
							"username":    "test",
							"user_email":  "123@123.com",
							"dingding_ak": "testakkkkk",
						},
					},
					"metrics_threshold": []map[string]interface{}{
						{
							"load5":       10,
							"heap5_usage": 100,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id":                        fmt.Sprintf("%s.terra", name),
						"app_name":                        name,
						"description":                     name,
						"max_jobs":                        "30",
						"max_concurrency":                 "10",
						"monitor_config.#":                "1",
						"monitor_config.0.send_channel":   "mail,ding",
						"monitor_config.0.alarm_type":     "CustomContacts",
						"contacts.#":                      "1",
						"contacts.0.username":             "test",
						"contacts.0.user_email":           "123@123.com",
						"contacts.0.dingding_ak":          "testakkkkk",
						"metrics_threshold.#":             "1",
						"metrics_threshold.0.load5":       "10",
						"metrics_threshold.0.heap5_usage": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"app_name":        "${var.name}_update",
					"description":     "${var.name}_update",
					"max_jobs":        "20",
					"max_concurrency": "5",
					"monitor_config": []map[string]interface{}{
						{
							"send_channel": "mail",
							"alarm_type":   "CustomContacts",
						},
					},
					"contacts": []map[string]interface{}{
						{
							"username":    "test1",
							"user_email":  "123@333.com",
							"dingding_ak": "testakkkkkkkk",
						},
					},
					"metrics_threshold": []map[string]interface{}{
						{
							"load5":       5,
							"heap5_usage": 90,
							"disk_usage":  90,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"app_name":                        fmt.Sprintf("%s_update", name),
						"description":                     fmt.Sprintf("%s_update", name),
						"max_jobs":                        "20",
						"max_concurrency":                 "5",
						"monitor_config.0.send_channel":   "mail",
						"contacts.0.username":             "test1",
						"contacts.0.user_email":           "123@333.com",
						"contacts.0.dingding_ak":          "testakkkkkkkk",
						"metrics_threshold.#":             "1",
						"metrics_threshold.0.load5":       "5",
						"metrics_threshold.0.heap5_usage": "90",
						"metrics_threshold.0.disk_usage":  "90",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"namespace"},
			},
		},
	})
}
