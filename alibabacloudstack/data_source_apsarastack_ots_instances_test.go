package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOtsInstancesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ots_instances.default"
	name := fmt.Sprintf("tf-otsinst-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOtsInstancesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ots_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ots_instance.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ots_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ots_instance.default.name}_fake",
		}),
	}

	specificationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_ots_instance.default.id}"},
			"specification": "HYBRID",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_ots_instance.default.id}"},
			"specification": "SSD",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_ots_instance.default.id}"},
			"name_regex":    "${alibabacloudstack_ots_instance.default.name}",
			"specification": "HYBRID",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_ots_instance.default.id}_fake"},
			"name_regex":    "${alibabacloudstack_ots_instance.default.name}_fake",
			"specification": "SSD",
		}),
	}

	var existOtsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                   "1",
			"instances.#":               "1",
			"instances.0.id":            CHECKSET,
			"instances.0.name":          name,
			"instances.0.user_id":       CHECKSET,
			"instances.0.description":   name,
			"instances.0.specification": "HYBRID",
			"instances.0.vcu_quota":     CHECKSET,
			"instances.0.create_time":   CHECKSET,
			"instances.0.tags.#":        CHECKSET,
		}
	}

	var fakeOtsInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var otsInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsInstancesMapFunc,
		fakeMapFunc:  fakeOtsInstancesMapFunc,
	}
	otsInstancesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, specificationConf, allConf)
}

func dataSourceOtsInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_ots_clusters" "default" {}

resource "alibabacloudstack_ots_instance" "default" {
  name = var.name
  description   = var.name
  specification  = "HYBRID"
}
`, name)
}
