package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackAscmUserGroupResourceSetBinding_Basic(t *testing.T) {
	var v *MembersInsideResourceSet
	resourceId := "alibabacloudstack_ascm_user_group_resource_set_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupResourceSetBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmresourceSetAddUusergroupconfigbasic)
	after7day := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscmUserGroupResourceSetBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_set_id":    "${alibabacloudstack_ascm_resource_group.default.rg_id}",
					"user_group_id":      "${alibabacloudstack_ascm_user_group.default.user_group_id}",
					"ascm_role_id":       "2",
					"enable_auth_expire": true,
					"expiration_time":    after7day,
					"expire_type":        "VALID_UNTIL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ascm_role_id":       CHECKSET,
						"enable_auth_expire": "true",
						"expire_type":        "VALID_UNTIL",
						"expiration_time":    after7day,
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"length_by_day": "30",
			// 		"expire_type":   "LENGTH",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"expire_type":   "LENGTH",
			// 			"length_by_day": "30",
			// 		}),
			// 	),
			// },
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_auth_expire", "length_by_day", "expiration_time", "expire_type"},
			},
		},
	})
}

func TestAccAlibabacloudStackAscmUserGroupResourceSetBinding_LengthByDays(t *testing.T) {
	var v *MembersInsideResourceSet
	resourceId := "alibabacloudstack_ascm_user_group_resource_set_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupResourceSetBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-ascmusergroup%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmresourceSetAddUusergroupconfigbasic)
	// after7day := time.Now().UTC().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscmUserGroupResourceSetBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_set_id":    "${alibabacloudstack_ascm_resource_group.default.rg_id}",
					"user_group_id":      "${alibabacloudstack_ascm_user_group_role_binding.default.user_group_id}",
					"ascm_role_id":       "2",
					"enable_auth_expire": true,
					"length_by_day":      "30",
					"expire_type":        "LENGTH",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ascm_role_id":       CHECKSET,
						"enable_auth_expire": "true",
						"length_by_day":      "30",
						"expire_type":        "LENGTH",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"expiration_time": after7day,
			// 		"expire_type":     "VALID_UNTIL",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"expiration_time": after7day,
			// 			"expire_type":     "VALID_UNTIL",
			// 		}),
			// 	),
			// },
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_auth_expire", "length_by_day", "expiration_time", "expire_type"},
			},
		},
	})
}

func testAccCheckAscmUserGroupResourceSetBindingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ascm_user_group_resource_set_binding" || rs.Type != "alibabacloudstack_ascm_user_group_resource_set_binding" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUserGroup(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.Message != "" {
			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
		}
	}

	return nil
}
func testascmresourceSetAddUusergroupconfigbasic(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}

resource "alibabacloudstack_ascm_user_group" "default" {
 group_name = "${var.name}"
}


resource "alibabacloudstack_ascm_resource_group" "default" {
  name = "${var.name}"
}

resource "alibabacloudstack_ascm_user_group_role_binding" "default" {
  role_ids = ["2"]
  user_group_id = "${alibabacloudstack_ascm_user_group.default.user_group_id}" 
}

`, name)
}

var testAccCheckUserGroupResourceSetBinding = map[string]string{
	"user_group_id":   CHECKSET,
	"resource_set_id": CHECKSET,
}
