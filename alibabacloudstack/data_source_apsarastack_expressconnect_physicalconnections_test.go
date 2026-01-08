package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectPhysicalconnectionsDataSource(t *testing.T) {

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"ids": `["${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"ids": `["${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}_fake"]`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"name_regex": `"^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}$"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"name_regex": `"^tf_fake_fake$"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"name_regex": `"^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}$"`,
			"ids": `["${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(map[string]string{
			"name_regex": `"^${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.physical_connection_name}_fake$"`,
			"ids": `["${data.alibabacloudstack_expressconnect_physical_connections.anyone.ids.0}_fake"]`,
		}),
	}
	AlibabacloudstackExpressconnectPhysicalConnectionsCheckInfo.dataSourceTestCheck(t, 0, nameConf)
	AlibabacloudstackExpressconnectPhysicalConnectionsCheckInfo.dataSourceTestCheck(t, 0, idsConf, nameConf, allConf)
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

func testAccCheckAlibabacloudstackExpressconnectPhysicalConnectionsSourceConfig(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

data "alibabacloudstack_expressconnect_physical_connections" "anyone" {
	name_regex = ".+"
	include_reservation_data = false
}

data "alibabacloudstack_expressconnect_physical_connections" "default" {
%s
}
`, strings.Join(pairs, "\n   "))
	return config
}
