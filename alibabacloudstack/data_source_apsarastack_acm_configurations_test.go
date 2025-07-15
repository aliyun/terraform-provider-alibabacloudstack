package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAcmConfigurationsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	resourceId := "data.alibabacloudstack_acm_configurations.default"

	name := fmt.Sprintf("tftestd%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAcmConfigurationsConfigDependence)

	namespaceConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}11",
		}),
	}
	
	appnameConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"app_name": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"app_name": "tftestacmconfigdataempty",
		}),
	}
	groupConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"group_id": "DEFAULT_GROUP",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"group_id": "tftestacmconfigdataempty",
		}),
	}
	dataIdConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"data_id": "name",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace_id": "${alibabacloudstack_edas_namespace.default.tenant_id}",
			"data_id": "tftestacmconfigdataempty",
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

	AcmConfigurationsCheckInfo.dataSourceTestCheck(t, rand, namespaceConfig, appnameConfig, groupConfig, dataIdConfig)
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

resource "alibabacloudstack_acm_configuration" "default" {
	app_name = "${var.name}"
	content = "test"
	data_id = "${var.name}"
	group = "DEFAULT_GROUP"
	type = "text"
	desc = "${var.name}"
	namespace_id = "${alibabacloudstack_edas_namespace.default.tenant_id}"
}
`, name, defaultRegionToTest, name)
}
