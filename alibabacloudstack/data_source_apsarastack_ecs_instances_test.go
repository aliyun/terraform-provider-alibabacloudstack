package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackInstancesDataSourceBasic(t *testing.T) {
	resourceId := "data.alibabacloudstack_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "Tf-EcsInstanceDataSource", dataSourceEcsInstancesConfigDependence)

	// Test with ids filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"i-fake-nonexistent"},
		}),
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^Tf-EcsInstanceDataSource_ecs$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^nonexistent-ecs-instance-name$",
		}),
	}

	var existInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         CHECKSET,
			"instances.#":                   CHECKSET,
			"instances.0.id":                CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.availability_zone": CHECKSET,
		}
	}

	var fakeInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var instancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstancesMapFunc,
		fakeMapFunc:  fakeInstancesMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "GET", "Ecs", "2014-05-26", "DescribeInstances", "")
		},
	}
	instancesCheckInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf)
}

func dataSourceEcsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s
`, name, ECSInstanceCommonTestCase)
}
