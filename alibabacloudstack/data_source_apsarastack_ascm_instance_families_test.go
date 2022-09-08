package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Instance_families_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Instance_families,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_ecs_instance_families.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ecs_instance_families.default", "families.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ecs_instance_families.default", "families.status"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ecs_instance_families.default", "families.resource_type"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_Instance_families = `

data "alibabacloudstack_ascm_ecs_instance_families" "default" {
status = "Available"
}
`
