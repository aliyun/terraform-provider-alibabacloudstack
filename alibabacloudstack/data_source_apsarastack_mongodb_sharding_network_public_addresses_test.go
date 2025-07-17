package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackMongodbShardingNetworkPublicAddressesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
		}),
	}

	nodeidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.node_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.id}"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.node_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.id}_fake"]`,
			"db_instance_id": `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.db_instance_id}"`,
			"node_id":        `"${alibabacloudstack_mongodb_sharding_network_publicaddrdss.default.node_id}_fake"`,
		}),
	}
	AlibabacloudstackMongodbShardingNetworkPublicAddressesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, nodeidConf, allConf)
}

var existAlibabacloudstackMongodbShardingNetworkPublicAddressesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"network_addresses.#":                 "1",
		"network_addresses.0.network_address": CHECKSET,
		"network_addresses.0.port":            CHECKSET,
		"network_addresses.0.node_id":         CHECKSET,
		"network_addresses.0.role":            CHECKSET,
		"network_addresses.0.ip_address":      CHECKSET,
		"network_addresses.0.network_type":    CHECKSET,
	}
}

var fakeAlibabacloudstackMongodbShardingNetworkPublicAddressesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"network_addresses.#": "0",
	}
}

var AlibabacloudstackMongodbShardingNetworkPublicAddressesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_mongodb_sharding_network_publicaddrdsses.default",
	existMapFunc: existAlibabacloudstackMongodbShardingNetworkPublicAddressesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackMongodbShardingNetworkPublicAddressesDataMapFunc,
}

func testAccCheckAlibabacloudstackMongodbShardingNetworkPublicAddressesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackMongodbShardingNetworkPublicAddresses%d"
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


resource "alibabacloudstack_mongodb_sharding_network_publicaddrdss" "default" {
	db_instance_id="${alibabacloudstack_mongodb_sharding_instance.default.id}"

	node_id="${alibabacloudstack_mongodb_sharding_instance.default.mongo_list.0.node_id}"
}

data "alibabacloudstack_mongodb_sharding_network_publicaddrdsses" "default" {
	%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
