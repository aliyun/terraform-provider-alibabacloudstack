package apsarastack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func SkipTestAccApsaraStackApigatewayServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackApigatewayServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_api_gateway_service.current"),
					resource.TestCheckResourceAttrSet("data.apsarastack_api_gateway_service.current", "id"),
					resource.TestCheckResourceAttr("data.apsarastack_api_gateway_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackApigatewayServiceDataSource = `
data "apsarastack_api_gateway_service" "current" {
	enable = "On"
}
`
