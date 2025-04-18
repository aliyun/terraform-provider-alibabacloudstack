package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeClusterDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceAlibabacloudstackMaxcomputeClusters,
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

const datasourceAlibabacloudstackMaxcomputeClusters = `
data "alibabacloudstack_maxcompute_clusters" "default"{
}
`
