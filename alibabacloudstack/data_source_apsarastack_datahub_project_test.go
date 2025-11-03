package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDatahubProjectsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"name_regex":`"^${alibabacloudstack_datahub_project.default.name}$"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"name_regex": `"^${alibabacloudstack_datahub_project.default.name}_fake$"`,
		}),
	}

	AlibabacloudstackDatahubProjectsCheckInfo.dataSourceTestCheck(t, rand, nameConf)
}

var existAlibabacloudstackDatahubProjectsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":    "1",
		"projects.0.name": CHECKSET,
	}
}

var fakeAlibabacloudstackDatahubProjectsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#": "0",
	}
}

var AlibabacloudstackDatahubProjectsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_datahub_projects.default",
	existMapFunc: existAlibabacloudstackDatahubProjectsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackDatahubProjectsMapFunc,
}

func testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf_testdatahubrojects%d"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = var.name
}

data "alibabacloudstack_datahub_projects" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
