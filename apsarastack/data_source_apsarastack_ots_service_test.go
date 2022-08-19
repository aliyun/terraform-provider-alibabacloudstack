package apsarastack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackOtsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackOtsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ots_service.current"),
					resource.TestCheckResourceAttrSet("data.apsarastack_ots_service.current", "id"),
					resource.TestCheckResourceAttr("data.apsarastack_ots_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackOtsServiceDataSource = `
data "apsarastack_ots_service" "current" {
	enable = "On"
}
`
