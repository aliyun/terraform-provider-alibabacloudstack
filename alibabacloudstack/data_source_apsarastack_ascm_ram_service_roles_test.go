package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_RamServiceRoles_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_RamServiceRoles,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_ram_service_roles.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.description"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.role_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.product"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.organization_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_ram_service_roles.default", "roles.aliyun_user_id"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_RamServiceRoles = `

data "alibabacloudstack_ascm_ram_service_roles" "default" {
  product = "ecs"
}

`
