package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksRemind_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_remind.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksRemindMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksRemind")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksremind%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksRemindBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_unit":    "OWNER",
					"remind_name":   name,
					"remind_type":   "FINISHED",
					"remind_unit":   "PROJECT",
					"project_id":    "10023",
					"alert_methods": "SMS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_unit":    "OWNER",
						"remind_name":   name,
						"remind_type":   "FINISHED",
						"remind_unit":   "PROJECT",
						"project_id":    "10023",
						"alert_methods": "SMS",
					}),
				),
			},
			//测试时, 只留一个 remind_unit 类型单独测试，其余注释掉
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_unit":  "BASELINE",
					"baseline_ids": "100016",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_unit":  "BASELINE",
						"baseline_ids": "100016",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_type": "TIMEOUT",
					"detail":      "1800",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "TIMEOUT",
						"detail":      "1800",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_type": "ERROR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "ERROR",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_type": "UNFINISHED",
					"detail": fmt.Sprintf("{\\%s\":%d,\\%s\":%d}",
						"\"hour\\", 23, "\"minu\\", 59,
					),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "UNFINISHED",
						"detail":      "{\"hour\":23,\"minu\":59}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_type": "CYCLE_UNFINISHED",
					"detail": fmt.Sprintf("{\\%s\":\\%s\",\\%s\":\\%s\"}",
						"\"1\\", "\"05:50\\", "\"2\\", "\"06:50\\",
					),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "CYCLE_UNFINISHED",
						"detail":      "{\"1\":\"05:50\",\"2\":\"06:50\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_flag": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_flag": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"use_flag": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"use_flag": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dnd_end":         "08:00",
					"alert_interval":  "1200",
					"max_alert_times": "4",
					"alert_methods":   "MAIL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dnd_end":         "08:00",
						"alert_interval":  "1200",
						"max_alert_times": "4",
						"alert_methods":   "MAIL",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDataWorksRemindMap0 = map[string]string{}

func AlibabacloudStackDataWorksRemindBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
