package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMaxcomputeProject_basic(t *testing.T) {
	resourceId := "alibabacloudstack_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	// rand := getAccTestRandInt(1000, 9999)
	// name := fmt.Sprintf("tf_testAcck%d", rand)
	name := "tf_testAcck3043"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMaxcomputeProjectDependence)

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
					"name":           "${var.name}",
					"disk":           "50",
					"account":        "ascm-dw-1755770304175",
					"account_pk":     "1784955770304192",
					"quota_id":       "29",
					"external_table": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":           name,
						"disk":           "50",
						"account":        "ascm-dw-1755770304175",
						"account_pk":     "1784955770304192",
						"quota_id":       "29",
						"external_table": "true",
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

func resourceMaxcomputeProjectDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

`, name)
}
