package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackOosTemplatesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
		}),
	}

	create_timeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_oos_templates.default.id}"]`,
			"create_time": `"${alibabacloudstack_oos_templates.default.CreateTime}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
			"create_time": `"${alibabacloudstack_oos_templates.default.CreateTime}_fake"`,
		}),
	}

	created_byConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_oos_templates.default.id}"]`,
			"created_by": `"${alibabacloudstack_oos_templates.default.CreatedBy}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
			"created_by": `"${alibabacloudstack_oos_templates.default.CreatedBy}_fake"`,
		}),
	}

	has_triggerConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
		}),
	}

	share_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_oos_templates.default.id}"]`,
			"share_type": `"${alibabacloudstack_oos_templates.default.ShareType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
			"share_type": `"${alibabacloudstack_oos_templates.default.ShareType}_fake"`,
		}),
	}

	template_formatConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_oos_templates.default.id}"]`,
			"template_format": `"${alibabacloudstack_oos_templates.default.TemplateFormat}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
			"template_format": `"${alibabacloudstack_oos_templates.default.TemplateFormat}_fake"`,
		}),
	}

	template_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_oos_templates.default.id}"]`,
			"template_name": `"${alibabacloudstack_oos_templates.default.TemplateName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_oos_templates.default.id}_fake"]`,
			"template_name": `"${alibabacloudstack_oos_templates.default.TemplateName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}"]`,

			"create_time":     `"${alibabacloudstack_oos_templates.default.CreateTime}"`,
			"created_by":      `"${alibabacloudstack_oos_templates.default.CreatedBy}"`,
			"share_type":      `"${alibabacloudstack_oos_templates.default.ShareType}"`,
			"template_format": `"${alibabacloudstack_oos_templates.default.TemplateFormat}"`,
			"template_name":   `"${alibabacloudstack_oos_templates.default.TemplateName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_oos_templates.default.id}_fake"]`,

			"create_time":     `"${alibabacloudstack_oos_templates.default.CreateTime}_fake"`,
			"created_by":      `"${alibabacloudstack_oos_templates.default.CreatedBy}_fake"`,
			"share_type":      `"${alibabacloudstack_oos_templates.default.ShareType}_fake"`,
			"template_format": `"${alibabacloudstack_oos_templates.default.TemplateFormat}_fake"`,
			"template_name":   `"${alibabacloudstack_oos_templates.default.TemplateName}_fake"`}),
	}

	AlibabacloudstackOosTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, create_timeConf, created_byConf, has_triggerConf, share_typeConf, template_formatConf, template_nameConf, allConf)
}

var existAlibabacloudstackOosTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#":    "1",
		"templates.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackOosTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#": "0",
	}
}

var AlibabacloudstackOosTemplatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_oos_templates.default",
	existMapFunc: existAlibabacloudstackOosTemplatesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackOosTemplatesMapFunc,
}

func testAccCheckAlibabacloudstackOosTemplatesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackOosTemplates%d"
}






data "alibabacloudstack_oos_templates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
