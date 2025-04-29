package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDmsenterpriseUser0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_dmsenterprise_user.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccDmsenterpriseUserCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Dms_enterpriseService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDms_EnterpriseGetuserRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdms_enterpriseuser%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccDmsenterpriseUserBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"uid": "265530631068325049",

					"user_name": "rdktest",

					"mobile": "11111111111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"uid": "265530631068325049",

						"user_name": "rdktest",

						"mobile": "11111111111",
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

var AlibabacloudTestAccDmsenterpriseUserCheckmap = map[string]string{

	"status": CHECKSET,

	"uid": CHECKSET,

	"role_names": CHECKSET,

	"user_name": CHECKSET,

	"user_id": CHECKSET,

	"role_ids": CHECKSET,

	"mobile": CHECKSET,

	"parent_uid": CHECKSET,
}

func AlibabacloudTestAccDmsenterpriseUserBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
