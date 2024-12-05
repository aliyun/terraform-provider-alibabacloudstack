package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackRosTemplatesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ros_templates.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ros_templates.default.id}_fake"]`,
		}),
	}

	template_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ros_templates.default.id}"]`,
			"template_name": `"${alibabacloudstack_ros_templates.default.TemplateName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_ros_templates.default.id}_fake"]`,
			"template_name": `"${alibabacloudstack_ros_templates.default.TemplateName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ros_templates.default.id}"]`,

			"template_name": `"${alibabacloudstack_ros_templates.default.TemplateName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ros_templates.default.id}_fake"]`,

			"template_name": `"${alibabacloudstack_ros_templates.default.TemplateName}_fake"`}),
	}

	AlibabacloudstackRosTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf, template_nameConf, allConf)
}

var existAlibabacloudstackRosTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#":    "1",
		"templates.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRosTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#": "0",
	}
}

var AlibabacloudstackRosTemplatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ros_templates.default",
	existMapFunc: existAlibabacloudstackRosTemplatesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRosTemplatesMapFunc,
}

func testAccCheckAlibabacloudstackRosTemplatesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRosTemplates%d"
}






data "alibabacloudstack_ros_templates" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
