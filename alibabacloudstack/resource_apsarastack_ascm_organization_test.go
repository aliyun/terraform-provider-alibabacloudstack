package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAscm_OrganizationBasic(t *testing.T) {
	var v *Organization

	resourceId := "alibabacloudstack_ascm_organization.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmOrg)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmorgbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testaccOrganizationBasic)
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
					"name":      "${var.name}",
					"parent_id": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
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

func testaccOrganizationBasic(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}
`, name)
}

var testAccCheckAscmOrg = map[string]string{
	"name":      CHECKSET,
	"parent_id": CHECKSET,
}
