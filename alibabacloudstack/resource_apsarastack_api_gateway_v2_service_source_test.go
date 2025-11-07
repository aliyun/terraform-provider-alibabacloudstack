package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2ServiceSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_service_source.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2ServiceSource")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ServiceSourceCommonTestCase)

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
					"source_name":      "${var.name}",
					"source_type":      "1",
					"instance_id":      "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"description":      "${var.name}",
					"check_type":       "1",
					"nacos_access_key": "root",
					"nacos_secret_key": "12345",
					"nacos_registry":   "127.0.0.1:8000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_name":      name,
						"source_type":      "1",
						"description":      name,
						"check_type":       "1",
						"nacos_access_key": "root",
						"nacos_secret_key": "12345",
						"nacos_registry":   "127.0.0.1:8000",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"nacos_secret_key", "edas_secret_key", "password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_name":      "${var.name}_update",
					"description":      "${var.name}_update",
					"check_type":       "2",
					"nacos_access_key": "admin",
					"nacos_secret_key": "newsecretkey",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_name":      fmt.Sprintf("%s_update", name),
						"description":      fmt.Sprintf("%s_update", name),
						"check_type":       "2",
						"nacos_access_key": "admin",
						"nacos_secret_key": "newsecretkey",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_type":          "2",
					"check_type":           "1",
					"edas_end_point_port":  "8080",
					"type":                 "1",
					"max_connection":       "10",
					"max_idle_connection":  "5",
					"connection_idle_time": "60",
					"database_type":        "0",
					"edas_name_space_id":   "test1234",
					"edas_access_key":      "root",
					"edas_secret_key":      "1234",
					"edas_end_point":       "127.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":          "2",
						"check_type":           "1",
						"edas_end_point_port":  "8080",
						"type":                 "1",
						"max_connection":       "10",
						"max_idle_connection":  "5",
						"connection_idle_time": "60",
						"database_type":        "0",
						"edas_name_space_id":   "test1234",
						"edas_access_key":      "root",
						"edas_end_point":       "127.0.0.1",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackApiGatewayV2ServiceSource_Eureka(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_service_source.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2ServiceSource")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("testtf-apigw-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2ServiceSourceCommonTestCase)

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
					"source_name":     "${var.name}",
					"source_type":     "3",
					"instance_id":     "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"description":     "${var.name}",
					"eureka_registry": "https://127.0.0.1:8000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_name":     name,
						"source_type":     "3",
						"description":     name,
						"eureka_registry": "https://127.0.0.1:8000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eureka_registry": "https://127.0.0.1:8888",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eureka_registry": "https://127.0.0.1:8888",
					}),
				),
			},
		},
	})
}

func ApiGatewayV2ServiceSourceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  	instance_name = "${var.name}"
	node_number = "1"
	instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
	broker_engine_type = "SCG"
	deploy_mode = "custom"
}

`, name)
}
