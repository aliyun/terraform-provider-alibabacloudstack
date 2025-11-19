package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func ApiGatewayV2RouteDependence(name string) string {
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

resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name = "${var.name}"
  description = "${var.name}"
  protocol = "HTTP"
  upstream_type = "1"
  load_balance_type = "1"
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
  service_nodes {
    ip = "127.0.0.1"
    port = "80"
    weight = "100"
    enable = "true"
  }
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain         = "${var.name}.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTP"
}

`, name)
}

func ApiGatewayV2SourceRouteDependence(name string) string {
	return fmt.Sprintf(`
%s

resource "alibabacloudstack_api_gateway_v2_instance" "cascade" {
  instance_name      = "${var.name}-cascade"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name      = "${var.name}"
  cascade_instance_id = "${alibabacloudstack_api_gateway_v2_instance.cascade.id}"
}

resource "alibabacloudstack_api_gateway_v2_cascade_link" "default" {
  cascade_instance_id = "${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"
  link_name = "${var.name}"
  source_instance_address =  "10.17.94.180"
  source_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}

`, ApiGatewayV2RouteDependence(name))
}

func TestAccAlibabacloudStackApiGatewayV2Route_basic(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_api_gateway_v2_route.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApigwV2Route")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testacc-route%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2RouteDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"route_name":     "${var.name}",
					"strip_prefix":   "2",
					"order":          "100",
					"route_path":     []string{"/testtc", "/test/aaa"},
					"methods":        []string{"GET", "POST", "PUT", "DELETE"},
					"header": []map[string]interface{}{
						{
							"key":   "header",
							"value": "aaaaa",
						},
					},

					"cookie": []map[string]interface{}{
						{
							"key":   "cookie",
							"value": "bbbbb",
						},
					},

					"query_param": []map[string]interface{}{
						{
							"key":   "query",
							"value": "ccccc",
						},
					},

					"domain_ids": []string{"${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"},

					"service_ids": []map[string]interface{}{
						{
							"service_id": "${alibabacloudstack_api_gateway_v2_service.default.service_id}",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":               name,
						"strip_prefix":             "2",
						"order":                    "100",
						"methods.#":                "4",
						"route_path.#":             "2",
						"header.#":                 "1",
						"header.0.key":             "header",
						"cookie.#":                 "1",
						"cookie.0.key":             "cookie",
						"query_param.#":            "1",
						"query_param.0.key":        "query",
						"domain_ids.#":             "1",
						"service_ids.#":            "1",
						"service_ids.0.service_id": CHECKSET,
						"service_ids.0.weight":     "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"strip_prefix": "3",
					"order":        "80",

					"route_path": []string{"/testtc/aaaa/*", "/test/aaa", "/test/bbb"},

					"methods": []string{"GET", "POST", "DELETE"},

					"header": []map[string]interface{}{
						{
							"key":   "header",
							"value": "aaaaa",
						},
						{
							"key":   "header2",
							"value": "aaaaa2",
						},
					},

					"domain_ids": []string{"${alibabacloudstack_api_gateway_v2_domain.default.domain_id}"},

					"service_ids": REMOVEKEY,
					"service_id":  "${alibabacloudstack_api_gateway_v2_service.default.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":               name,
						"strip_prefix":             "3",
						"order":                    "80",
						"methods.#":                "3",
						"route_path.#":             "3",
						"header.#":                 "2",
						"service_ids":              REMOVEKEY,
						"service_ids.#":            REMOVEKEY,
						"service_ids.0.service_id": REMOVEKEY,
						"service_ids.0.weight":     REMOVEKEY,
						"service_id":               CHECKSET,
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

func TestAccAlibabacloudStackApiGatewayV2Route_Cascade(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_route.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApigwV2Route")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testacc-route%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2SourceRouteDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id":   "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"route_name":       "${var.name}",
					"strip_prefix":     "2",
					"order":            "100",
					"route_path":       []string{"/testtc", "/test/aaa"},
					"methods":          []string{"GET", "POST", "PUT", "DELETE"},
					"cascade_link_ids": []string{"${alibabacloudstack_api_gateway_v2_cascade_link.default.id}"},
					"header": []map[string]interface{}{
						{
							"key":   "header",
							"value": "aaaaa",
						},
					},

					"cookie": []map[string]interface{}{
						{
							"key":   "cookie",
							"value": "bbbbb",
						},
					},

					"query_param": []map[string]interface{}{
						{
							"key":   "query",
							"value": "ccccc",
						},
					},

					"service_id": "${alibabacloudstack_api_gateway_v2_service.default.service_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":        name,
						"strip_prefix":      "2",
						"order":             "100",
						"methods.#":         "4",
						"route_path.#":      "2",
						"header.#":          "1",
						"header.0.key":      "header",
						"cookie.#":          "1",
						"cookie.0.key":      "cookie",
						"query_param.#":     "1",
						"query_param.0.key": "query",
						"service_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"strip_prefix": "3",
					"order":        "80",

					"route_path": []string{"/testtc/aaaa/*", "/test/aaa", "/test/bbb"},

					"methods": []string{"GET", "POST", "DELETE"},

					"header": []map[string]interface{}{
						{
							"key":   "header",
							"value": "aaaaa",
						},
						{
							"key":   "header2",
							"value": "aaaaa2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_name":   name,
						"strip_prefix": "3",
						"order":        "80",
						"methods.#":    "3",
						"route_path.#": "3",
						"header.#":     "2",
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
