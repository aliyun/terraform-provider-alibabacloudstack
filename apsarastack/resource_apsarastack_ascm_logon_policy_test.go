package apsarastack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackLogonPolicy_basic(t *testing.T) {
	var v *LoginPolicy
	resourceId := "apsarastack_ascm_logon_policy.default"
	ra := resourceAttrInit(resourceId, ascmLogonPolicyBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-ascmlogonpolicybasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testacclogonpolicyconfigBasic)
	resource.Test(t, resource.TestCase{
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
			//{
			//	ResourceName:      resourceId,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
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
