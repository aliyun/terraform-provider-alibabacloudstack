package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_E_OrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscm_e_Organization_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
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
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.RequestID != "" {
			return WrapError(Error("organization still exist"))
		}
	}

	return nil
}

const testAccAscm_e_Organization_resource = `
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org"
  parent_id = "1"
}`

var testAccCheckAscmOrg = map[string]string{
	"name":      CHECKSET,
	"parent_id": CHECKSET,
}
