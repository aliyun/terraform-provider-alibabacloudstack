package alibabacloudstack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackAscmMaxcomputeUserDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccAlibabacloudStack%d", rand)
	resource.Test(t, resource.TestCase{
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
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
data "alibabacloudstack_maxcompute_users" "default"{
	name_regex = "tf_testAccAlibabacloudStack"
}
`
