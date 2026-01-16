package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccalibabacloudstackEssAlarm(t *testing.T) {
	var v ess.Alarm
	rand := getAccTestRandInt(10000, 999999)
	var basicMap = map[string]string{
		// "name":                   fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand),
		// "description":            "Acc alarm test",
		// "alarm_actions.#":        "1",
		"scaling_group_id": CHECKSET,
		// "metric_type":            "system",
		// "evaluation_count":       "2",
		"cloud_monitor_group_id": NOSET,
		// "enable":                 "true",
		// "expressions.#":          "1",
	}
	resourceId := "alibabacloudstack_ess_alarm.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEssAlarm_basic-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEssAlarmConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                name,
					"description":         "Acc alarm test",
					"alarm_actions":       []string{"${alibabacloudstack_ess_scaling_rule.default.0.ari}"},
					"scaling_group_id":    "${alibabacloudstack_ess_scaling_group.default.id}",
					"metric_type":         "system",
					"evaluation_count":    "2",
					"period":              "300",
					"statistics":          "Average",
					"metric_name":         "CpuUtilization",
					"threshold":           "200.3",
					"comparison_operator": ">=",
					"enable":              true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                name,
						"description":         "Acc alarm test",
						"alarm_actions.#":     "1",
						"metric_type":         "system",
						"evaluation_count":    "2",
						"period":              "300",
						"statistics":          "Average",
						"metric_name":         "CpuUtilization",
						"threshold":           "200.3",
						"comparison_operator": ">=",
						"enable":              "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alarm_actions": "${alibabacloudstack_ess_scaling_rule.default.*.ari}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alarm_actions.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        name + "_update",
					"description": "Acc alarm test update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name + "_update",
						"description": "Acc alarm test update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period":              "120",
					"statistics":          "Minimum",
					"metric_name":         "MemoryUtilization",
					"threshold":           "99.9",
					"comparison_operator": ">",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period":              "120",
						"statistics":          "Minimum",
						"metric_name":         "MemoryUtilization",
						"threshold":           "99.9",
						"comparison_operator": ">",
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

func resourceEssAlarmConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	%s

	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}",]
	}

	resource "alibabacloudstack_ess_scaling_rule" "default" {
		count = 2
		scaling_rule_name = "${var.name}-${count.index}"
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}

`, name, VSwitchCommonTestCase)
}
