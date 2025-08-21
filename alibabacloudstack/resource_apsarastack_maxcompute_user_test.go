package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeUser(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_user.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 9999)
	name := fmt.Sprintf("tf_testAcck%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMaxcomputeUserDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"user_name":   "${var.name}",
					"description": "TestAccAlibabacloudStackMaxcomputeUser",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   name,
						"description": "TestAccAlibabacloudStackMaxcomputeUser",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_name":   "${var.name}_update",
					"description": "TestAccAlibabacloudStackMaxcomputeUser_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":   fmt.Sprintf("%s_update", name),
						"description": "TestAccAlibabacloudStackMaxcomputeUser_update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceMaxcomputeUserDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}
