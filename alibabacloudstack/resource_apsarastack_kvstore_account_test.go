package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisAccount0(t *testing.T) {
	var v *r_kvstore.Account

	resourceId := "alibabacloudstack_redis_account.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisAccountCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR_KvstoreDescribeaccountsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-redisaccount%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisAccountBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,

		CheckDestroy: nil,

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"description": "rdk_test_description",
					"instance_id": "${local.kv_instance_id}",
					"account_name": "rdk_test_name_01",
					"account_password": "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "rdk_test_description",
						"account_name": "rdk_test_name_01",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Sensitive information is not read back
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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
		},
	})
}

var AlibabacloudTestAccRedisAccountCheckmap = map[string]string{
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

%s

`, name,  KVRInstanceCommonTestCase("enterprise", string(KVStoreRedis)))
}
