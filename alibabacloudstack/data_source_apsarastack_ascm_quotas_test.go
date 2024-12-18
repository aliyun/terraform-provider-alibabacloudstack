package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_quotasDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Quota,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_quotas.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.quota_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.quota_type_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.region"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.cluster_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.used_vip_public"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_quotas.default", "groups.allocate_vip_internal"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_Quota = `

data "alibabacloudstack_ascm_quotas" "default" {
  quota_type = "organization"
  quota_type_id = 1
  product_name = "SLB"
}
`
