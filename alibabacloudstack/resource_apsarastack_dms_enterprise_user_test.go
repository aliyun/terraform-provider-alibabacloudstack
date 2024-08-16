package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_dms_enterprise_user", &resource.Sweeper{
		Name: "alibabacloudstack_dms_enterprise_user",
		F:    testSweepDMSEnterpriseUsers,
	})
}

func testSweepDMSEnterpriseUsers(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"testacc",
	}
	request := map[string]interface{}{
		"UserState":  "NORMAL",
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	var response map[string]interface{}
	action := "ListUsers"
	conn, err := client.NewDmsenterpriseClient()
	if err != nil {
		return WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_dms_enterprise_users", action, AlibabacloudStackSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.UserList.User", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.UserList.User", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["NickName"]; !ok {
				skip = false
			} else {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(fmt.Sprintf("%v", item["NickName"])), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping DMS Enterprise User: %v", item["NickName"])
				continue
			}
			action := "DeleteUser"
			request := map[string]interface{}{
				"Uid": item["Uid"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-11-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				log.Printf("[ERROR] Failed to delete DMS Enterprise User (%v): %s", item["NickName"], err)
				continue
			}

			log.Printf("[INFO] Delete DMS Enterprise User Success: %v ", item["NickName"])
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

func TestAccAlibabacloudStackDMSEnterpriseUser_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dms_enterprise_user.default"
	ra := resourceAttrInit(resourceId, DmsEnterpriseUserMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDmsEnterpriseUser")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDmsEnterpriseUser%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DmsEnterpriseUserBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"uid":        "${alibabacloudstack_ascm_user.user.user_id}",
					"nick_name":  name,
					"role_names": []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name":    name,
						"role_names.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_execute_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_result_count": "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_result_count": "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"nick_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nick_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"role_names": []string{"USER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "NORMAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_execute_count": "1000",
					"max_result_count":  "1000",
					"nick_name":         name + "change",
					"role_names":        []string{"DBA"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_execute_count": "1000",
						"max_result_count":  "1000",
						"nick_name":         name + "change",
						"role_names.#":      "1",
					}),
				),
			},
		},
	})
}

var DmsEnterpriseUserMap = map[string]string{
	"status": CHECKSET,
}

func DmsEnterpriseUserBasicdependence(name string) string {
	return fmt.Sprintf(`
//   data "alibabacloudstack_ascm_organizations" "default" {
//     name_regex = "cxt"
//}
	resource "alibabacloudstack_ascm_organization" "default" {
	 name = "Test_binder124"
	 parent_id = "1"
	}
	
	resource "alibabacloudstack_ascm_user" "user" {
	 cellphone_number = "13900000000"
	 email = "admin@ascm.com"
	 display_name = "admin"
	 organization_id = alibabacloudstack_ascm_organization.default.org_id
	 mobile_nation_code = "91"
	 login_name = "User_Role_Test%s"
	 login_policy_id = 1
	}
	
`, name)
}
