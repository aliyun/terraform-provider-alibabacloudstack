package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackOosTemplatesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_oos_templates.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccOosTemplate-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOosTemplatesDependence)

	oneDay, _ := time.ParseDuration("24h")
	oneDayAfter := time.Now().Add(oneDay).Format("2006-01-02T15:04Z")
	oneDayDefore := time.Now().Add(-oneDay).Format("2006-01-02T15:04Z")

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oos_template.default.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oos_template.default.template_name}-fake",
		}),
	}
	
	formatConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_format": "JSON",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"template_format": "YAML",
		}),
	}

	createByConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"created_by": "${alibabacloudstack_oos_template.default.created_by}",
			"ids":        []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"created_by": "${alibabacloudstack_oos_template.default.created_by}-fake",
			"ids":        []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
	}

	createDataBeforeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"created_date": oneDayAfter,
			"ids":          []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"created_date": oneDayDefore,
			"ids":          []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
	}

	createDataAfterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"created_date_after": oneDayDefore,
			"ids":                []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"created_date_after": oneDayAfter,
			"ids":                []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oos_template.default.template_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oos_template.default.template_name}-fake"},
		}),
	}
	shareTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"share_type": "Private",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"share_type": "Public",
		}),
	}
	hasTriggerConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"has_trigger": "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"has_trigger": "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"has_trigger": "false",
			"share_type":  "Private",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oos_template.default.template_name}"},
			"has_trigger": "false",
			"share_type":  "Public",
		}),
	}
	var existOosTemplateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"ids.0":                        CHECKSET,
			"templates.#":                  "1",
			"templates.0.category":         CHECKSET,
			"templates.0.created_date":     CHECKSET,
			"templates.0.description":      CHECKSET,
			"templates.0.has_trigger":      CHECKSET,
			"templates.0.created_by":       CHECKSET,
			"templates.0.share_type":       "Private",
			"templates.0.template_format":  "JSON",
			"templates.0.template_id":      CHECKSET,
			"templates.0.id":               name,
			"templates.0.template_name":    name,
			"templates.0.template_version": CHECKSET,
			"templates.0.updated_by":       CHECKSET,
			"templates.0.updated_date":     CHECKSET,
		}
	}

	var fakeOosTemplateMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"templates.#": "0",
		}
	}

	var oosTemplatesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOosTemplateMapFunc,
		fakeMapFunc:  fakeOosTemplateMapFunc,
	}

	oosTemplatesInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, formatConf, createByConf, createDataBeforeConf, createDataAfterConf, shareTypeConf, hasTriggerConf, allConf)
}

func dataSourceOosTemplatesDependence(name string) string {
	return fmt.Sprintf(`
		resource "alibabacloudstack_oos_template" "default" {
		  content= <<EOF
		  {
			"FormatVersion": "OOS-2019-06-01",
			"Description": "Update Describe instances of given status",
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
	`, name)
}
