package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmVendorsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hsm_vendors.default"
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHsmVendorsDataSource,
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

const testAccCheckAlibabacloudStackHsmVendorsDataSource = `
data "alibabacloudstack_hsm_vendors" "default" {
}
`
