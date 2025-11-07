package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApiGatewayV2Signature_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_api_gateway_v2_signature.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApiGateWayV2Service{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeApiGatewayV2Signature")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: ApiGatewayV2SignatureBasicTestCase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sig_scheme_name": "test",
						"sig_alg":         "HmacSM3",
					}),
				),
			},
			{
				Config: ApiGatewayV2SignatureUpdateNameTestCase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sig_scheme_name": "test-update",
					}),
				),
			},
			{
				Config: ApiGatewayV2SignatureUpdateStatusTestCase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: ApiGatewayV2SignatureUpdateStatus2TestCase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
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

const ApiGatewayV2SignatureBasicTestCase = `
variable "signature_name" {
  default = "test"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
}
`

const ApiGatewayV2SignatureUpdateNameTestCase = `
variable "signature_name" {
  default = "test-update"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
}
`

const ApiGatewayV2SignatureUpdateStatusTestCase = `
variable "signature_name" {
  default = "test-update"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
  status          = "0"
}
`
const ApiGatewayV2SignatureUpdateStatus2TestCase = `
variable "signature_name" {
  default = "test-update"
}

variable "signature_algorithm" {
  default = "HmacSM3"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "testtf-apigw-instance"
  node_number        = "1"
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

data "alibabacloudstack_api_gateway_v2_instances" "default" {
  name_regex = alibabacloudstack_api_gateway_v2_instance.default.instance_name
}

resource "alibabacloudstack_api_gateway_v2_signature" "default" {
  sig_scheme_name = var.signature_name
  sig_alg         = var.signature_algorithm
  gw_instance_id  = alibabacloudstack_api_gateway_v2_instance.default.id
  status          = "1"
}
`
