package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_RamServiceRoles_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_RamServiceRoles,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_ram_service_roles.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.role_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.product"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.organization_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ram_service_roles.default", "roles.aliyun_user_id"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_RamServiceRoles = `

data "apsarastack_ascm_ram_service_roles" "default" {
  product = "ecs"
}

`
