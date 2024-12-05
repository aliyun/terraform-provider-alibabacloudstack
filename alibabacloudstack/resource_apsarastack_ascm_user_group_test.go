package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccAlibabacloudStackAscm_User_Group_Basic(t *testing.T) {
	var v *UserGroup
	var org_id string
	resourceId := "alibabacloudstack_ascm_user_group.default"
	ra := resourceAttrInit(resourceId, ascmusergroupBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
	if os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT") != "" {
		org_id = os.Getenv("ALIBABACLOUDSTACK_DEPARTMENT")
	} else {
		org_id = "${alibabacloudstack_ascm_organization.default.org_id}"
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmusergroupconfigbasic)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_User_Group_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":      name,
					"organization_id": org_id,
					"role_in_ids":     []string{"2", "6"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_User_Group_Destroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if true {
			continue
		}
		_, err := ascmService.DescribeAscmUserGroup(rs.Primary.ID)
		if err == nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
	}

	return nil
}
func testascmusergroupconfigbasic(name string) string {
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

var ascmusergroupBasicMap = map[string]string{
	"group_name":      CHECKSET,
	"organization_id": CHECKSET,
}
