package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisConnection0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_redis_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisConnectionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoR_KvstoreDescribedbinstancenetinfoRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredisconnection%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisConnectionBasicdependence)
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

					"instance_id": "r-8vb6ces3yk5huhxoek",

					"port": "6379",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_id": "r-8vb6ces3yk5huhxoek",

						"port": "6379",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccRedisConnectionCheckmap = map[string]string{

	"status": CHECKSET,

	"instance_id": CHECKSET,

	"create_time": CHECKSET,

	"port": CHECKSET,

	"vswitch_id": CHECKSET,

	"vpc_id": CHECKSET,

	"expired_time": CHECKSET,

	"db_instance_net_type": CHECKSET,

	"ip_address": CHECKSET,

	"connection_string": CHECKSET,
}

func AlibabacloudTestAccRedisConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
