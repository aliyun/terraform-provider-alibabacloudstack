package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksUserRoleBinding_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_user_role_binding.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksUserRoleBindingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksUserRoleBinding")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_userrole%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksUserRoleBindingBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id": "${alibabacloudstack_data_works_project.default.id}",
					"user_id":    "${alibabacloudstack_data_works_user.default.project_member_id}",
					"role_code":  "role_project_guest",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_code": "role_project_guest",
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

var AlibabacloudStackDataWorksUserRoleBindingMap0 = map[string]string{}

func AlibabacloudStackDataWorksUserRoleBindingBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alibabacloudstack_account" "current" {
}

data "alibabacloudstack_ascm_users" "default" {
 organization_id = "${data.alibabacloudstack_account.current.organization_id}"
}

resource "alibabacloudstack_data_works_project" "default" {
	name =           "${var.name}"
	description =    "${var.name}_desc"
	task_auth_type = "PROJECT"
}

resource "alibabacloudstack_data_works_user" "default" {
	project_id= "${alibabacloudstack_data_works_project.default.id}"
	user_id=    "${data.alibabacloudstack_ascm_users.default.users.0.primary_key}"
	role_code = ["role_project_admin"]
	lifecycle {
	    ignore_changes = [
	      role_code
	    ]
	}
}
`, name)
}
