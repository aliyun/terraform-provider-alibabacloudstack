package alibabacloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscm_RegionsByProduct_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
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
  product_name = "ecs"
}
`
