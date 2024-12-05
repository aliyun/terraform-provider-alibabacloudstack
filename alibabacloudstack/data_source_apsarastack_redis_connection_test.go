package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackRedisConnectionsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_connections.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_connections.default.id}_fake"]`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_connections.default.id}"]`,
			"instance_id": `"${alibabacloudstack_redis_connections.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_redis_connections.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_redis_connections.default.InstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_connections.default.id}"]`,

			"instance_id": `"${alibabacloudstack_redis_connections.default.InstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_redis_connections.default.id}_fake"]`,

			"instance_id": `"${alibabacloudstack_redis_connections.default.InstanceId}_fake"`}),
	}

	AlibabacloudstackRedisConnectionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_idConf, allConf)
}

var existAlibabacloudstackRedisConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#":    "1",
		"connections.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackRedisConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#": "0",
	}
}

var AlibabacloudstackRedisConnectionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_redis_connections.default",
	existMapFunc: existAlibabacloudstackRedisConnectionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackRedisConnectionsMapFunc,
}

func testAccCheckAlibabacloudstackRedisConnectionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackRedisConnections%d"
}






data "alibabacloudstack_redis_connections" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
