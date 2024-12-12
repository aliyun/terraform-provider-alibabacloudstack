package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackSlbListenersDataSource(t *testing.T) {
	// 根据test_meta自动生成的tasecase

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_listeners.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_listeners.default.id}_fake"]`,
		}),
	}

	listener_protocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_listeners.default.id}"]`,
			"listener_protocol": `"${alibabacloudstack_slb_listeners.default.ListenerProtocol}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_listeners.default.id}_fake"]`,
			"listener_protocol": `"${alibabacloudstack_slb_listeners.default.ListenerProtocol}_fake"`,
		}),
	}

	load_balancer_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_listeners.default.id}"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_listeners.default.LoadBalancerId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_slb_listeners.default.id}_fake"]`,
			"load_balancer_id": `"${alibabacloudstack_slb_listeners.default.LoadBalancerId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_listeners.default.id}"]`,

			"listener_protocol": `"${alibabacloudstack_slb_listeners.default.ListenerProtocol}"`,
			"load_balancer_id":  `"${alibabacloudstack_slb_listeners.default.LoadBalancerId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_listeners.default.id}_fake"]`,

			"listener_protocol": `"${alibabacloudstack_slb_listeners.default.ListenerProtocol}_fake"`,
			"load_balancer_id":  `"${alibabacloudstack_slb_listeners.default.LoadBalancerId}_fake"`}),
	}

	AlibabacloudstackSlbListenersDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, listener_protocolConf, load_balancer_idConf, allConf)
}

var existAlibabacloudstackSlbListenersDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"listeners.#":    "1",
		"listeners.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbListenersDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"listeners.#": "0",
	}
}

var AlibabacloudstackSlbListenersDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_listeners.default",
	existMapFunc: existAlibabacloudstackSlbListenersDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbListenersDataMapFunc,
}

func testAccCheckAlibabacloudstackSlbListenersDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbListeners%d"
}






data "alibabacloudstack_slb_listeners" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

