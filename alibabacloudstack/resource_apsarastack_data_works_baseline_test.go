package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksBaseline_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_baseline.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksBaselineMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksBaseline")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_baseline%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksBaselineBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name":          "${var.name}",
					"project_id":             "${alibabacloudstack_data_works_project.default.id}",
					"owner":                  "${alibabacloudstack_data_works_user.default.user_id}",
					"priority":               "5",
					"baseline_type":          "DAILY",
					"alert_margin_threshold": "30",
					"enabled":                true,
					"alert_enabled":          true,
					"overtime_settings": []map[string]interface{}{
						{
							"cycle": 1,
							"time":  "08:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name":          name,
						"owner":                  CHECKSET,
						"priority":               "5",
						"baseline_type":          "DAILY",
						"alert_margin_threshold": "30",
						"enabled":                "true",
						"alert_enabled":          "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_name":          "${var.name}_update",
					"priority":               "3",
					"alert_margin_threshold": "60",
					"enabled":                "false",
					"alert_enabled":          "false",
					"overtime_settings": []map[string]interface{}{
						{
							"cycle": 1,
							"time":  "10:00",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_name":          name + "_update",
						"priority":               "3",
						"alert_margin_threshold": "60",
						"enabled":                "false",
						"alert_enabled":          "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_type":          "HOURLY",
					"alert_margin_threshold": "15",
					"overtime_settings": []map[string]interface{}{
						{
							"cycle": 1,
							"time":  "00:30",
						},
						{
							"cycle": 2,
							"time":  "01:45",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_type":          "HOURLY",
						"alert_margin_threshold": "15",
					}),
				),
			},
		},
	})
}

// AlibabacloudStackDataWorksBaselineMap0 is the expected state map for the DataWorks baseline resource
var AlibabacloudStackDataWorksBaselineMap0 = map[string]string{
	"project_id":  CHECKSET,
	"baseline_id": CHECKSET,
}

// AlibabacloudStackDataWorksBaselineBasicDependence0 returns the basic dependencies for the DataWorks baseline resource test
func AlibabacloudStackDataWorksBaselineBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alibabacloudstack_account" "current" {
}

data "alibabacloudstack_ascm_users" "default" {
 organization_id = "${data.alibabacloudstack_account.current.organization_id}"
}

resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
}

resource "alibabacloudstack_data_works_user" "default" {
	project_id= "${alibabacloudstack_data_works_project.default.id}"
	user_id=    "${data.alibabacloudstack_ascm_users.default.users.0.primary_key}"
	role_code = ["role_project_admin"]
	lifecycle {
	    ignore_changes = [
	      role_code
	    ]
	}
}
`, name)
}
