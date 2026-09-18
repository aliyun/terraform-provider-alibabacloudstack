package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDatahubProjectsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_datahub_projects.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testdatahubrojects%d", rand),
		dataSourceDatahubProjectsConfigDependence)

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_datahub_project.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_datahub_project.default.name}_fake$",
		}),
	}

	AlibabacloudstackDatahubProjectsCheckInfo.dataSourceTestCheck(t, rand, nameConf)
}

var existAlibabacloudstackDatahubProjectsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"projects.#":      "1",
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

func dataSourceDatahubProjectsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = var.name
}
`, name)
}
