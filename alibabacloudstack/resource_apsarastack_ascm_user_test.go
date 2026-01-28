package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscmUserBasic(t *testing.T) {
	var v *User
	resourceId := "alibabacloudstack_ascm_user.default"
	ra := resourceAttrInit(resourceId, ascmuserBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmusers%v", rand)
	// if os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT") != "" {
	// 	org_id = os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT")
	// } else {
	// 	org_id = "${alibabacloudstack_ascm_organization.default.org_id}"
	// }
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmuserconfigbasic)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cellphone_number":   "13612345678",
					"email":              "test01@gmail.com",
					"display_name":       "Test-Apsara",
					"mobile_nation_code": "86",
					"login_name":         name,
					"login_policy_id":    "1",
					"role_ids":           []string{"8", "9"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cellphone_number":   "13612345678",
						"email":              "test01@gmail.com",
						"display_name":       "Test-Apsara",
						"mobile_nation_code": "86",
						"login_name":         name,
						"login_policy_id":    "2",
						"role_ids.#":         "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cellphone_number":   "13600000000",
					"email":              "test02@gmail.com",
					"display_name":       "Test-Apsara1",
					"mobile_nation_code": "85",
					"login_name":         name + "_update",
					"login_policy_id":    "2",
					"role_ids":           []string{"8", "9"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cellphone_number":   "13600000000",
						"email":              "test02@gmail.com",
						"display_name":       "Test-Apsara1",
						"mobile_nation_code": "85",
						"login_name":         name + "_update",
						"login_policy_id":    "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_ids": []string{"2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_ids.#": "1",
						"role_ids.0": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// init_password is only returned once during creation, and will not be returned during subsequent imports
				ImportStateVerifyIgnore: []string{"init_password"},
			},
		},
	})

}

func testascmuserconfigbasic(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}
`, name)
}

var ascmuserBasicMap = map[string]string{
	"cellphone_number":   CHECKSET,
	"email":              CHECKSET,
	"display_name":       CHECKSET,
	"mobile_nation_code": CHECKSET,
	"login_name":         CHECKSET,
}
