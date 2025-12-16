package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmVendorsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHsmVendorsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_hsm_vendors.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.name"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.0.code"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsm_vendors.default", "vendors.0.products.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackHsmVendorsDataSource = `
data "alibabacloudstack_hsm_vendors" "default" {
}
`
