package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Resource_Groups_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_Resource_Group_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_resource_groups.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_resource_groups.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_resource_groups.default", "groups.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_resource_groups.default", "groups.organization_id"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_Resource_Group_Organization = `
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org"
  parent_id = "1"
}
 resource "alibabacloudstack_ascm_resource_group" "default" {
  organization_id = alibabacloudstack_ascm_organization.default.org_id
  name = "alibabacloudstack-Datasource-resourceGroup"
}
data "alibabacloudstack_ascm_resource_groups" "default" {
  name_regex = alibabacloudstack_ascm_resource_group.default.name
}
`
