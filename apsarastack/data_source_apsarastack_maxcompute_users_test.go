package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackAscmMaxcomputeUserDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccApsaraStack%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(datasourceApsarastackMaxcomputeUsers, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_maxcompute_users.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.user_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.user_pk"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.user_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.user_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.organization_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.organization_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_maxcompute_users.default", "users.description"),
				),
			},
		},
	})
}

const datasourceApsarastackMaxcomputeUsers = `
resource "apsarastack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccApsaraStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
data "apsarastack_maxcompute_users" "default"{
	name_regex = "tf_testAccApsaraStack"
}
`
