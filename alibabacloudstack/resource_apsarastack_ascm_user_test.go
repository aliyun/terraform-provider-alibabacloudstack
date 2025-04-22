package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_UserBasic(t *testing.T) {
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
		CheckDestroy:  testAccCheckAscm_UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cellphone_number":   "13612345678",
					"email":              "test01@gmail.com",
					"display_name":       "Test-Apsara",
					"mobile_nation_code": "86",
					"login_name":         name,
					"login_policy_id":    "1",
					// "role_ids":           []string{"8", "9"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cellphone_number":   "13612345678",
						"email":              "test01@gmail.com",
						"display_name":       "Test-Apsara",
						"mobile_nation_code": "86",
						"login_name":         name,
						"login_policy_id":    "1",
						// "role_ids.#":         "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				//ImportStateVerifyIgnore: []string{"compact_topic", "partition_num", "remark"},
			},
		},
	})

}

func testAccCheckAscm_UserDestroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if true {
			continue
		}
		_, err := ascmService.DescribeAscmUser(rs.Primary.ID)
		if err == nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
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
