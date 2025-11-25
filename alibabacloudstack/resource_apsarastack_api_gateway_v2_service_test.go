package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2Service_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_service.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Service")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ServiceCommonTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id":    "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"name":              "${var.name}",
					"description":       "${var.name}",
					"upstream_type":     "1",
					"load_balance_type": "1",
					"protocol":          "HTTP",
					"service_nodes": []map[string]interface{}{
						{
							"ip":     "127.0.0.1",
							"port":   "80",
							"weight": "100",
							"enable": "true",
						},
					},
					"health_check_struct": []map[string]interface{}{
						{
							"type":               "1",
							"health_path":        "/check",
							"http_statuses":      "200",
							"timeout":            "20000",
							"health_interval":    "30",
							"un_health_interval": "30",
							"http_successes":     "1",
							"http_failures":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                              name,
						"description":                       name,
						"upstream_type":                     "1",
						"load_balance_type":                 "1",
						"protocol":                          "HTTP",
						"service_nodes.#":                   "1",
						"service_nodes.0.ip":                "127.0.0.1",
						"service_nodes.0.weight":            "100",
						"health_check_struct.#":             "1",
						"health_check_struct.0.type":        "1",
						"health_check_struct.0.health_path": "/check",
						"health_check_struct.0.timeout":     "20000",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.name}update",
					"load_balance_type": "2",
					"upstream_type":     "1",
					"protocol":          "HTTPS",
					"service_nodes": []map[string]interface{}{
						{
							"ip":     "192.168.1.1",
							"port":   "8080",
							"weight": "1",
							"enable": "true",
						},
					},
					"health_check_struct": []map[string]interface{}{
						{
							"type":               "2",
							"timeout":            "20000",
							"health_interval":    "30",
							"un_health_interval": "30",
							"http_successes":     "1",
							"http_failures":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                              fmt.Sprintf("%supdate", name),
						"load_balance_type":                        "2",
						"upstream_type":                            "1",
						"protocol":                                 "HTTPS",
						"service_nodes.0.ip":                       "192.168.1.1",
						"service_nodes.0.port":                     "8080",
						"service_nodes.0.weight":                   "1",
						"health_check_struct.#":                    "1",
						"health_check_struct.0.type":               "2",
						"health_check_struct.0.health_path":        REMOVEKEY,
						"health_check_struct.0.timeout":            "20000",
						"health_check_struct.0.un_health_interval": "30",
					}),
				),
			},
		},
	})
}

func TestUatAlibabacloudStackApiGatewayV2Service_HSF(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_service.hsf"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Service")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(2000, 3000)
	name := fmt.Sprintf("testtf-hsf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ServiceForServiceSourceTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccApigwV2ServicePreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id":    "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"name":              "${var.name}",
					"description":       "${var.name}",
					"upstream_type":     "2",
					"load_balance_type": "4",
					"protocol":          "HSF",
					"real_service_name": "com.alibaba.edas.carshop.itemcenter.ItemService",
					"service_group":     "HSF",
					"service_version":   "1.0.0",
					"source_id":         "${alibabacloudstack_api_gateway_v2_service_source.default.source_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"description":       name,
						"upstream_type":     "2",
						"load_balance_type": "4",
						"protocol":          "HSF",
						"real_service_name": "com.alibaba.edas.carshop.itemcenter.ItemService",
						"service_group":     "HSF",
						"service_version":   "1.0.0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "${var.name}update",
					"real_service_name": "com.alibaba.edas.carshop.itemcenter.service.ItemService2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       fmt.Sprintf("%supdate", name),
						"real_service_name": "com.alibaba.edas.carshop.itemcenter.service.ItemService2",
					}),
				),
			},
		},
	})
}

func TestUatAlibabacloudStackApiGatewayV2Service_AIGwService(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_service.aigw_service"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Service")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(2000, 3000)
	name := fmt.Sprintf("testtf-hsf-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ServiceAIgwServiceTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id":      "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"name":                "${var.name}",
					"service_source_type": "ip",
					"protocol":            "HTTP",
					"service_nodes": []map[string]interface{}{
						{
							"ip":   "127.0.0.1",
							"port": "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                 name,
						"service_source_type":  "ip",
						"protocol":             "HTTP",
						"service_nodes.#":      "1",
						"service_nodes.0.ip":   "127.0.0.1",
						"service_nodes.0.port": "80",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTPS",
					"service_nodes": []map[string]interface{}{
						{
							"ip":   "172.16.0.1",
							"port": "443",
						},
						{
							"ip":   "192.168.1.1",
							"port": "9999",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":             "HTTPS",
						"service_nodes.#":      "2",
						"service_nodes.0.ip":   CHECKSET,
						"service_nodes.0.port": CHECKSET,
						"service_nodes.1.ip":   CHECKSET,
						"service_nodes.1.port": CHECKSET,
					}),
				),
			},
		},
	})
}

func ApiGatewayV2ServiceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, ApiGatwayV2K8sInstanceTestCase("SCG", "apig_k8s"))
}

func ApiGatewayV2ServiceAIgwServiceTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

`, name, ApiGatwayV2K8sInstanceTestCase("HIGRESS", "k8s"))
}

func ApiGatewayV2ServiceForServiceSourceTestCase(name string) string {
	edasAccessKey := os.Getenv("ALIBABACLOUDSTACK_EDAS_ACCESS_KEY")
	edasSecretKey := os.Getenv("ALIBABACLOUDSTACK_EDAS_SECRET_KEY")
	edasEndPointPort := os.Getenv("ALIBABACLOUDSTACK_EDAS_ENDPOINT_PORT")
	edasNameSpaceId := os.Getenv("ALIBABACLOUDSTACK_EDAS_NAMESPACE_ID")
	edasEndPoint := os.Getenv("ALIBABACLOUDSTACK_EDAS_ENDPOINT")
	return fmt.Sprintf(`

%s

resource "alibabacloudstack_api_gateway_v2_service_source" "default" {
  instance_id       = alibabacloudstack_api_gateway_v2_instance.default.id
  source_name       = var.name
  source_type       = "2"
  description       = var.name
  
  edas_end_point_port    = %s
  type                   = 1
  edas_name_space_id     = "%s"
  edas_access_key        = "%s"
  edas_secret_key        = "%s"
  edas_end_point         = "%s"
}
`, ApiGatewayV2ServiceCommonTestCase(name), edasEndPointPort, edasNameSpaceId, edasAccessKey, edasSecretKey, edasEndPoint)
}
