package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbShardingNetworkPublicAddress0(t *testing.T) {
	var v *DdsDescribeshardingnetworkaddressResponse

	resourceId := "alibabacloudstack_mongodb_sharding_network_publicaddrdss.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccMongodbShardingNetworkPublicAddressCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoDdsDescribeshardingnetworkaddressRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfmongodb_sharding_network_publicaddrdss%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbShardingNetworkPublicAddressBasicdependence)
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

					"db_instance_id": "${alibabacloudstack_mongodb_sharding_instance.default.id}",

					"node_id": "${alibabacloudstack_mongodb_sharding_instance.default.mongo_list.0.node_id}",

					// "instance_id": "${alibabacloudstack_mongodb_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"network_addresses.#":                 "1",
						"network_addresses.0.network_address": CHECKSET,
						"network_addresses.0.port":            CHECKSET,
						"network_addresses.0.node_id":         CHECKSET,
						"network_addresses.0.role":            CHECKSET,
						"network_addresses.0.ip_address":      CHECKSET,
						"network_addresses.0.network_type":    CHECKSET,
					}),
				),
			},

			// {
			// 	Config: testAccConfig(map[string]interface{}{

			// 		"account_password": modify_passward,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"account_name": name,
			// 		}),
			// 	),
			// },

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_instance_id", "node_id"},
			},
		},
	})
}

var AlibabacloudTestAccMongodbShardingNetworkPublicAddressCheckmap = map[string]string{}

func AlibabacloudTestAccMongodbShardingNetworkPublicAddressBasicdependence(name string) string {
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
