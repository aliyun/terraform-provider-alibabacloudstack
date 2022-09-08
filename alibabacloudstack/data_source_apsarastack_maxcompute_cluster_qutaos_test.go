package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeClusterQutaoDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceAlibabacloudstackMaxcomputeClusterQutaos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_maxcompute_clusters.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_clusters.default", "clusters.cluster"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_clusters.default", "clusters.core_arch"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_clusters.default", "clusters.project"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_clusters.default", "clusters.region"),
				),
			},
		},
	})
}

const datasourceAlibabacloudstackMaxcomputeClusterQutaos = `
data "alibabacloudstack_maxcompute_clusters" "default"{
}

data "alibabacloudstack_maxcompute_cluster_qutaos" "default"{
    cluster = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}
`
