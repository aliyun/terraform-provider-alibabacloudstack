package alibabacloudstack

import (
	"fmt"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackMaxcomputeUser(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_user.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAccAlibabacloudStack%d", rand)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMaxcomputeUser, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   name,
						"description": "TestAccAlibabacloudStackMaxcomputeUser",
					}),
				),
			},
			{
				Config: fmt.Sprintf(testAccMaxcomputeUserUpdate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   name,
						"description": "TestAccAlibabacloudStackMaxcomputeUserUpdate",
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeUser = `
resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
`

const testAccMaxcomputeUserUpdate = `
resource "alibabacloudstack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccAlibabacloudStackMaxcomputeUserUpdate"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
`
