package apsarastack

import (
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccApsaraStackAscm_CustomRoleBasic(t *testing.T) {
	var v *AscmCustomRole

	resourceId := "apsarastack_ascm_custom_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmCustomRole)
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
		CheckDestroy:  testAccCheckAscm_CustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAscm_CustomRole_resource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_CustomRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)

	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_custom_role" || rs.Type != "apsarastack_ascm_custom_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmCustomRole(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			return WrapError(Error("custom role still exist"))
		}
	}

	return nil
}

const testAccAscm_CustomRole_resource = `
resource "apsarastack_ascm_custom_role" "default" {
	role_name = "Test_custom_Role"
	description = "TestRole"
	organization_visibility = "organizationVisibility.global"
	role_range = "roleRange.allOrganizations"
	privileges = ["PRIG_SYS_BILLING_CLOUDPRODUCTBILL_READ","PRIG_SYS_BILLING_ORGRSBILL_READ","PRIG_SYS_BILLING_BILL_EXPORT","PRIG_SYS_BILLING_BILL_MODIFY","PRIG_SYS_CHANGEOWN_READ","PRIG_SYS_CHANGEOWN_ORGANIZATION","PRIG_SYS_CHANGEOWN_RESOURCESET","PRIG_SYS_CHANGEOWN_USER","PRIG_SYS_CHANGEOWN_RESOURCE","PRIG_SYS_CHARGING_PRICE_READ","PRIG_SYS_CHARGING_PRICE_OPERATE","PRIG_SYS_CHARGING_PRICE_CREATE_DELETE","PRIG_SYS_DOWNLOAD_CENTER_TASK_READ","PRIG_SYS_DOWNLOAD_CENTER_TASK_CREATE","PRIG_SYS_DOWNLOAD_CENTER_TASK_DELETE","PRIG_SYS_DOWNLOAD_CENTER_REPORT_DOWNLOAD","PRIG_SYS_LOGINPOLICY_READ","PRIG_SYS_LOGINPOLICY_CREATE_DELETE","PRIG_SYS_LOGINPOLICY_OPERATE","PRIG_SYS_MENU_MANAGE","PRIG_SYS_METERING_READ","PRIG_SYS_METERING_EXPORT","PRIG_SYS_MSGCENTER","PRIG_SYS_OPLOG_READ","PRIG_SYS_OPLOG_OPERATE","PRIG_SYS_ORG_READ","PRIG_SYS_ORG_CREATE_DELETE","PRIG_SYS_ORG_OPERATE","PRIG_SYS_ORG_AK_READ","PRIG_SYS_QUOTA_READ","PRIG_SYS_QUOTA_OPERATE","PRIG_SYS_RESOURCESET_READ","PRIG_SYS_RESOURCESET_CREATE_DELETE","PRIG_SYS_RESOURCESET_OPERATE","PRIG_SYS_ROLE_READ","PRIG_SYS_ROLE_CREATE_DELETE","PRIG_SYS_ROLE_OPERATE","PRIG_SYS_SYSCONF","PRIG_SYS_USER_READ","PRIG_SYS_USER_CREATE_DELETE","PRIG_SYS_USER_OPERATE","PRIG_SYS_USERGROUP_READ","PRIG_SYS_USERGROUP_CREATE_DELETE","PRIG_SYS_USERGROUP_OPERATE"]
}
`

var testAccCheckAscmCustomRole = map[string]string{
	"role_name":               CHECKSET,
	"role_range":              CHECKSET,
	"organization_visibility": CHECKSET,
}
