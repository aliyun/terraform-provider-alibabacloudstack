package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_RegionsByProduct_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_RegionsByProduct,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_regions_by_product.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_regions_by_product.default", "roles.region_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_regions_by_product.default", "roles.region_type"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_RegionsByProduct = `

data "apsarastack_ascm_regions_by_product" "default" {
  output_file = "product_regions1"
  product_name = "ecs"
}
`
