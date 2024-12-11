package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisAccount0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR_KvstoreDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredisaccount%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisAccountBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "rdk_test_description",

					"instance_id": "r-bp1db1a29f56e904",

					"account_name": "rdk_test_name_01",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "rdk_test_description",

						"instance_id": "r-bp1db1a29f56e904",

						"account_name": "rdk_test_name_01",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"account_privilege": "RoleReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"account_privilege": "RoleReadWrite",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "testescription",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "testescription",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRedisAccountCheckmap = map[string]string{

	"status": CHECKSET,

	"description": CHECKSET,

	"account_privilege": CHECKSET,

	"instance_id": CHECKSET,

	"account_type": CHECKSET,

	"account_name": CHECKSET,
}

func AlibabacloudTestAccRedisAccountBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
