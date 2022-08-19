package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Service_ClusterByProductDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_ServiceClusterByProductbasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_service_cluster_by_product.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_service_cluster_by_product.default", "cluster_list"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_ServiceClusterByProductbasic = `

data "apsarastack_ascm_service_cluster_by_product" "default" {
  product_name = "ecs"
}
`
