package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_quotasDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Quota,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_quotas.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.quota_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.quota_type_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.region"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.cluster_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.used_vip_public"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quotas.default", "groups.allocate_vip_internal"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Quota = `

data "apsarastack_ascm_quotas" "default" {
  quota_type = "organization"
  quota_type_id = 1
  product_name = "SLB"
}
`
