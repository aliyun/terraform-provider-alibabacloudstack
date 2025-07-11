package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAcmConfigurationsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	resourceId := "data.alibabacloudstack_acm_configurations.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("tf_testacmconfig_%d", rand), dataSourceAcmConfigurationsConfigDependence)

	namespaceConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.id}",
			"data_id":      "${var.name}",
			"group_id":     "DEFAULT_GROUP",
		}),
	}

	var existAcmConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            CHECKSET,
			"configurations.#": CHECKSET,
		}
	}

	var fakeAcmConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"configurations.#": "0",
		}
	}

	var AcmConfigurationsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAcmConfigurationsMapFunc,
		fakeMapFunc:  fakeAcmConfigurationsMapFunc,
	}

	AcmConfigurationsCheckInfo.dataSourceTestCheck(t, rand, namespaceConfig)
}

func dataSourceAcmConfigurationsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "logical_id" {
  default = "%s:%s"
}

resource "alibabacloudstack_edas_namespace" "default" {
  	description = "${var.name}"
	namespace_name = "${var.name}"
	namespace_logical_id = "${var.logical_id}"
}

resource "apsarastack_acm_configuration" "default" {
	app_name = "${var.name}"
	content = "test"
	data_id = "${var.name}"
	group = "DEFAULT_GROUP"
	type = "text"
	namespace_id = "${alibabacloudstack_edas_namespace.default.id}"
}
`, name, defaultRegionToTest, name)
}
