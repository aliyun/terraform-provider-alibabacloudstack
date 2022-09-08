package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_EcsInstanceFamilies_DataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Ecs_instance_families,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_ecs_instance_families.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ecs_instance_families.default", "families.instance_type_family_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ecs_instance_families.default", "families.generation"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_Ecs_instance_families = `

data "alibabacloudstack_ascm_ecs_instance_families" "default" {
  status = "Available"
}
`
