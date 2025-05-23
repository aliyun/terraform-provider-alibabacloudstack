package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApigatewayAppAttachment(t *testing.T) {
	var v *cloudapi.AuthorizedApp

	resourceId := "alibabacloudstack_api_gateway_app_attachment.default"
	ra := resourceAttrInit(resourceId, apigatewayAppAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccApp_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceApigatewayAppAttachmentConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"api_id":     "${alibabacloudstack_api_gateway_api.default.api_id}",
					"group_id":   "${alibabacloudstack_api_gateway_group.default.id}",
					"stage_name": "PRE",
					"app_id":     "${alibabacloudstack_api_gateway_app.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func resourceApigatewayAppAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
resource "alibabacloudstack_api_gateway_group" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiGroup Description"
}
resource "alibabacloudstack_api_gateway_api" "default" {
  name        = "${var.name}"
  group_id    = "${alibabacloudstack_api_gateway_group.default.id}"
  description = "description"
  auth_type   = "APP"

  request_config {
    protocol = "HTTP"
    method   = "GET"
    path     = "/test/path"
    mode     = "MAPPING"
  }

  service_type = "HTTP"

  http_service_config {
    address   = "http://apigateway-backend.alicloudapi.com:8080"
    method    = "GET"
    path      = "/web/cloudapi"
    timeout   = 22
    aone_name = "cloudapi-openapi"
  }

  request_parameters {
      name         = "aa"
      type         = "STRING"
      required     = "OPTIONAL"
      in           = "QUERY"
      in_service   = "QUERY"
      name_service = "testparams"
    }
}

resource "alibabacloudstack_api_gateway_app" "default" {
  name        = "${var.name}"
  description = "tf_testAccApiAPP Description"
}

 `, name)
}

var apigatewayAppAttachmentBasicMap = map[string]string{
	"api_id":     CHECKSET,
	"group_id":   CHECKSET,
	"stage_name": "PRE",
	"app_id":     CHECKSET,
}
