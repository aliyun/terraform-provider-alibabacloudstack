package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2McpserversDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_mcpservers.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfaccmcp%d", rand),
		dataSourceApiGatewayV2McpserversConfigDependence)

	// Define the basic test attributes
	testAcc := &dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                     "1",
				"mcp_servers.#":             "1",
				"mcp_servers.0.name":        fmt.Sprintf("tfaccmcp%d", rand),
				"mcp_servers.0.type":        "OPEN_API",
				"mcp_servers.0.description": fmt.Sprintf("tfaccmcp%d", rand),
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":         "0",
				"mcp_servers.#": "0",
			}
		},
		ExternalProviders: testAccExternalProviders,
	}

	// Run tests using customized test method
	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"gw_instance_id": "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}",
				"name_regex":     "${alibabacloudstack_api_gateway_v2_mcpserver.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"gw_instance_id": "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}",
				"name_regex":     "^fake-name.*$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"gw_instance_id": "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}",
				"ids":            []string{"${alibabacloudstack_api_gateway_v2_mcpserver.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"gw_instance_id": "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}",
				"ids":            []string{"fake-id"},
			}),
		},
	)
}

func dataSourceApiGatewayV2McpserversConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  	domain = "${var.name}.com"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	protocol = "HTTP"
}

resource "alibabacloudstack_api_gateway_v2_service" "default" {
	name = "${var.name}"
	service_source_type = "ip"
	protocol = "HTTP"
	service_nodes {
		ip = "192.168.1.1"
		port = 80
	}
	gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}
resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  name             = "${var.name}"
  description      = "${var.name}"
  type             = "OPEN_API"
  service          = "${alibabacloudstack_api_gateway_v2_service.default.service_id}"
  domains          = ["${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"]
  consumer_auth    = true
  gw_instance_id   = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}
`, name, ApiGatwayV2K8sInstanceTestCase("HIGRESS", "k8s"))
}
