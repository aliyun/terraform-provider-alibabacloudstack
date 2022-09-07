package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_RegionsByProduct_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_RegionsByProduct,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_regions_by_product.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_regions_by_product.default", "roles.region_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_regions_by_product.default", "roles.region_type"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_RegionsByProduct = `

data "alibabacloudstack_ascm_regions_by_product" "default" {
  output_file = "product_regions1"
  product_name = "ecs"
}
`
