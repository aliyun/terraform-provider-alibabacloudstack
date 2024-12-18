package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscm_User_GroupsDataSource(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlibabacloudStackAscm_User_Group_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ascm_user_groups.default"),
				),
			},
		},
	})
}

const dataSourceAlibabacloudStackAscm_User_Group_Organization = `
data "alibabacloudstack_ascm_user_groups" "default" {
	name_regex = "cxt"
	
}
output "groups" {
	  value = "${data.alibabacloudstack_ascm_user_groups.default.groups}"
	}
`
