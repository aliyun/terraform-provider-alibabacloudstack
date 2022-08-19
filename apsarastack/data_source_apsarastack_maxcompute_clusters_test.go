package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmMaxcomputeClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: datasourceApsarastackMaxcomputeClusters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_maxcompute_clusters.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_clusters.default", "clusters.cluster"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_clusters.default", "clusters.core_arch"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_clusters.default", "clusters.project"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_clusters.default", "clusters.region"),
				),
			},
		},
	})
}

const datasourceApsarastackMaxcomputeClusters = `
data "apsarastack_maxcompute_clusters" "default"{
}
`
