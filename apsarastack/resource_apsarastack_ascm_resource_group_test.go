package apsarastack

import (
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccApsaraStackAscm_Resource_GroupBasic(t *testing.T) {
	var v *ResourceGroup
	resourceId := "apsarastack_ascm_resource_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckResourceGroup)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
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
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscm_Resource_GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscmResource_Group_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_Resource_GroupDestroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_resource_group" || rs.Type != "apsarastack_ascm_resource_group" {
			continue
		}
		ascm, err := ascmService.DescribeAscmResourceGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.Message != "" {
			return WrapError(Error("resource  still exist"))
		}
	}

	return nil
}

const testAccAscmResource_Group_resource = `
resource "apsarastack_ascm_organization" "default" {
  name = "Tf-testingresource-org"
  parent_id = "1"
} 
 resource "apsarastack_ascm_resource_group" "default" {
  organization_id = apsarastack_ascm_organization.default.org_id
  name = "apsarastack-Datasource-resourceGroup"
}`

var testAccCheckResourceGroup = map[string]string{
	"name":            CHECKSET,
	"organization_id": CHECKSET,
}
