package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_EcsInstanceFamilies_DataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Ecs_instance_families,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_ecs_instance_families.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ecs_instance_families.default", "families.instance_type_family_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ecs_instance_families.default", "families.generation"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Ecs_instance_families = `

data "apsarastack_ascm_ecs_instance_families" "default" {
  status = "Available"
}
`
