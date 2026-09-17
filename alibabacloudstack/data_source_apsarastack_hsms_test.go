package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hsms.default"
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		Providers: testYunDunProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackHsmsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "hsms.#"),
					resource.TestCheckResourceAttrSet(resourceId, "hsms.0.id"),
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
