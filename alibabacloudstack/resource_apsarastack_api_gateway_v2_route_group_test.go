package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2RouteGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_route_group.default"
	ra := resourceAttrInit(resourceId, ApiGatewayV2RouteGroupBasicMap)
	serviceFunc := func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeRouteGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testacc-routegroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApiGatewayV2RouteGroupCommonTestCase)

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
					"instance_id": "${alibabacloudstack_api_gateway_v2_instance.default.id}",
					"name":        "${var.name}",
					"base_path":   "/test",
					"description": "${var.name}",
					"domain_ids": []string{
						"${alibabacloudstack_api_gateway_v2_domain.domain0.domain_id}",
						"${alibabacloudstack_api_gateway_v2_domain.domain1.domain_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"base_path":   "/test",
						"description": name,
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
					"name":        "${var.name}-update",
					"base_path":   "/test1",
					"description": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("%s-update", name),
						"base_path":   "/test1",
						"description": fmt.Sprintf("%s-update", name),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_ids": []string{
						"${alibabacloudstack_api_gateway_v2_domain.domain0.domain_id}",
						"${alibabacloudstack_api_gateway_v2_domain.domain1.domain_id}",
						"${alibabacloudstack_api_gateway_v2_domain.domain2.domain_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_ids.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_ids": []string{
						"${alibabacloudstack_api_gateway_v2_domain.domain1.domain_id}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_ids.#": "1",
					}),
				),
			},
		},
	})
}

var ApiGatewayV2RouteGroupBasicMap = map[string]string{
	"instance_id": CHECKSET,
	"group_id":    CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
	"editable":    "true",
	"domains.#":   CHECKSET,
}

func ApiGatewayV2RouteGroupCommonTestCase(name string) string {
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

resource "alibabacloudstack_api_gateway_v2_domain" "domain1" {
  domain         = "${var.name}2.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTP"
  client_auth    = "0"
}

resource "alibabacloudstack_api_gateway_v2_domain" "domain2" {
  domain         = "${var.name}3.com"
  instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol       = "HTTP"
  client_auth    = "0"
}

`, name)
}
