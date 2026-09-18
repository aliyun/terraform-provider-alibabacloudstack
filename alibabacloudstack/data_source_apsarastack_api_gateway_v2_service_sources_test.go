package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2ServiceSourcesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_service_sources.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-serviceSource%d", rand),
		dataSourceApiGatewayV2ServiceSourcesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"name_regex":  "tf-testacc-serviceSource*",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"name_regex":  "tf-testacc-fakeserviceSource*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_service_source.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"ids":         []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"name_regex":  "tf-testacc-serviceSource*",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_service_source.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_service_source.default.instance_id}",
			"name_regex":  "tf-testacc-serviceSource*",
			"ids":         []string{"fake-id"},
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"sources.#":                  CHECKSET,
			"sources.0.instance_id":      CHECKSET,
			"sources.0.source_id":        CHECKSET,
			"sources.0.source_name":      fmt.Sprintf("tf-testacc-serviceSource%d", rand),
			"sources.0.source_type":      "1",
			"sources.0.description":      fmt.Sprintf("tf-testacc-serviceSource%d", rand),
			"sources.0.create_time":      CHECKSET,
			"sources.0.source_type_name": "NACOS",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"sources.#": "0",
		}
	}

	dsa := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	dsa.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceApiGatewayV2ServiceSourcesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_service_source" "default" {
	source_name = "${var.name}"
	source_type = "1"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	description = "${var.name}"
	check_type = "1"
	nacos_access_key = "root"
	nacos_secret_key = "12345"
	nacos_registry = "172.16.100.1:8000"
}
`, name)
}
