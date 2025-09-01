package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeUsersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(datasourceAlibabacloudstackMaxcomputeUsers, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_maxcompute_users.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.user_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.user_pk"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.user_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.user_type"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.organization_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.organization_name"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_maxcompute_users.default", "users.description"),
				),
			},
		},
	})
}

const datasourceAlibabacloudstackMaxcomputeUsers = `
resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
}

data "alibabacloudstack_maxcompute_users" "default"{
	name_regex = "tf_testAcck"
}
`
