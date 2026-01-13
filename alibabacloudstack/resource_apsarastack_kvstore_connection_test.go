package alibabacloudstack

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisConnection0(t *testing.T) {
	var v r_kvstore.InstanceNetInfo

	resourceId := "alibabacloudstack_redis_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisConnectionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKvstoreConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-kvcon%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisConnectionBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":              "${local.kv_instance_id}",
					"connection_string_prefix": name,
					"port": "6379",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "6379",
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

var AlibabacloudTestAccRedisConnectionCheckmap = map[string]string{

	"instance_id": CHECKSET,
	"port": CHECKSET,
	"connection_string": CHECKSET,
}

func AlibabacloudTestAccRedisConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}

	%s

	`, name, KVRInstanceCommonTestCase("enterprise", string(KVStoreRedis), ))
}
