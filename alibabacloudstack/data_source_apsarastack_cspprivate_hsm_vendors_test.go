package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCspprivateVendorsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cspprivate_hsm_vendors.default"
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackCspprivateVendorsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.#"),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.0.code"),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.0.name"),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.0.products.#"),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.0.products.0.code"),
					resource.TestCheckResourceAttrSet(resourceId, "vendors.0.products.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackCspprivateVendorsDataSource = `
data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}
`
