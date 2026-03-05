package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOnsInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ons_instances.default"
	name := fmt.Sprintf("tf-instancedata%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ons_instance.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-name",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ons_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id-12345"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ons_instance.default.id}"},
			"name_regex": "^" + name + "$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"fake-id-12345"},
			"name_regex": "another-fake-name",
		}),
	}

	var existOnsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "1",
			"instances.#": "1",
			"instances.0.id":   CHECKSET,
			"instances.0.instance_name": name,
		}
	}

	var fakeOnsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var onsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsInstancesMapFunc,
		fakeMapFunc:  fakeOnsInstancesMapFunc,
	}
	onsInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceOnsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

%s

`, name, DataZoneCommonTestCase, OnsCommonTestCase)
}
