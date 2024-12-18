package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func SkipTestAccAlibabacloudStackApigatewayServiceDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackApigatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_api_gateway_service.current"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_api_gateway_service.current", "id"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_api_gateway_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackApigatewayServiceDataSource = `
data "alibabacloudstack_api_gateway_service" "current" {
	enable = "On"
}
`
