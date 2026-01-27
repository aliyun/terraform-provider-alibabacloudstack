package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmOrganizationBasic(t *testing.T) {
	var v *OrganizationResponse

	resourceId := "alibabacloudstack_ascm_organization.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmOrg)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmorg%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAscm_e_Organization_resource)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_E_OrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":      "${var.name}",
					"parent_id": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"parent_id": "1",
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

func testAccCheckAscm_E_OrganizationDestroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_organization" || rs.Type != "alibabacloudstack_ascm_organization" {
			continue
		}
		ascm, err := ascmService.DescribeAscmOrganization(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.RequestID != "" {
			return errmsgs.WrapError(errmsgs.Error("organization still exist"))
		}
	}

	return nil
}

func testAccAscm_e_Organization_resource(name string) string {
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
