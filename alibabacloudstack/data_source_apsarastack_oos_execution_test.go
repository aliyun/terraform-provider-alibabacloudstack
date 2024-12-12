package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackOosExecutionsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_executions.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
		}),
	}

	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}"]`,
			"category": `"${alibabacloudstack_oos_executions.default.Category}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"category": `"${alibabacloudstack_oos_executions.default.Category}_fake"`,
		}),
	}

	end_dateConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}"]`,
			"end_date": `"${alibabacloudstack_oos_executions.default.EndDate}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"end_date": `"${alibabacloudstack_oos_executions.default.EndDate}_fake"`,
		}),
	}

	executed_byConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_oos_executions.default.id}"]`,
			"executed_by": `"${alibabacloudstack_oos_executions.default.ExecutedBy}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"executed_by": `"${alibabacloudstack_oos_executions.default.ExecutedBy}_fake"`,
		}),
	}

	execution_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_oos_executions.default.id}"]`,
			"execution_id": `"${alibabacloudstack_oos_executions.default.ExecutionId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"execution_id": `"${alibabacloudstack_oos_executions.default.ExecutionId}_fake"`,
		}),
	}

	modeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_oos_executions.default.id}"]`,
			"mode": `"${alibabacloudstack_oos_executions.default.Mode}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"mode": `"${alibabacloudstack_oos_executions.default.Mode}_fake"`,
		}),
	}

	parent_execution_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_oos_executions.default.id}"]`,
			"parent_execution_id": `"${alibabacloudstack_oos_executions.default.ParentExecutionId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"parent_execution_id": `"${alibabacloudstack_oos_executions.default.ParentExecutionId}_fake"`,
		}),
	}

	ram_roleConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}"]`,
			"ram_role": `"${alibabacloudstack_oos_executions.default.RamRole}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"ram_role": `"${alibabacloudstack_oos_executions.default.RamRole}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_oos_executions.default.id}"]`,
			"status": `"${alibabacloudstack_oos_executions.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"status": `"${alibabacloudstack_oos_executions.default.Status}_fake"`,
		}),
	}

	template_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_oos_executions.default.id}"]`,
			"template_name": `"${alibabacloudstack_oos_executions.default.TemplateName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_oos_executions.default.id}_fake"]`,
			"template_name": `"${alibabacloudstack_oos_executions.default.TemplateName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_executions.default.id}"]`,

			"category":            `"${alibabacloudstack_oos_executions.default.Category}"`,
			"end_date":            `"${alibabacloudstack_oos_executions.default.EndDate}"`,
			"executed_by":         `"${alibabacloudstack_oos_executions.default.ExecutedBy}"`,
			"execution_id":        `"${alibabacloudstack_oos_executions.default.ExecutionId}"`,
			"mode":                `"${alibabacloudstack_oos_executions.default.Mode}"`,
			"parent_execution_id": `"${alibabacloudstack_oos_executions.default.ParentExecutionId}"`,
			"ram_role":            `"${alibabacloudstack_oos_executions.default.RamRole}"`,
			"status":              `"${alibabacloudstack_oos_executions.default.Status}"`,
			"template_name":       `"${alibabacloudstack_oos_executions.default.TemplateName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_executions.default.id}_fake"]`,

			"category":            `"${alibabacloudstack_oos_executions.default.Category}_fake"`,
			"end_date":            `"${alibabacloudstack_oos_executions.default.EndDate}_fake"`,
			"executed_by":         `"${alibabacloudstack_oos_executions.default.ExecutedBy}_fake"`,
			"execution_id":        `"${alibabacloudstack_oos_executions.default.ExecutionId}_fake"`,
			"mode":                `"${alibabacloudstack_oos_executions.default.Mode}_fake"`,
			"parent_execution_id": `"${alibabacloudstack_oos_executions.default.ParentExecutionId}_fake"`,
			"ram_role":            `"${alibabacloudstack_oos_executions.default.RamRole}_fake"`,
			"status":              `"${alibabacloudstack_oos_executions.default.Status}_fake"`,
			"template_name":       `"${alibabacloudstack_oos_executions.default.TemplateName}_fake"`}),
	}

	AlibabacloudstackOosExecutionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, categoryConf, end_dateConf, executed_byConf, execution_idConf, modeConf, parent_execution_idConf, ram_roleConf, statusConf, template_nameConf, allConf)
}

var existAlibabacloudstackOosExecutionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"executions.#":    "1",
		"executions.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackOosExecutionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"executions.#": "0",
	}
}

var AlibabacloudstackOosExecutionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_oos_executions.default",
	existMapFunc: existAlibabacloudstackOosExecutionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackOosExecutionsMapFunc,
}

func testAccCheckAlibabacloudstackOosExecutionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackOosExecutions%d"
}






data "alibabacloudstack_oos_executions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
