package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackOnsGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.apsarastack_ons_groups.default"
	name := fmt.Sprintf("GID-tf-testacconsgroup%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${apsarastack_ons_instance.default.id}",
			"group_id_regex": "${apsarastack_ons_group.default.group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${apsarastack_ons_instance.default.id}",
			"group_id_regex": "${apsarastack_ons_group.default.group_id}_fake",
		}),
	}

	var existOnsGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                    "1",
			"groups.0.independent_naming": "true",
			"groups.0.remark":             "apsarastack_ons_group_remark",
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
variable "group_id" {
 default = "%v"
}

resource "apsarastack_ons_instance" "default" {
  name = var.group_id
  remark = "default-remark"
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
}

resource "apsarastack_ons_group" "default" {
  instance_id = "${apsarastack_ons_instance.default.id}"
  group_id = "${var.group_id}"
  remark = "apsarastack_ons_group_remark"
  read_enable = "true"
}
`, name)
}
