package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOnsGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ons_groups.default"
	name := fmt.Sprintf("tf-groupdata%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_ons_instance.default.id}",
			"group_id_regex": "${alibabacloudstack_ons_group.default.group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alibabacloudstack_ons_instance.default.id}",
			"group_id_regex": "${alibabacloudstack_ons_group.default.group_id}_fake",
		}),
	}

	var existOnsGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                    "1",
			"groups.0.independent_naming": "true",
			"groups.0.remark":             "alibabacloudstack_ons_group_remark",
		}
	}

	var fakeOnsGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
		}
	}

	var onsGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsGroupsMapFunc,
		fakeMapFunc:  fakeOnsGroupsMapFunc,
	}

	onsGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceOnsGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

%s

resource "alibabacloudstack_ons_group" "default" {
  instance_id = "${alibabacloudstack_ons_instance.default.id}"
  group_id = "GID-${var.name}"
  remark = "alibabacloudstack_ons_group_remark"
}
`, name, OnsCommonTestCase)
}
