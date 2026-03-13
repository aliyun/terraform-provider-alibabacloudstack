package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackOosExecutionsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_oos_executions.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosExecutions-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOosExecutionsDependence)

	oneDay, _ := time.ParseDuration("24h")
	oneDayAfter := time.Now().Add(oneDay).Format("2006-01-02T15:04Z")
	oneDayDefore := time.Now().Add(-oneDay).Format("2006-01-02T15:04Z")

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oos_execution.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_oos_execution.default.id}"},
			"status": "Success",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_oos_execution.default.id}"},
			"status": "Cancelled",
		}),
	}
	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_oos_execution.default.id}"},
			"category": "Other",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":      []string{"${alibabacloudstack_oos_execution.default.id}"},
			"category": "TimerTrigger",
		}),
	}

	templateNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_oos_execution.default.id}"},
			"template_name": "${alibabacloudstack_oos_template.default.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_oos_execution.default.id}"},
			"template_name": "${alibabacloudstack_oos_template.default.template_name}-fake",
		}),
	}

	// FIXME: API has a bug, retry to test at 3.21.0
	/*
	executedByConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"executed_by": "${alibabacloudstack_oos_execution.default.executed_by}",
			"ids":         []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"executed_by": "${alibabacloudstack_oos_execution.default.executed_by}-fake",
			"ids":         []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
	}
	*/

	endDataBeforeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"end_date": oneDayAfter,
			"ids":      []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"end_date": oneDayDefore,
			"ids":      []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
	}

	endDataAfterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"end_date_after": oneDayDefore,
			"ids":            []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"end_date_after": oneDayAfter,
			"ids":            []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
	}

	startDataBeforeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"start_date_before": oneDayAfter,
			"ids":               []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"start_date_before": oneDayDefore,
			"ids":               []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
	}

	startDataAfterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"start_date_after": oneDayDefore,
			"ids":              []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"start_date_after": oneDayAfter,
			"ids":              []string{"${alibabacloudstack_oos_execution.default.id}"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_oos_execution.default.id}"},
			"status":        "Success",
			"template_name": "${alibabacloudstack_oos_template.default.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_oos_execution.default.id}"},
			"status":        "Cancelled",
			"template_name": "${alibabacloudstack_oos_template.default.template_name}",
		}),
	}
	var existOosExecutionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"executions.#":                  "1",
			"executions.0.category":         "Other",
			"executions.0.counters":         CHECKSET,
			"executions.0.create_date":      CHECKSET,
			"executions.0.end_date":         CHECKSET,
			"executions.0.executed_by":      CHECKSET,
			"executions.0.id":               CHECKSET,
			"executions.0.execution_id":     CHECKSET,
			"executions.0.is_parent":        "false",
			"executions.0.mode":             "Automatic",
			"executions.0.outputs":          CHECKSET,
			"executions.0.parameters":       CHECKSET,
			"executions.0.start_date":       CHECKSET,
			"executions.0.status":           "Success",
			"executions.0.template_id":      CHECKSET,
			"executions.0.template_name":    name,
			"executions.0.template_version": CHECKSET,
			"executions.0.update_date":      CHECKSET,
		}
	}

	var fakeOosExecutionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"executions.#": "0",
		}
	}

	var oosExecutionsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOosExecutionMapFunc,
		fakeMapFunc:  fakeOosExecutionMapFunc,
	}

	oosExecutionsInfo.dataSourceTestCheck(t, 0, statusConf, idsConf, categoryConf, endDataBeforeConf, endDataAfterConf, startDataBeforeConf, startDataAfterConf, templateNameConf, allConf)
}

func dataSourceOosExecutionsDependence(name string) string {
	return fmt.Sprintf(`
		resource "alibabacloudstack_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Describe instances of given status",
			"Parameters":{
			  "Status":{
				"Type": "String",
				"Description": "(Required) The status of the Ecs instance."
			  }
			},
			"Tasks": [
			  {
				"Properties" :{
				  "Parameters":{
					"Status": "{{ Status }}"
				  },
				  "API": "DescribeInstances",
				  "Service": "Ecs"
				},
				"Name": "foo",
				"Action": "ACS::ExecuteApi"
			  }]
		  }
		  EOF
		  template_name = "%s"
		  version_name = "test"
		}
		
		resource "alibabacloudstack_oos_execution" "default"{
			template_name = alibabacloudstack_oos_template.default.template_name
			description = "From TF Test"
			parameters = <<EOF
				{"Status":"Running"}
		  	EOF
		}
	`, name)
}
