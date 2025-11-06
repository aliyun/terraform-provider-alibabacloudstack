package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2RouteGroupsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_route_groups.default"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-routegroup.*"`,
		}),
		fakeConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-fakeroutegroup.*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_api_gateway_v2_route_group.default.id}"]`,
		}),
		fakeConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-routegroup.*"`,
			"ids":        `["${alibabacloudstack_api_gateway_v2_route_group.default.id}"]`,
		}),
		fakeConfig: ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"tf-testacc-routegroup.*"`,
			"ids":        `["fake-id"]`,
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"route_groups.#":                     CHECKSET,
			"route_groups.0.id":                  CHECKSET,
			"route_groups.0.instance_id":         CHECKSET,
			"route_groups.0.group_id":            CHECKSET,
			"route_groups.0.name":                fmt.Sprintf("tf-testacc-routegroup%d", rand),
			"route_groups.0.base_path":           "/test",
			"route_groups.0.description":         "test description",
			"route_groups.0.create_time":         CHECKSET,
			"route_groups.0.update_time":         CHECKSET,
			"route_groups.0.editable":            "true",
			"route_groups.0.domains.#":           "1",
			"route_groups.0.domains.0.protocol":  "HTTP",
			"route_groups.0.domains.0.domain":    fmt.Sprintf("tf-testacc-routegroup%d1.com", rand),
			"route_groups.0.domains.0.domain_id": CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"route_groups.#": "0",
		}
	}

	dsa := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	dsa.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func ApiGatewayV2RouteGroupsDataSourceCommonTestCaseNew(rand int, attrMap map[string]string) string {
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
  default = "tf-testacc-routegroup%d"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = var.name
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_domain" "domain0" {
  domain         = "${var.name}1.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTP"
  client_auth    = "0"
}

resource "alibabacloudstack_api_gateway_v2_route_group" "default" {
  name           = var.name
  base_path      = "/test"
  description    = "test description"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  domain_ids     = [alibabacloudstack_api_gateway_v2_domain.domain0.domain_id]
}

data "alibabacloudstack_api_gateway_v2_route_groups" "default" {
  instance_id = alibabacloudstack_api_gateway_v2_route_group.default.instance_id
  
%s
}
`, rand, strings.Join(pairs, "\n"))

	return config
}
