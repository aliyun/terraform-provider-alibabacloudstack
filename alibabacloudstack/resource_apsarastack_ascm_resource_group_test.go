package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAlibabacloudStackAscm_Resource_GroupBasic(t *testing.T) {
	var v *ResourceGroup
	resourceId := "alibabacloudstack_ascm_resource_group.default"
	ra := resourceAttrInit(resourceId, testAccCheckResourceGroup)
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
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_resource_group" || rs.Type != "alibabacloudstack_ascm_resource_group" {
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
resource "alibabacloudstack_ascm_organization" "default" {
  name = "Tf-testingresource-org"
  parent_id = "1"
} 
 resource "alibabacloudstack_ascm_resource_group" "default" {
  organization_id = alibabacloudstack_ascm_organization.default.org_id
  name = "alibabacloudstack-Datasource-resourceGroup"
}`

var testAccCheckResourceGroup = map[string]string{
	"name":            CHECKSET,
	"organization_id": CHECKSET,
}
