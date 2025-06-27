package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpngatewaySslVpnServersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_vpngateway_ssl_vpnservers.default"
	name := fmt.Sprintf("tf_testAccVpngatewaySslvpnserverDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpngatewaySslvpnserverConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}-fake"},
		}),
	}

	ascmcreateuserConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"ascm_create_user": "admin",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"ascm_create_user": "admin" + "fake",
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"name_regex": name + "fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"},
			"name_regex":       name,
			"ascm_create_user": "admin",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"},
			"name_regex":       name + "_fake",
			"ascm_create_user": "admin" + "_fake",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"ids.0":                                 CHECKSET,
			"ssl_vpn_servers.#":                     "1",
			"ssl_vpn_servers.0.ssl_vpn_server_name": CHECKSET,
			"ssl_vpn_servers.0.vpn_gateway_id":      CHECKSET,
			"ssl_vpn_servers.0.id":                  CHECKSET,
			"ssl_vpn_servers.0.connections":         CHECKSET,
			"ssl_vpn_servers.0.local_subnet":        CHECKSET,
			"ssl_vpn_servers.0.create_time":         CHECKSET,
			"ssl_vpn_servers.0.ascm_create_user":    CHECKSET,
			"ssl_vpn_servers.0.cipher":              CHECKSET,
			"ssl_vpn_servers.0.port":                CHECKSET,
			"ssl_vpn_servers.0.ssl_vpn_server_id":   CHECKSET,
			"ssl_vpn_servers.0.compress":            CHECKSET,
			"ssl_vpn_servers.0.proto":               CHECKSET,
			"ssl_vpn_servers.0.client_ip_pool":      CHECKSET,
			"ssl_vpn_servers.0.internet_ip":         CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "0",
			"ssl_vpn_servers.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, ascmcreateuserConf, nameRegexConf, allConf)
}

func dataSourceVpngatewaySslvpnserverConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_vpngateway_vpngateway" "default" {
description= "${var.name}"

vpn_gateway_name = "${var.name}"

bandwidth = "5"

vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"

vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"

enable_ssl = "true"

instance_charge_type = "PostPaid"
}

resource "alibabacloudstack_vpngateway_ssl_vpnserver" "default" {
ssl_vpn_server_name = "${var.name}"

vpn_gateway_id = "${alibabacloudstack_vpngateway_vpngateway.default.id}"

local_subnet = "192.168.1.0/24"

client_ip_pool = "10.8.0.0/24"
}


`, name, VSwitchCommonTestCase)
}
