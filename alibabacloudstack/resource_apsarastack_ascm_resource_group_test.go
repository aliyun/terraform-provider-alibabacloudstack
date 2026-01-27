package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackascmResourceGroupBasic(t *testing.T) {
	var v *ResourceGroupData
	resourceId := "alibabacloudstack_ascm_resource_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckResourceGroup)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmregp%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscmResourceGroup)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		// CheckDestroy: testAccCheckAscm_Resource_GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            "${var.name}",
					"organization_id": "${alibabacloudstack_ascm_organization.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":  name,
						"rg_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_update",
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

// func testAccCheckAscm_Resource_GroupDestroy(s *terraform.State) error { //destroy function
// 	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
// 	ascmService := AscmService{client}

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "alibabacloudstack_ascm_resource_group" {
// 			continue
// 		}
// 		ascm, err := ascmService.DescribeAscmResourceGroup(rs.Primary.ID)
// 		if err != nil {
// 			if errmsgs.NotFoundError(err) {
// 				continue
// 			}
// 			return errmsgs.WrapError(err)
// 		}
// 		if len(ascm.Data) > 0 {
// 			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
// 		}
// 	}

// 	return nil
// }

func testAccAscmResourceGroup(name string) string {
	return fmt.Sprintf(`
variable name{
	default = "%s"
}

resource "alibabacloudstack_ascm_organization" "default" {
  name = "${var.name}"
  parent_id = "1"
} 

`, name)
}

var testAccCheckResourceGroup = map[string]string{
	"name":            CHECKSET,
	"organization_id": CHECKSET,
}
