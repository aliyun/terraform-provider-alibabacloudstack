package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackLogonPolicy_basic(t *testing.T) {
	var v *LoginPolicy
	resourceId := "alibabacloudstack_ascm_logon_policy.default"
	ra := resourceAttrInit(resourceId, ascmLogonPolicyBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-ascmlogonpolicybasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testacclogonpolicyconfigBasic)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		//CheckDestroy: testAccCheckLogonPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "testing purpose",
					"rule":        "ALLOW",
					"name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
func testacclogonpolicyconfigBasic(name string) string {
	return fmt.Sprintf(`
variable logonPolicy{
 default = "%s"
}
`, name)
}

var ascmLogonPolicyBasicMap = map[string]string{
	"description": CHECKSET,
	"rule":        CHECKSET,
	"name":        CHECKSET,
}
