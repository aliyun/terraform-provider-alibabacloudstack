package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_ess_scheduled_task", &resource.Sweeper{
		Name: "alibabacloudstack_ess_scheduled_task",
		F:    testSweepEssSchedules,
	})
}

func testSweepEssSchedules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ess.ScheduledTask
	req := ess.CreateDescribeScheduledTasksRequest()
	req.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.PageSize = requests.NewInteger(PageSizeLarge)

	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScheduledTasks(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Scheduled Tasks: %s", err)
		}
		resp, _ := raw.(*ess.DescribeScheduledTasksResponse)
		if resp == nil || len(resp.ScheduledTasks.ScheduledTask) < 1 {
			break
		}
		groups = append(groups, resp.ScheduledTasks.ScheduledTask...)

		if len(resp.ScheduledTasks.ScheduledTask) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	for _, v := range groups {
		name := v.ScheduledTaskName
		id := v.ScheduledTaskId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Scheduled Task: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Scheduled Task: %s (%s)", name, id)
		req := ess.CreateDeleteScheduledTaskRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ess", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		req.QueryParams["Department"] = client.Department
		req.QueryParams["ResourceGroup"] = client.ResourceGroup
		req.ScheduledTaskId = id
		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScheduledTask(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Scheduled Task (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAlibabacloudStackEssScheduledTask_basic(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alibabacloudstack_ess_scheduled_task.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerCommon + testAccEssScheduleConfig(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_action":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerCommon + testAccEssScheduleUpdateScheduledTaskName(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_task_name": fmt.Sprintf("tf-testAccEssSchedule-%d", rand),
					}),
				),
			},
			{
				Config: providerCommon + testAccEssScheduleUpdateDescription(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform test",
					}),
				),
			},
			{
				Config: providerCommon + testAccEssScheduleUpdateLaunchExpirationTime(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"launch_expiration_time": "500",
					}),
				),
			},
			{
				Config: providerCommon + testAccEssScheduleUpdateRecurrenceType(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence_type":     "Weekly",
						"recurrence_value":    CHECKSET,
						"recurrence_end_time": CHECKSET,
					}),
				),
			},
			{
				Config: providerCommon + testAccEssScheduleUpdateTaskEnabled(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_enabled": "false",
					}),
				),
			},
			{
				Config: providerCommon + testAccEssScheduleConfig(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEssScheduledTask_multi(t *testing.T) {
	var v ess.ScheduledTask
	resourceId := "alibabacloudstack_ess_scheduled_task.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	// Setting schedule time to more than one day
	oneDay, _ := time.ParseDuration("24h")
	rand := acctest.RandIntRange(1000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScheduledTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScheduleConfigMulti(EcsInstanceCommonTestCase,
					time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), rand),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduled_action":       CHECKSET,
						"launch_time":            CHECKSET,
						"scheduled_task_name":    fmt.Sprintf("tf-testAccEssScheduleConfig-%d-9", rand),
						"launch_expiration_time": "600",
						"task_enabled":           "true",
					}),
				),
			},
		},
	})
}

func testAccCheckEssScheduledTaskDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_scheduled_task" {
			continue
		}
		if _, err := essService.DescribeEssScheduledTask(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Schedule %s still exist", rs.Primary.ID)
	}

	return nil
}

func testAccEssScheduleConfig(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateScheduledTaskName(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateDescription(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
	}
	`, common, rand, scheduleTime)
}

func testAccEssScheduleUpdateLaunchExpirationTime(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
	}
	`, common, rand, scheduleTime)
}
func testAccEssScheduleUpdateRecurrenceType(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}

	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
		recurrence_type = "Weekly"
		recurrence_value = "0,1,2"
		recurrence_end_time = "%s"
	}
	`, common, rand, scheduleTime, scheduleTime)
}

func testAccEssScheduleUpdateTaskEnabled(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssSchedule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}"
		description = "terraform test"
		launch_expiration_time = 500
		//recurrence_type = "Weekly"
		//recurrence_value = "0,1,2"
		//recurrence_end_time = "%s"
		task_enabled = false
	}
	`, common, rand, scheduleTime, scheduleTime)
}
func testAccEssScheduleConfigMulti(common, scheduleTime string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScheduleConfig-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
		removal_policies = ["OldestInstance", "NewestInstance"]
	}
	
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 2
		cooldown = 60
	}
	
	resource "alibabacloudstack_ess_scheduled_task" "default" {
		count = 10
		scheduled_action = "${alibabacloudstack_ess_scaling_rule.default.ari}"
		launch_time = "%s"
		scheduled_task_name = "${var.name}-${count.index}"
	}
	`, common, rand, scheduleTime)
}
