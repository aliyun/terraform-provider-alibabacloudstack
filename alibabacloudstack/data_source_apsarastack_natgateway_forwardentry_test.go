package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNatgatewayForwardEntriesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_natgateway_forward_entries.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackNatgatewayForwardEntries%d", rand),
		dataSourceNatgatewayForwardEntriesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
		}),
	}

	external_ipConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"external_ip": "${alibabacloudstack_natgateway_forward_entries.default.ExternalIp}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"external_ip": "${alibabacloudstack_natgateway_forward_entries.default.ExternalIp}_fake",
		}),
	}

	external_portConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"external_port": "${alibabacloudstack_natgateway_forward_entries.default.ExternalPort}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"external_port": "${alibabacloudstack_natgateway_forward_entries.default.ExternalPort}_fake",
		}),
	}

	forward_entry_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"forward_entry_id": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"forward_entry_id": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryId}_fake",
		}),
	}

	forward_entry_nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"forward_entry_name": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryName}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"forward_entry_name": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryName}_fake",
		}),
	}

	forward_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"forward_table_id": "${alibabacloudstack_natgateway_forward_entries.default.ForwardTableId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"forward_table_id": "${alibabacloudstack_natgateway_forward_entries.default.ForwardTableId}_fake",
		}),
	}

	internal_ipConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"internal_ip": "${alibabacloudstack_natgateway_forward_entries.default.InternalIp}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"internal_ip": "${alibabacloudstack_natgateway_forward_entries.default.InternalIp}_fake",
		}),
	}

	internal_portConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"internal_port": "${alibabacloudstack_natgateway_forward_entries.default.InternalPort}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":           []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"internal_port": "${alibabacloudstack_natgateway_forward_entries.default.InternalPort}_fake",
		}),
	}

	ip_protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"ip_protocol": "${alibabacloudstack_natgateway_forward_entries.default.IpProtocol}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"ip_protocol": "${alibabacloudstack_natgateway_forward_entries.default.IpProtocol}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                []string{"${alibabacloudstack_natgateway_forward_entries.default.id}"},
			"external_ip":        "${alibabacloudstack_natgateway_forward_entries.default.ExternalIp}",
			"external_port":      "${alibabacloudstack_natgateway_forward_entries.default.ExternalPort}",
			"forward_entry_id":   "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryId}",
			"forward_entry_name": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryName}",
			"forward_table_id":   "${alibabacloudstack_natgateway_forward_entries.default.ForwardTableId}",
			"internal_ip":        "${alibabacloudstack_natgateway_forward_entries.default.InternalIp}",
			"internal_port":      "${alibabacloudstack_natgateway_forward_entries.default.InternalPort}",
			"ip_protocol":        "${alibabacloudstack_natgateway_forward_entries.default.IpProtocol}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                []string{"${alibabacloudstack_natgateway_forward_entries.default.id}_fake"},
			"external_ip":        "${alibabacloudstack_natgateway_forward_entries.default.ExternalIp}_fake",
			"external_port":      "${alibabacloudstack_natgateway_forward_entries.default.ExternalPort}_fake",
			"forward_entry_id":   "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryId}_fake",
			"forward_entry_name": "${alibabacloudstack_natgateway_forward_entries.default.ForwardEntryName}_fake",
			"forward_table_id":   "${alibabacloudstack_natgateway_forward_entries.default.ForwardTableId}_fake",
			"internal_ip":        "${alibabacloudstack_natgateway_forward_entries.default.InternalIp}_fake",
			"internal_port":      "${alibabacloudstack_natgateway_forward_entries.default.InternalPort}_fake",
			"ip_protocol":        "${alibabacloudstack_natgateway_forward_entries.default.IpProtocol}_fake",
		}),
	}

	AlibabacloudstackNatgatewayForwardEntriesCheckInfo.dataSourceTestCheck(t, rand, idsConf, external_ipConf, external_portConf, forward_entry_idConf, forward_entry_nameConf, forward_table_idConf, internal_ipConf, internal_portConf, ip_protocolConf, allConf)
}

var existAlibabacloudstackNatgatewayForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":    "1",
		"entries.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackNatgatewayForwardEntriesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var AlibabacloudstackNatgatewayForwardEntriesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_natgateway_forward_entries.default",
	existMapFunc: existAlibabacloudstackNatgatewayForwardEntriesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackNatgatewayForwardEntriesMapFunc,
}

func dataSourceNatgatewayForwardEntriesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
