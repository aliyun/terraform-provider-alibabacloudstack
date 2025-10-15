package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackAckTemplatesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ack_template.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_ack_template.default.name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_ack_template.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_ack_template.default.id}_fake" ]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_ack_template.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_ack_template.default.description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_ack_template.default.name}"`,
			"description_regex": `"${alibabacloudstack_ack_template.default.description}"`,
			"ids":               `[ "${alibabacloudstack_ack_template.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_ack_template.default.name}"`,
			"ids":               `[ "${alibabacloudstack_ack_template.default.id}" ]`,
			"description_regex": `"${alibabacloudstack_ack_template.default.description}_fake"`,
		}),
	}

	AckTemplatesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStackAckTemplatesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccAckTemplatesDatasource%d"
}

resource "alibabacloudstack_ack_template" "default" {
    template = "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  labels:\n    vsw: test\n  name: nginx-deployment-basic\n  namespace: default\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      vsw: test\n  template:\n    metadata:\n      labels:\n        vsw: test\n    spec:\n      containers:\n        - command:\n            - sleep\n            - '1000000'\n          image: >-\n            registry.acs.intra.env35.shuguang.com/acs/busybox:1.33.1\n          imagePullPolicy: IfNotPresent\n          name: vsw"
	name = "${var.name}"
	description="${var.name}"
	template_type="kubernetes"
}

data "alibabacloudstack_ack_templates" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existAckTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#":                            "1",
		"ids.#":                                  "1",
		"templates.0.template":                   CHECKSET,
		"templates.0.description":                CHECKSET,
		"templates.0.name":                       CHECKSET,
		"templates.0.template_id":                CHECKSET,
		"templates.0.template_type":              CHECKSET,
		"templates.0.template_with_hist_id":      CHECKSET,
		"templates.0.template_hash_code_version": CHECKSET,
		"templates.0.created":                    CHECKSET,
		"templates.0.acl":                        CHECKSET,
		"templates.0.version":                    CHECKSET,
		"templates.0.ali_uid":                    CHECKSET,
		"templates.0.updated":                    CHECKSET,
	}
}

var fakeAckTemplatesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"templates.#": "0",
		"ids.#":       "0",
	}
}

var AckTemplatesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ack_templates.default",
	existMapFunc: existAckTemplatesMapFunc,
	fakeMapFunc:  fakeAckTemplatesMapFunc,
}
