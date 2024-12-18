package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_Users_DataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_User_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_users.default"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_User_Organization = `
data "alibabacloudstack_ascm_users" "default" {
}
`
