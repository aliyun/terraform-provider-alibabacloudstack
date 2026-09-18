package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscm_CustomRoleBasic(t *testing.T) {
	var v *AscmCustomRole

	resourceId := "alibabacloudstack_ascm_custom_role.default"
	ra := resourceAttrInit(resourceId, testAccCheckAscmCustomRole)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "Tf-testAccCustomRole"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAscmCustomRoleConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscmCustomRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":               "${var.name}",
					"description":             "TestRole",
					"organization_visibility": "organizationVisibility.global",
					"role_range":              "roleRange.allOrganizations",
					"privileges": []string{
						"PRIG_SYS_BILLING_CLOUDPRODUCTBILL_READ", "PRIG_SYS_BILLING_ORGRSBILL_READ",
						"PRIG_SYS_BILLING_BILL_EXPORT", "PRIG_SYS_BILLING_BILL_MODIFY",
						"PRIG_SYS_CHANGEOWN_READ", "PRIG_SYS_CHANGEOWN_ORGANIZATION",
						"PRIG_SYS_CHANGEOWN_RESOURCESET", "PRIG_SYS_CHANGEOWN_USER",
						"PRIG_SYS_CHANGEOWN_RESOURCE", "PRIG_SYS_CHARGING_PRICE_READ",
						"PRIG_SYS_CHARGING_PRICE_OPERATE", "PRIG_SYS_CHARGING_PRICE_CREATE_DELETE",
						"PRIG_SYS_DOWNLOAD_CENTER_TASK_READ", "PRIG_SYS_DOWNLOAD_CENTER_TASK_CREATE",
						"PRIG_SYS_DOWNLOAD_CENTER_TASK_DELETE", "PRIG_SYS_DOWNLOAD_CENTER_REPORT_DOWNLOAD",
						"PRIG_SYS_LOGINPOLICY_READ", "PRIG_SYS_LOGINPOLICY_CREATE_DELETE",
						"PRIG_SYS_LOGINPOLICY_OPERATE", "PRIG_SYS_MENU_MANAGE",
						"PRIG_SYS_METERING_READ", "PRIG_SYS_METERING_EXPORT",
						"PRIG_SYS_MSGCENTER", "PRIG_SYS_OPLOG_READ",
						"PRIG_SYS_OPLOG_OPERATE", "PRIG_SYS_ORG_READ",
						"PRIG_SYS_ORG_CREATE_DELETE", "PRIG_SYS_ORG_OPERATE",
						"PRIG_SYS_ORG_AK_READ", "PRIG_SYS_QUOTA_READ",
						"PRIG_SYS_QUOTA_OPERATE", "PRIG_SYS_RESOURCESET_READ",
						"PRIG_SYS_RESOURCESET_CREATE_DELETE", "PRIG_SYS_RESOURCESET_OPERATE",
						"PRIG_SYS_ROLE_READ", "PRIG_SYS_ROLE_CREATE_DELETE",
						"PRIG_SYS_ROLE_OPERATE", "PRIG_SYS_SYSCONF",
						"PRIG_SYS_USER_READ", "PRIG_SYS_USER_CREATE_DELETE",
						"PRIG_SYS_USER_OPERATE", "PRIG_SYS_USERGROUP_READ",
						"PRIG_SYS_USERGROUP_CREATE_DELETE", "PRIG_SYS_USERGROUP_OPERATE",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name,
						"role_range":              CHECKSET,
						"organization_visibility": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"privileges"},
			},
		},
	})
}

func testAccCheckAscmCustomRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ascm_custom_role" {
			continue
		}
		ascm, err := ascmService.DescribeAscmCustomRole(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.AsapiErrorCode != "200" {
			continue
		}
		return errmsgs.WrapError(errmsgs.Error("custom role still exist"))
	}
	return nil
}

func resourceAscmCustomRoleConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

var testAccCheckAscmCustomRole = map[string]string{
	"role_name":               CHECKSET,
	"role_range":              CHECKSET,
	"organization_visibility": CHECKSET,
}
