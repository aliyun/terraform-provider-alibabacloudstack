package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackSlbsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackSlbsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_slbs.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_slbs.default", "slbs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_slbs.default", "ids.#"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckAlibabacloudStackSlbsDataSource = `
variable "name" {
	default = "tf-SlbDataSourceSlbs"
}
` + SlbCommonTestCase + `

data "alibabacloudstack_slbs" "default" {
 ids = ["${alibabacloudstack_slb.default.id}"]
}
`
