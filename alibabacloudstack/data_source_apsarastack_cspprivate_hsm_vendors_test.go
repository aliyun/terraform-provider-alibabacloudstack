package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCspprivateVendorsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackCspprivateVendorsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_cspprivate_hsm_vendors.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.0.name"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.0.products.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.0.products.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_cspprivate_hsm_vendors.default", "vendors.0.products.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackCspprivateVendorsDataSource = `
data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}
`
