package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeCu(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_cu.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMaxcomputeCu, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu_name":      name,
						"cu_num":       "1",
						"cluster_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMaxcomputeCu = `
data "alibabacloudstack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "alibabacloudstack_maxcompute_cu" "default"{
  cu_name      = "%s"
  cu_num       = "1"
  cluster_name = data.alibabacloudstack_maxcompute_clusters.default.clusters.0.cluster
}
`
