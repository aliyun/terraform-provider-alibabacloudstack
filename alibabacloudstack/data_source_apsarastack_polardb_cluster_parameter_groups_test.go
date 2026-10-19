package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterParameterGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_cluster_parameter_groups.default"
	name := fmt.Sprintf("tf_cluster_pg%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePolardbClusterParameterGroupsPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_cluster_parameter_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_polardb_cluster_parameter_group.default.id}_fake"},
		}),
	}
	nameprexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardb_cluster_parameter_group.default.parameter_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_polardb_cluster_parameter_group.default.parameter_group_name}_fake",
		}),
	}

	var existPolardbClusterParameterGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         CHECKSET,
			"ids.0":                         CHECKSET,
			"groups.#":                      CHECKSET,
			"groups.0.id":                   CHECKSET,
			"groups.0.parameter_group_id":   CHECKSET,
			"groups.0.parameter_group_name": CHECKSET,
			"groups.0.parameter_group_desc": CHECKSET,
			"groups.0.engine_version":       CHECKSET,
			"groups.0.engine":               CHECKSET,
			"groups.0.parameter_group_type": CHECKSET,
			"groups.0.force_restart":        CHECKSET,
			"groups.0.created":              CHECKSET,
			"groups.0.param_counts":         CHECKSET,
		}
	}

	var fakePolardbClusterParameterGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}

	var PolardbClusterParameterGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbClusterParameterGroupsMapFunc,
		fakeMapFunc:  fakePolardbClusterParameterGroupsMapFunc,
	}

	PolardbClusterParameterGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameprexConf)
}

func dataSourcePolardbClusterParameterGroupsPresetDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_polardb_cluster_parameter_group" "default" {
  engine = "MySQL"
  engine_version = "8.0"
  parameter_group_name = "${var.name}_test"
  parameter_group_desc = var.name
  parameters = {
    wait_timeout = "3600"
  }
}
`, name)
}
