package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAccountDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackAccountDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_account.current"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_account.current", "id"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackAccountDataSourceBasic = `
data "alibabacloudstack_account" "current" {
}
`
