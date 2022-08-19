package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackEhpcJobTemplatesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_ehpc_job_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${apsarastack_ehpc_job_template.default.id}_fake"]`,
		}),
	}

	var existApsaraStackEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"templates.#":                   "1",
			"templates.0.command_line":      "./LammpsTest/lammps.pbs",
			"templates.0.job_template_name": fmt.Sprintf("tf-testAccTemplates-%d", rand),
		}
	}
	var fakeApsaraStackEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var apsarastackEhpcJobTemplatesCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_ehpc_job_templates.default",
		existMapFunc: existApsaraStackEhpcJobTemplatesDataSourceNameMapFunc,
		fakeMapFunc:  fakeApsaraStackEhpcJobTemplatesDataSourceNameMapFunc,
	}
	apsarastackEhpcJobTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckApsaraStackEhpcJobTemplatesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTemplates-%d"
}

resource "apsarastack_ehpc_job_template" "default"{
  job_template_name =  var.name
  command_line=       "./LammpsTest/lammps.pbs"
}

data "apsarastack_ehpc_job_templates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
