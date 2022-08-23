package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackDataWorksUserRoleBinding_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_data_works_user_role_binding.default"
	ra := resourceAttrInit(resourceId, ApsaraStackDataWorksUserRoleBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDataWorksUserRoleBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworksuser%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackDataWorksUserRoleBindingBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id": "10023",
					"user_id":    "5247087457099176824",
					"role_code":  "role_project_guest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id": "10023",
						"user_id":    "5247087457099176824",
						"role_code":  "role_project_guest",
					}),
				),
			},
		},
	})
}

var ApsaraStackDataWorksUserRoleBindingMap0 = map[string]string{}

func ApsaraStackDataWorksUserRoleBindingBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
