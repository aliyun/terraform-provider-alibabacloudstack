package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackMaxcomputeUser(t *testing.T) {
	resourceId := "apsarastack_maxcompute_user.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccApsaraStack%d", rand)

	resource.Test(t, resource.TestCase{
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
						"description": "TestAccApsaraStackMaxcomputeUser",
					}),
				),
			},
			{
				Config: fmt.Sprintf(testAccMaxcomputeUserUpdate, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   name,
						"description": "TestAccApsaraStackMaxcomputeUserUpdate",
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeUser = `
resource "apsarastack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccApsaraStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
`

const testAccMaxcomputeUserUpdate = `
resource "apsarastack_maxcompute_user" "default"{
  user_name             = "%s"
  description           = "TestAccApsaraStackMaxcomputeUserUpdate"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
`
