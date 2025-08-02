package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbShardingInstanceCsNodeAddress(t *testing.T) {
	var v map[string]interface {}

	resourceId := "alibabacloudstack_mongodb_shardinginstance_csnode_address.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"enable_public_connection":  CHECKSET,
		"enable_private_connection": CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeShardingInstanceNode")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfmongodb_csnode%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccMongodbShardingInstanceNodeAddressBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":            "${local.shard_instance.id}",
					"node_id":                   "${[for cs in local.shard_instance.configserver_list : cs if cs.description == \"cs1\"][0].node_id}",
					"enable_public_connection":  true,
					"enable_private_connection": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_connect_string":  CHECKSET,
						"public_connect_port":    CHECKSET,
						"private_connect_string": CHECKSET,
						"private_connect_port":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"db_instance_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection":  false,
					"enable_private_connection": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_connect_string":  REMOVEKEY,
						"public_connect_port":    REMOVEKEY,
						"private_connect_string": REMOVEKEY,
						"private_connect_port":   REMOVEKEY,
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccMongodbShardingInstanceNodeAddressBasicdependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

variable "existed_db_instance_id" {
	default = "%s"
}

data "alibabacloudstack_mongodb_instance_types" "mongos" {
  db_instnace_type = "sharding"
  node_type        = "mongos"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

data "alibabacloudstack_mongodb_instance_types" "configserver" {
  db_instnace_type = "sharding"
  node_type        = "configserver"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

data "alibabacloudstack_mongodb_instance_types" "shard" {
  db_instnace_type = "sharding"
  node_type        = "shard"
  sorted_by        = "CPU"
  engine_version   = "4.0"
}

resource "random_password" "password" {
  count            = 2
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

data "alibabacloudstack_mongodb_instances" "default" {
  ids = var.existed_db_instance_id == "" ? [] : ["${var.existed_db_instance_id}",]
  instance_type = "sharding"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count = length(data.alibabacloudstack_mongodb_instances.default.instances) == 0 ? 1 : 0
  db_account_password = random_password.password.0.result
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  engine_version      = "3.4"
  shard_list {
    node_storage = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.storage_min
    description  = "shard1"
    node_class   = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.id
  }
  shard_list {
    description  = "shard2"
    node_class   = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.id
    node_storage = data.alibabacloudstack_mongodb_instance_types.shard.instance_types.0.storage_min
  }

  mongo_list {
    node_class  = data.alibabacloudstack_mongodb_instance_types.mongos.instance_types.0.id
    description = "mongo1"
  }
  mongo_list {
    description = "mongo2"
    node_class  = data.alibabacloudstack_mongodb_instance_types.mongos.instance_types.0.id
  }

  configserver_list {
    node_storage = data.alibabacloudstack_mongodb_instance_types.configserver.instance_types.0.storage_min
    description  = "cs1"
    node_class   = data.alibabacloudstack_mongodb_instance_types.configserver.instance_types.0.id
  }

	db_account_name = "tf_testacc"

}

locals {
    shard_instance = length(data.alibabacloudstack_mongodb_instances.default.instances) == 0 ? alibabacloudstack_mongodb_sharding_instance.default.0 : data.alibabacloudstack_mongodb_instances.default.instances.0
}

`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_MONGOS_SHARD_ID"))
}
