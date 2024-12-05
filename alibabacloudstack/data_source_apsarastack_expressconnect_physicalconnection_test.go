package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackExpressconnectPhysicalConnectionsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_physical_connections.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_physical_connections.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_physical_connections.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_physical_connections.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackExpressconnectPhysicalConnectionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackExpressconnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#":    "1",
		"connections.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackExpressconnectPhysicalConnectionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"connections.#": "0",
	}
}

var AlibabacloudstackExpressconnectPhysicalConnectionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_expressconnect_physical_connections.default",
	existMapFunc: existAlibabacloudstackExpressconnectPhysicalConnectionsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackExpressconnectPhysicalConnectionsMapFunc,
}

func testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackExpressconnectPhysicalConnections%d"
}






data "alibabacloudstack_expressconnect_physical_connections" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
