package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_OrganizationDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_E_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_organizations.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.cuser_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.muser_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.alias"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.parent_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ascm_organizations.default", "organizations.internal"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_E_Organization = `

resource "alibabacloudstack_ascm_organization" "org" {
 name = "Tf-testing-DataSourceascm-organization"
 parent_id = "1"
}

data "alibabacloudstack_ascm_organizations" "default" {
  name_regex = alibabacloudstack_ascm_organization.org.name
}
`
