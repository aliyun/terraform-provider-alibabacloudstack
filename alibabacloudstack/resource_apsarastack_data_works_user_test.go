package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDataWorksUser_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_data_works_user.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDataWorksUserMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDataWorksUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_dataworksuser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDataWorksUserBasicDependence0)
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
					"user_id":    "${data.alibabacloudstack_ascm_users.default.users.0.primary_key}",
					"role_code":  []string{"role_project_pe", "role_project_admin"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_code.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_code": []string{"role_project_admin", "role_project_dev"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_code.#": "2",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackDataWorksUserMap0 = map[string]string{}

func AlibabacloudStackDataWorksUserBasicDependence0(name string) string {
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
`, name)
}
