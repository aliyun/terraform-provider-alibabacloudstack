package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackMongodbShardingInstanceShardNodeAddress(t *testing.T) {
	var v map[string]interface {}

	resourceId := "alibabacloudstack_mongodb_shardinginstance_shardnode_address.default"
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
	name := fmt.Sprintf("tfmongodb_shardNode%d", rand)

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
					"node_id":                   "${[for shard in local.shard_instance.shard_list : shard if shard.description == \"shard1\"][0].node_id}",
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

