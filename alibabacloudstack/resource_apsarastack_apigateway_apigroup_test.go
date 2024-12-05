package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackApigatewayApigroup0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_apigateway_apigroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccApigatewayApigroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudApiService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCloudapiDescribeapigroupRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testaccapi_gatewayapi_group%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccApigatewayApigroupBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"api_group_name": "testGroupNameByAMP",

					"description": "test",

					"custom_trace_config": "{\"parameterLocation\":\"QUERY\",\"parameterName\":\"traceId\"}",

					"compatible_flags": "supportSSE",

					"user_log_config": "{\"requestBody\":false,\"responseBody\":false,\"queryString\":\"\",\"requestHeaders\":\"\",\"responseHeaders\":\"\",\"jwtClaims\":\"\"}",

					"passthrough_headers": "eagleeye-rpcid,x-b3-traceid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"api_group_name": "testGroupNameByAMP",

						"description": "test",

						"custom_trace_config": "{\"parameterLocation\":\"QUERY\",\"parameterName\":\"traceId\"}",

						"compatible_flags": "supportSSE",

						"user_log_config": "{\"requestBody\":false,\"responseBody\":false,\"queryString\":\"\",\"requestHeaders\":\"\",\"responseHeaders\":\"\",\"jwtClaims\":\"\"}",

						"passthrough_headers": "eagleeye-rpcid,x-b3-traceid",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"api_group_name": "modifyApiGroupName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"api_group_name": "modifyApiGroupName",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test modify",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"compatible_flags": "disableInnerDomain",

					"custom_trace_config": "{\"parameterLocation\":\"HEADER\",\"parameterName\":\"traceId\"}",

					"user_log_config": "{\"requestBody\":true,\"responseBody\":true,\"queryString\":\"\",\"requestHeaders\":\"\",\"responseHeaders\":\"\",\"jwtClaims\":\"\"}",

					"passthrough_headers": "eagleeye-rpcid,x-b3-traceid,host",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"compatible_flags": "disableInnerDomain",

						"custom_trace_config": "{\"parameterLocation\":\"HEADER\",\"parameterName\":\"traceId\"}",

						"user_log_config": "{\"requestBody\":true,\"responseBody\":true,\"queryString\":\"\",\"requestHeaders\":\"\",\"responseHeaders\":\"\",\"jwtClaims\":\"\"}",

						"passthrough_headers": "eagleeye-rpcid,x-b3-traceid,host",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccApigatewayApigroupCheckmap = map[string]string{

	"api_group_name": CHECKSET,

	"description": CHECKSET,

	"https_policy": CHECKSET,

	"custom_trace_config": CHECKSET,

	"user_log_config": CHECKSET,

	"passthrough_headers": CHECKSET,

	"vpc_domain": CHECKSET,

	"sub_domain": CHECKSET,

	"modified_time": CHECKSET,

	"illegal_status": CHECKSET,

	"instance_type": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"instance_id": CHECKSET,

	"compatible_flags": CHECKSET,

	"create_time": CHECKSET,

	"billing_status": CHECKSET,

	"traffic_limit": CHECKSET,

	"group_id": CHECKSET,

	"ipv6_status": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccApigatewayApigroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
