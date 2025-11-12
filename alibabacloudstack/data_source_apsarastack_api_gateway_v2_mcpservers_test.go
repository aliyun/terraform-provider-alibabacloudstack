package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackAPIGatewayV2McpServersDataSource_basic(t *testing.T) {
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
%s

data "alibabacloudstack_api_gateway_v2_mcpservers" "default" {
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_mcpserver.default.gw_instance_id}"
  %s
}
`, resourceApiGatewayV2McpserverDependenceNew(rand, nil), strings.Join(pairs, "\n  "))
}

// Dependency template generation method as required
func resourceApiGatewayV2McpserverDependenceNew(rand int, _ map[string]string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

%s

resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
	cs_cluster_id =   "${local.k8s_cluster_id}"
	k8s_cluster_name = "${var.name}"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}-apigw"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "HIGRESS"
  deploy_mode        = "k8s"
  deploy_cluster_code = "${alibabacloudstack_api_gateway_v2_k8s_cluster.default.id}"
  deploy_cluster_namespace = "${var.name}-namespace"
  ingress_class_name = "${var.name}-class"
  sls_enabled = "true"
  prometheus_enabled = "true"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  	domain = "${var.name}.com"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	protocol = "HTTP"
}

resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  name             = "${var.name}"
  description      = "${var.name}"
  type             = "OPEN_API"
  service          = "kubernetes.default.svc.cluster.local"
  domains          = ["testtf.com"]
  consumer_auth    = true
  gw_instance_id   = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}
`, rand, AckK8sCommonTestCase())
}
