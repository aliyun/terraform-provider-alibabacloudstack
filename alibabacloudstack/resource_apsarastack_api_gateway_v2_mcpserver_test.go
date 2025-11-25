package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func buildBasicGwInstance(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name = "${var.name}"
  service_source_type = "ip"
  protocol = "HTTP"
  service_nodes {
	ip = "192.168.1.1"
	port = 443
  }
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  	domain = "${var.name}.com"
	instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
	protocol = "HTTP"
}
`, name, ApiGatwayV2K8sInstanceTestCase("HIGRESS", "k8s"))
}

func TestAccAlibabacloudStackApiGatewayV2Mcpserver_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_mcpserver.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"name":              "testtf",
		"description":       "testddd",
		"type":              "OPEN_API",
		"domains.#":         "1",
		"domains.0":         CHECKSET,
		"consumer_auth":     "true",
		"services.#":        "1",
		"services.0.port":   "443",
		"services.0.weight": "100",
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeMcpserver")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildBasicGwInstance)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "testddd",
					"type":        "OPEN_API",
					"service":     "${alibabacloudstack_api_gateway_v2_service.default.id}",
					"domains": []string{
						"${alibabacloudstack_api_gateway_v2_domain.default.id}",
					},
					"consumer_auth":  "true",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"description":       "testddd",
						"type":              "OPEN_API",
						"service":           CHECKSET,
						"domains.#":         "1",
						"domains.0":         CHECKSET,
						"consumer_auth":     "true",
						"services.#":        "1",
						"services.0.name":   CHECKSET,
						"services.0.port":   "443",
						"services.0.weight": "100",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Mcpserver_database(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_mcpserver.database"
	ra := resourceAttrInit(resourceId, map[string]string{
		"description":       CHECKSET,
		"type":              "DATABASE",
		"services.#":        CHECKSET,
		"services.0.name":   CHECKSET,
		"services.0.port":   CHECKSET,
		"services.0.weight": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeMcpserver")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-mcp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildBasicGwInstance)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "${var.name}",
					"type":        "DATABASE",
					"service":     "${alibabacloudstack_api_gateway_v2_service.default.id}",
					"domains": []string{
						"${alibabacloudstack_api_gateway_v2_domain.default.id}",
					},
					"consumer_auth": "true",
					"db_host":       "172.24.1.220",
					"db_port":       "81",
					"db_username":   "root",
					"db_password":   "test1234",
					"db_name":       "db",
					"other_params": map[string]string{
						"aaa": "bbb",
						"ccc": "ddd",
					},
					"db_type":        "MYSQL",
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"description":       name,
						"service":           CHECKSET,
						"domains.#":         "1",
						"domains.0":         CHECKSET,
						"consumer_auth":     "true",
						"services.#":        "1",
						"services.0.name":   CHECKSET,
						"services.0.port":   "443",
						"services.0.weight": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
					"db_username": "root22",
					"db_password": "test123444",
					"db_name":     "db2",
					"other_params": map[string]string{
						"aaa": "111",
						"ddd": "ddd",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      fmt.Sprintf("%s_update", name),
						"db_username":      "root22",
						"db_password":      "test123444",
						"db_name":          "db2",
						"other_params.%":   "2",
						"other_params.aaa": "111",
						"other_params.ddd": "ddd",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2Mcpserver_directRoute(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_mcpserver.direct_route"
	ra := resourceAttrInit(resourceId, map[string]string{
		"description":       CHECKSET,
		"type":              "DIRECT_ROUTE",
		"services.#":        CHECKSET,
		"services.0.name":   CHECKSET,
		"services.0.port":   CHECKSET,
		"services.0.weight": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeMcpserver")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-mcp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, buildBasicGwInstance)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"description": "${var.name}",
					"type":        "DIRECT_ROUTE",
					"service":     "${alibabacloudstack_api_gateway_v2_service.default.id}",
					"domains": []string{
						"${alibabacloudstack_api_gateway_v2_domain.default.id}",
					},
					"consumer_auth":     "false",
					"direct_route_path": "/sse",
					"direct_route_type": "sse",
					"gw_instance_id":    "${alibabacloudstack_api_gateway_v2_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"description":       name,
						"type":              "DIRECT_ROUTE",
						"service":           CHECKSET,
						"domains.#":         "1",
						"domains.0":         CHECKSET,
						"consumer_auth":     "false",
						"direct_route_path": "/sse",
						"direct_route_type": "sse",
						"services.#":        "1",
						"services.0.name":   CHECKSET,
						"services.0.port":   "443",
						"services.0.weight": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.name}_update",
					"direct_route_path": "/test",
					"direct_route_type": "streamable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       fmt.Sprintf("%s_update", name),
						"direct_route_path": "/test",
						"direct_route_type": "streamable",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
