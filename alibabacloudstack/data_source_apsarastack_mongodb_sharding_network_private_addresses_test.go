package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMongodbShardingNetworkPrivateAddressesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	passward := getAccTestPassword(12)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_private_address.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_private_address.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
		}),
	}

	nodeidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_private_address.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_private_address.default.node_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_private_address.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_private_address.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand, passward, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_private_address.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_private_address.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_private_address.default.node_id}_fake"`,
		}),
	}
	AlibabacloudstackMongodbShardingNetworkPrivateAddressesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, nodeidConf, allConf)
}

var existAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"network_addresses.#":                 "2",
		"network_addresses.0.network_address": CHECKSET,
		"network_addresses.0.port":            CHECKSET,
		"network_addresses.0.node_id":         CHECKSET,
		"network_addresses.0.role":            CHECKSET,
		"network_addresses.0.ip_address":      CHECKSET,
		"network_addresses.0.network_type":    CHECKSET,
	}
}

var fakeAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"network_addresses.#": "0",
	}
}

var AlibabacloudstackMongodbShardingNetworkPrivateAddressesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_mongodb_sharding_network_private_addresses.default",
	existMapFunc: existAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataMapFunc,
}

func testAccCheckAlibabacloudstackMongodbShardingNetworkPrivateAddressesDataSourceConfig(rand int, password string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackMongodbShardingNetworkPrivateAddresses%d"
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

resource "alibabacloudstack_mongodb_sharding_network_private_address" "default" {
	db_instance_id="${alibabacloudstack_mongodb_sharding_instance.default.id}"
	account_name="terraform"
	account_password = "%s"
	node_id="${alibabacloudstack_mongodb_sharding_instance.default.mongo_list.0.node_id}"
}

data "alibabacloudstack_mongodb_sharding_network_private_addresses" "default" {
	%s
}

`, rand, password, strings.Join(pairs, "\n   "))
	return config
}
