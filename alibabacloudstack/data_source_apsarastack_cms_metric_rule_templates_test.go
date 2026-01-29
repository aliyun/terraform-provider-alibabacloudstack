package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCmsMetricRuleTemplatesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cms_metric_rule_templates.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacccn-cms-metricruletemplate%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlibabacloudStackCmsMetricRuleTemplatesDepenced)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_metric_rule_template.default.metric_rule_template_name}",
			"is_default": false,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
			"is_default": false,
		}),
	}

	templateIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"template_id": "${alibabacloudstack_cms_metric_rule_template.default.id}",
			"is_default":  false,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"template_id": "999999999",
			"is_default":  false,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_cms_metric_rule_template.default.id}"},
			"is_default": false,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"999999999"},
			"is_default": false,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "${alibabacloudstack_cms_metric_rule_template.default.metric_rule_template_name}",
			"template_id": "${alibabacloudstack_cms_metric_rule_template.default.id}",
			"ids":         []string{"${alibabacloudstack_cms_metric_rule_template.default.id}"},
			"is_default":  false,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "fake_*",
			"template_id": "999999999",
			"ids":         []string{"999999999"},
			"is_default":  false,
		}),
	}

	var existCmsMetricRuleTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"templates.#":              "1",
			"templates.0.name":         name,
			"templates.0.id":           CHECKSET,
			"templates.0.description":  name,
			"templates.0.rest_version": CHECKSET,
		}
	}

	var fakeCmsMetricRuleTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"templates.#": "0",
		}
	}

	var cmsMetricRuleTemplatesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCmsMetricRuleTemplatesMapFunc,
		fakeMapFunc:  fakeCmsMetricRuleTemplatesMapFunc,
	}

	cmsMetricRuleTemplatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, templateIdConf, idsConf, allConf)
}

func dataSourceAlibabacloudStackCmsMetricRuleTemplatesDepenced(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_cms_metric_rule_template" "default" {
  metric_rule_template_name = var.name
  alert_templates {
    category      = "ecs"
    metric_name   = "cpu_total"
    namespace     = "acs_ecs_dashboard"
    rule_name     = "tf_testAcc_new"
    escalations {
      critical {
        comparison_operator = "GreaterThanThreshold"
        statistics          = "Average"
        threshold           = "90"
        times               = "3"
      }
    }
  }
  description = var.name
}
`, name)
}
