package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackEhpcJobTemplatesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_ehpc_job_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_ehpc_job_template.default.id}_fake"]`,
		}),
	}

	var existAlibabacloudStackEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"templates.#":                   "1",
			"templates.0.command_line":      "./LammpsTest/lammps.pbs",
			"templates.0.job_template_name": fmt.Sprintf("tf-testAccTemplates-%d", rand),
		}
	}
	var fakeAlibabacloudStackEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alibabacloudstackEhpcJobTemplatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_ehpc_job_templates.default",
		existMapFunc: existAlibabacloudStackEhpcJobTemplatesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackEhpcJobTemplatesDataSourceNameMapFunc,
	}
	alibabacloudstackEhpcJobTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckAlibabacloudStackEhpcJobTemplatesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTemplates-%d"
}

resource "alibabacloudstack_ehpc_job_template" "default"{
  job_template_name =  var.name
  command_line=       "./LammpsTest/lammps.pbs"
}

data "alibabacloudstack_ehpc_job_templates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
