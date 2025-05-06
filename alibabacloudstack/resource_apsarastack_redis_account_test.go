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
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: nil,

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "rdk_test_description",

					"instance_id": "${alibabacloudstack_kvstore_instance.default.id}",

					"account_name": "rdk_test_name_01",

					"account_password": "1qaz@WSX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "rdk_test_description",

						"account_name": "rdk_test_name_01",

						"account_password": "1qaz@WSX",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// 敏感信息不回读
				ImportStateVerifyIgnore: []string{"account_password"},
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

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "%s"
}

%s

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
  }

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "%s"
	node_type = "double"
	architecture_type = "standard"
	password       = "1qaz@WSX"
}



`, name, string(KVStoreRedis), KVRInstanceClassCommonTestCase, string(KVStore4Dot0))
}
