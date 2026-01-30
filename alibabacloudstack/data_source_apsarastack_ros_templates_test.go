package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackRosTemplatesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_ros_templates.default"
	name := fmt.Sprintf("tf-test-ros-template-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplatesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ros_template.basic.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-template-id-12345"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ros_template.basic.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ros_template.basic.template_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ros_template.basic.id}"},
			"name_regex": "${alibabacloudstack_ros_template.basic.template_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"fake-template-id-12345"},
			"name_regex": "${alibabacloudstack_ros_template.basic.template_name}_fake",
		}),
	}

	var existAlibabacloudstackRosTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"templates.#":               "1",
			"templates.0.id":            CHECKSET,
			"templates.0.template_name": name,
			"templates.0.description":   "Created by Terraform AccTest",
			"templates.0.create_time":   CHECKSET,
		}
	}

	var fakeAlibabacloudstackRosTemplatesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"templates.#": "0",
		}
	}

	var rosTemplatesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlibabacloudstackRosTemplatesMapFunc,
		fakeMapFunc:  fakeAlibabacloudstackRosTemplatesMapFunc,
	}

	rosTemplatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceRosTemplatesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_ros_template" "basic" {
  template_name = var.name
  description   = "Created by Terraform AccTest"
  template_body = jsonencode({
    ROSTemplateFormatVersion = "2015-09-01"
    Description              = "Basic empty template"
    Resources                = {}
  })
}
`, name)
}
