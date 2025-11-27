package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackApiGatewayV2McpserversDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_api_gateway_v2_mcpservers.default"

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
	}

	// Run tests using customized test method
	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfigNew(rand, map[string]string{
				"name_regex": `"${alibabacloudstack_api_gateway_v2_mcpserver.default.name}"`,
			}),
			fakeConfig: testAccConfigNew(rand, map[string]string{
				"name_regex": `"^fake-name.*$"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfigNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_api_gateway_v2_mcpserver.default.id}"]`,
			}),
			fakeConfig: testAccConfigNew(rand, map[string]string{
				"ids": `["fake-id"]`,
			}),
		},
	)
}

// Helper function to generate test configurations dynamically based on input attributes
func testAccConfigNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
	variable "name" {
	  default = "tfaccmcp%d"
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
data "alibabacloudstack_api_gateway_v2_mcpservers" "default" {
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}"
  %s
}
`, rand, ApiGatwayV2K8sInstanceTestCase("HIGRESS", "k8s"), strings.Join(pairs, "\n  "))
}
