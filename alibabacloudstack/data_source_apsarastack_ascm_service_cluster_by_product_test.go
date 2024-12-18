package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Service_ClusterByProductDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_ServiceClusterByProductbasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_service_cluster_by_product.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_service_cluster_by_product.default", "cluster_list"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_ServiceClusterByProductbasic = `

data "alibabacloudstack_ascm_service_cluster_by_product" "default" {
  product_name = "ecs"
}
`
