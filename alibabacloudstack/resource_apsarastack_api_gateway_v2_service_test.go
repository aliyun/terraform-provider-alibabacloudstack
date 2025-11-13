package alibabacloudstack

import (
	"fmt"
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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gw_instance_id":    "i-bp7e6ztkm1f5jarsp8eo",
					"name":              "${var.name}",
					"description":       "${var.name}",
					"upstream_type":     "1",
					"load_balance_type": "1",
					"protocol":          "HTTP",
					"service_type":      "0",
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
						"service_type":                      "0",
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

func ApiGatewayV2ServiceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
