package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHsmsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_hsms.default"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsms.default", "hsms.#"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_hsms.default", "hsms.0.id"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackHsmsDataSource = `
data "alibabacloudstack_hsm_vendors" "default" {
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_hsms" "default" {
  vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].code}"
  product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].products[0].code}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

`
