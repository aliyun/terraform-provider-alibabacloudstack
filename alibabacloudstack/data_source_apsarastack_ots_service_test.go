package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOtsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackOtsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ots_service.current"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_ots_service.current", "id"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_ots_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackOtsServiceDataSource = `
data "alibabacloudstack_ots_service" "current" {
	enable = "On"
}
`
