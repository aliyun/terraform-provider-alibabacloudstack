package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmMaxcomputeCuDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccApsaraStack%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(datasourceApsarastackMaxcomputeCus, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_maxcompute_cus.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_cus.default", "cus.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_cus.default", "cus.cu_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_cus.default", "cus.cu_num"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_cus.default", "cus.cluster_name"),
				),
			},
		},
	})
}

const datasourceApsarastackMaxcomputeCus = `
data "apsarastack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "apsarastack_maxcompute_cu" "default"{
  cu_name      = "%s"
  cu_num       = "1"
  cluster_name = data.apsarastack_maxcompute_clusters.default.clusters.0.cluster
}

data "apsarastack_maxcompute_cus" "default"{
	name_regex = "tf_testAccApsaraStack"
}
`
