package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2RouteGroupsDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_api_gateway_v2_route_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-routegroup%d", rand),
		dataSourceApiGatewayV2RouteGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"name_regex":  "tf-testacc-routegroup.*",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"name_regex":  "tf-testacc-fakeroutegroup.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_route_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"ids":         []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"name_regex":  "tf-testacc-routegroup.*",
			"ids":         []string{"${alibabacloudstack_api_gateway_v2_route_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alibabacloudstack_api_gateway_v2_route_group.default.instance_id}",
			"name_regex":  "tf-testacc-routegroup.*",
			"ids":         []string{"fake-id"},
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

func dataSourceApiGatewayV2RouteGroupsConfigDependence(name string) string {
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
`, name)
}
