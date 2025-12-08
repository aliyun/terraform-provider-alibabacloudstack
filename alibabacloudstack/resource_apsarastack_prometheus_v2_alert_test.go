package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPrometheusV2Alert_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_alert.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePrometheusV2Alert")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcekPrometheusV2AlertDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                 "${var.name}",
					"notify_recovered":     "true",
					"is_check_all":         "false",
					"tag_set":              []string{"aaa", "ccc"},
					"trigger_clusters":     []string{"${alibabacloudstack_prometheus_v2_instance.default.id}"},
					"trigger_promql":       "select testfield from testtable where testfield >= 0",
					"trigger_period":       "5m",
					"trigger_severity":     "warning",
					"trigger_cron":         "0 /5 * * * ?",
					"recover_notification": "Trigger condition\\\\$${alert_source} \\\\Hit record\\\\$${alert_time}",
					"notification":         "Trigger condition: {condition}\\nHit record :{alert_result}",
					"notify_types":         []string{"EMAIL", "SMS"},
					"notify_group_ids":     []string{"${alibabacloudstack_prometheus_v2_notify_group.default.id}"},
					"notify_interval":      "10m",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"notify_recovered":   "true",
						"is_check_all":       "false",
						"tag_set.#":          "2",
						"trigger_clusters.#": "1",
						"trigger_period":     "5m",
						"trigger_promql":     "select testfield from testtable where testfield >= 0",
						"trigger_severity":   "warning",
						"trigger_cron":       "0 /5 * * * ?",
						"notify_types.#":     "2",
						"notify_group_ids.#": "1",
						"notification":       "Trigger condition: {condition}\nHit record :{alert_result}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":             "${var.name}_updated",
					"notify_recovered": false,
					"trigger_period":   "5m",
					"trigger_promql":   "select testfield from testtable where testfield == 0",
					"notification":     "Trigger condition: {condition}\\nHit record :{alert_result}111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":             fmt.Sprintf("%s_updated", name),
						"notify_recovered": "false",
						"trigger_period":   "5m",
						"trigger_promql":   "select testfield from testtable where testfield == 0",
						"notification":     "Trigger condition: {condition}\nHit record :{alert_result}111",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourcekPrometheusV2AlertDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
  cluster_name = "${var.name}"
  tags         = ["test1", "test2"]
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  name        = "${var.name}_notify_group"
  type        = "WEBHOOK"
  description = "${var.name}_notify_group_description"
  webhook_url = "https://test.com"
  webhook_header_params {
    key   = "Content-Type"
    value = "application/json"
  }
}
 `, name)
}
