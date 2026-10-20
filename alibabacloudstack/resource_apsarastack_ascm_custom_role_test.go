package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_ascm_custom_role", &resource.Sweeper{
		Name: "alibabacloudstack_ascm_custom_role",
		F: func(region string) error { return SweeperFunc(
			sweeperEntry{fn: testSweepAscmCustomRoles, precise: []string{"tf-testAccCustomRole"}, broad: broadPrefixes, fatal: true},
		) },
	})
}

func testSweepAscmCustomRoles(prefixes []string) error {
	rawClient, err := sharedClientForRegion("")
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRoles", "/ascm/auth/role/listRoles")
	request.QueryParams["roleType"] = "ROLETYPE_ASCM"
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("error listing ASCM custom roles: %s", err)
	}

	var resp AscmCustomRole
	if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &resp); err != nil {
		return fmt.Errorf("error unmarshaling roles: %s", err)
	}

	for _, d := range resp.Data {
		matched := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(d.RoleName, prefix) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		log.Printf("[INFO] Deleting ASCM custom role: %s", d.RoleName)
		removeReq := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRole", "/ascm/auth/role/removeRole")
		removeReq.QueryParams["roleName"] = d.RoleName
		if _, err := client.ProcessCommonRequest(removeReq); err != nil {
			log.Printf("[WARN] Failed to delete role %s: %s", d.RoleName, err)
		}
	}
	return nil
}

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
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccCustomRole%d", rand)
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
						"PRIG_SYS_ORG_READ",
						"PRIG_SYS_ROLE_READ",
						"PRIG_SYS_USER_READ",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":               name,
						"role_range":              CHECKSET,
						"organization_visibility": CHECKSET,
						"privileges.#":            "3",
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
	"privileges.#":            CHECKSET,
}
