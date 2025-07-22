package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbShardingNetworkPrivateAddress0(t *testing.T) {
	var v *DdsDescribeshardingnetworkaddressResponse

	resourceId := "alibabacloudstack_mongodb_sharding_network_private_address.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccMongodbShardingNetworkPrivateAddressCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDdsDescribeshardingnetworkaddressRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	// passward := getAccTestPassword(12)
	passward := "123@123qwe"

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfmongodb_sharding_network_privateaddrdss%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbShardingNetworkPrivateAddressBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"account_name":     "terraform",
					"account_password": passward,
					"db_instance_id":   "${alibabacloudstack_mongodb_sharding_instance.default.id}",

					"node_id": "${alibabacloudstack_mongodb_sharding_instance.default.shard_list.0.node_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_addresses.#":                 "2",
						"network_addresses.0.network_address": CHECKSET,
						"network_addresses.0.port":            CHECKSET,
						"network_addresses.0.node_id":         CHECKSET,
						"network_addresses.0.role":            CHECKSET,
						"network_addresses.0.ip_address":      CHECKSET,
						"network_addresses.0.network_type":    CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_instance_id", "node_id", "port", "network_address", "account_name", "account_password"},
			},
		},
	})
}

var AlibabacloudTestAccMongodbShardingNetworkPrivateAddressCheckmap = map[string]string{}

func AlibabacloudTestAccMongodbShardingNetworkPrivateAddressBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "MongoDB"
  }

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
	zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "3.4"
	shard_list {
	  node_class   = "dds.shard.mid"
	  node_storage = 10
	}
	shard_list {
	  node_class   = "dds.shard.standard"
	  node_storage = 20
	}
	mongo_list {
	  node_class = "dds.mongos.mid"
	}
	mongo_list {
	  node_class = "dds.mongos.mid"
	}
}

`, name)
}
