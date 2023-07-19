package alibabacloudstack

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackKvstoreConnection_basic(t *testing.T) {
	var v r_kvstore.InstanceNetInfo
	resourceId := "alibabacloudstack_kvstore_connection.default"
	ra := resourceAttrInit(resourceId, RedisConnectionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKvstoreConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKvstoreConnection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KvstoreConnectionBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetest",
					"instance_id":              "${alibabacloudstack_kvstore_instance.default.id}",
					"port":                     "6370",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":              CHECKSET,
						"port":                     "6370",
						"connection_string_prefix": "allocatetest",
					}),
				),
			},
			//{
			//	ResourceName:            resourceId,
			//	ImportState:             true,
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"connection_string_prefix"},
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "allocatetestupdate",
					"port":                     "6371",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": "allocatetestupdate",
						"port":                     "6371",
					}),
				),
			},
		},
	})
}

var RedisConnectionMap = map[string]string{}

func KvstoreConnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
provider "alibabacloudstack" {
	assume_role {}
}
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstancesDataSource2"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vpc" "default" {
name       = var.name
cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}
	
resource "alibabacloudstack_kvstore_instance" "default" {
instance_name = "%s"
instance_class = "redis.master.stand.default"
vswitch_id     = alibabacloudstack_vswitch.default.id
private_ip     = "172.16.0.10"
security_ips   = ["10.0.0.1"]
instance_type  = "Redis"
engine_version = "4.0"	
 architecture_type = "standard"
cpu_type = "intel"
}
	`, name)
}
