package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackDatahubProjectsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_projects.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_projects.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_projects.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackDatahubProjectsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_datahub_projects.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackDatahubProjectsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackDatahubProjectsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":    "1",
		"projects.0.id": CHECKSET,
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
	default = "tf-testAlibabacloudstackDatahubProjects%d"
}






data "alibabacloudstack_datahub_projects" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
