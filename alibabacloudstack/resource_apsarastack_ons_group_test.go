package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackOnsGroup_basic(t *testing.T) {
	var v *OnsGroup
	resourceId := "alibabacloudstack_ons_group.default"
	ra := resourceAttrInit(resourceId, onsGroupBasicMap)
	serviceFunc := func() interface{} {
		return &OnsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-groupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccOnsGroupConfigBasic)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id": "${alibabacloudstack_ons_instance.default.id}",
					"group_id":    "GID-${var.name}",
					"remark":      "Ons_Group",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_id": "GID-" + name,
						"remark":   "Ons_Group",
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

func testAccOnsGroupConfigBasic(name string) string {
	return fmt.Sprintf(`

variable "name" {
 default = "%s"
}

%s

`, name, OnsCommonTestCase)
}

var onsGroupBasicMap = map[string]string{
	"instance_id": CHECKSET,
	"group_id":    CHECKSET,
	"remark":      CHECKSET,
}
