package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEhpcJobTemplatesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ehpc_job_templates.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccTemplates-%d", rand),
		dataSourceEhpcJobTemplatesConfigDependence)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"templates.#":                   "1",
			"templates.0.command_line":      "./LammpsTest/lammps.pbs",
			"templates.0.job_template_name": fmt.Sprintf("tf-testAccTemplates-%d", rand),
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_ehpc_job_template.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_ehpc_job_template.default.id}_fake"},
			}),
		},
	)
}

func dataSourceEhpcJobTemplatesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ehpc_job_template" "default" {
	job_template_name = var.name
	command_line      = "./LammpsTest/lammps.pbs"
}
`, name)
}
