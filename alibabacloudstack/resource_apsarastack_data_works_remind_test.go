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
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksRemind")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_remind%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksRemindBasicDependence0)
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
					"alert_unit":    "OWNER",
					"remind_name":   name,
					"remind_type":   "FINISHED",
					"remind_unit":   "PROJECT",
					"project_id":    "${alibabacloudstack_data_works_project.default.id}",
					"alert_methods": []string{"SMS"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_unit":    "OWNER",
						"remind_name":   name,
						"remind_type":   "FINISHED",
						"remind_unit":   "PROJECT",
						"project_id":    CHECKSET,
					}),
				),
			},
			// During testing, only keep one remind_unit type for individual testing, comment out the rest
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":   REMOVEKEY,
					"remind_unit":  "BASELINE",
					"baseline_ids": []string{"${alibabacloudstack_data_works_baseline.default.baseline_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":   REMOVEKEY,
						"remind_unit":  "BASELINE",
						"baseline_ids": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"baseline_ids": REMOVEKEY,
					"remind_unit":  "NODE",
					"node_ids":     []string{"${alibabacloudstack_data_works_file.default.file_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"baseline_ids": REMOVEKEY,
						"remind_unit":  "NODE",
						"node_ids":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_ids":        REMOVEKEY,
					"remind_unit":     "BIZPROCESS",
					"biz_process_ids": []string{"${alibabacloudstack_data_works_business.default.business_id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_ids":        REMOVEKEY,
						"remind_unit":     "BIZPROCESS",
						"biz_process_ids": CHECKSET,
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
					"detail": TfRawString(`jsonencode({
	hour = 23
	minu = 59
})`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "UNFINISHED",
						"detail":      `{"hour":23,"minu":59}`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remind_type": "CYCLE_UNFINISHED",
					"detail": TfRawString(`jsonencode({
						"1" = "05:50"
						"2" = "06:50"
					})`),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remind_type": "CYCLE_UNFINISHED",
						"detail":      `{"1":"05:50","2":"06:50"}`,
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
					"alert_methods":   []string{"SMS","MAIL"},
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

data "alibabacloudstack_account" "current" {
}

data "alibabacloudstack_ascm_users" "default" {
  organization_id = data.alibabacloudstack_account.current.organization_id
}

resource "alibabacloudstack_data_works_project" "default" {
  name           = var.name
  description    = "${var.name}_desc"
  task_auth_type = "PROJECT"
}

resource "alibabacloudstack_data_works_user" "default" {
  project_id = alibabacloudstack_data_works_project.default.id
  user_id    = data.alibabacloudstack_ascm_users.default.users.0.primary_key
  role_code  = ["role_project_admin"]
  lifecycle {
    ignore_changes = [
      role_code
    ]
  }
}


resource "alibabacloudstack_data_works_baseline" "default" {
  alert_enabled          = true
  project_id             = alibabacloudstack_data_works_project.default.id
  owner                  = alibabacloudstack_data_works_user.default.user_id
  alert_margin_threshold = "30"
  baseline_name          = var.name
  enabled                = true
  overtime_settings {
    time  = "08:00"
    cycle = 1
  }

  priority      = "5"
  baseline_type = "DAILY"
}

resource "alibabacloudstack_data_works_business" "default" {
	project_id=  "${alibabacloudstack_data_works_project.default.id}"
	name=        "${var.name}"
	description= "${var.name}_desc"
}

resource "alibabacloudstack_data_works_folder" "default" {
	project_id    = "${alibabacloudstack_data_works_project.default.id}"
	business_name = "${alibabacloudstack_data_works_business.default.name}"
	engine_type   = "General"
	folder_path   = "folder_test"
}

data "alibabacloudstack_dataworks_file_types" "default" {
	name = "Shell"
	project_id = "${alibabacloudstack_data_works_project.default.id}"
}

resource "alibabacloudstack_data_works_file" "default" {
  project_id = "${alibabacloudstack_data_works_project.default.id}"
  folder_id = "${alibabacloudstack_data_works_folder.default.folder_id}"
  file_name = "${var.name}"
  file_type = "${data.alibabacloudstack_dataworks_file_types.default.file_types.0.node_type_id}"
}

`, name)
}
