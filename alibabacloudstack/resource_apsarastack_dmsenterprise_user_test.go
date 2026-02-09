package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkDmsUserDestory(rc *resourceCheck) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		module := s.RootModule()
		if module == nil {
			return fmt.Errorf("root module not found in state")
		}

		resource := module.Resources[rc.resourceId]
		if resource == nil {
			return fmt.Errorf("resource %q not found in state", rc.resourceId)
		}

		if resource.Primary == nil {
			return fmt.Errorf("primary instance not found for resource %q", rc.resourceId)
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		dmsService := DmsService{client}
		object, err := dmsService.DescribeDmsEnterpriseUser(resource.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil
			}
			return err
		}

		if object["State"].(string) == "DELETE" {
			return nil
		}

		return errmsgs.WrapError(errmsgs.Error("the resource %s %s was not destroyed ! ", rc.resourceId, resource.Primary.ID))
	}
}
func TestAccAlibabacloudStackDmsenterpriseUser0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_dmsenterprise_user.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDmsenterpriseUserCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDms_EnterpriseGetuserRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-dmsuser%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDmsenterpriseUserBasicdependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: checkDmsUserDestory(rc),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"uid":               "${alibabacloudstack_ascm_user.user.user_uid}",
					"user_name":         "${alibabacloudstack_ascm_user.user.login_name}",
					"max_execute_count": 10,
					"max_result_count":  10,
					"role_names":        []string{"USER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_names.#": "1",
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

func TestAccAlibabacloudStackDmsenterpriseUser1(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_dmsenterprise_user.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDmsenterpriseUserCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDms_EnterpriseGetuserRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-dmsuser%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDmsenterpriseUserBasicdependence1)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: checkDmsUserDestory(rc),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"uid":               "${data.alibabacloudstack_ascm_users.default.users.0.primary_key}",
					"user_name":         "${data.alibabacloudstack_ascm_users.default.users.0.login_name}",
					"max_execute_count": 10,
					"max_result_count":  10,
					"role_names":        []string{"USER"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
					"role_names": []string{"USER", "ADMIN"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"status": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlibabacloudTestAccDmsenterpriseUserCheckmap = map[string]string{
	"status":       CHECKSET,
	"uid":          CHECKSET,
	"role_names.#": CHECKSET,
	"user_name":    CHECKSET,
}

func AlibabacloudTestAccDmsenterpriseUserBasicdependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alibabacloudstack_account" "current" {
}

resource "alibabacloudstack_ascm_user" "user" {
 cellphone_number = "13900000000"
 email = "test@gmail.com"
 display_name = "${var.name}"
 organization_id = data.alibabacloudstack_account.current.organization_id
 mobile_nation_code = "86"
 login_name = "${var.name}"
 login_policy_id = 1
}

`, name)
}

func AlibabacloudTestAccDmsenterpriseUserBasicdependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alibabacloudstack_ascm_users" "default" {
}

`, name)
}
