package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2ServiceSourcesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_service_sources.default"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-serviceSource*"`,
		}),
		fakeConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-fakeserviceSource*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_api_gateway_v2_service_source.default.id}"]`,
		}),
		fakeConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-serviceSource*"`,
			"ids":        `["${alibabacloudstack_api_gateway_v2_service_source.default.id}"]`,
		}),
		fakeConfig: APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-serviceSource*"`,
			"ids":        `["fake-id"]`,
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

func APIGatewayV2ServiceSourcesDataSourceCommonTestCaseNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		if k == "ids" {
			pairs = append(pairs, fmt.Sprintf(`%s = %s`, k, v))
		} else {
			pairs = append(pairs, fmt.Sprintf(`%s = %s`, k, v))
		}
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-serviceSource%d"
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
	nacos_registry = "127.0.0.1:8000"
}

data "alibabacloudstack_api_gateway_v2_service_sources" "default" {
  instance_id = alibabacloudstack_api_gateway_v2_service_source.default.instance_id
  
%s
}
`, rand, strings.Join(pairs, "\n"))

	return config
}
