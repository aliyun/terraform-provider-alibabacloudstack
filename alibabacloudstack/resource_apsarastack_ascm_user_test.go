package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccAlibabacloudStackAscm_UserBasic(t *testing.T) {
	var v *User
	var org_id string
	resourceId := "alibabacloudstack_ascm_user.default"
	ra := resourceAttrInit(resourceId, ascmuserBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-ascmusers%v", rand)
	if os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT") != "" {
		org_id = os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT")
	} else {
		org_id = "${alibabacloudstack_ascm_organization.default.org_id}"
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmuserconfigbasic)
	resource.Test(t, resource.TestCase{
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
					"cellphone_number":   "13900000000",
					"email":              "test01@gmail.com",
					"display_name":       "Test-Apsara",
					"organization_id":    org_id,
					"mobile_nation_code": "91",
					"login_name":         name,
					"login_policy_id":    "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
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
	if os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT") != "" {
		return fmt.Sprintf(`
variable name{
 default = "%s"
}
`, name)
	} else {
		return fmt.Sprintf(`
variable name{
 default = "%s"
}
resource "alibabacloudstack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}
`, name)
	}
}

var ascmuserBasicMap = map[string]string{
	"cellphone_number":   CHECKSET,
	"email":              CHECKSET,
	"display_name":       CHECKSET,
	"organization_id":    CHECKSET,
	"mobile_nation_code": CHECKSET,
	"login_name":         CHECKSET,
}
