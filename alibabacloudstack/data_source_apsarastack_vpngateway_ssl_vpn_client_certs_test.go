package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpngatewaySslVpnClientCertsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_vpngateway_sslvpnclientcerts.default"
	name := fmt.Sprintf("tf_testAccVpngatewaySslvpnclientcertDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpngatewaySslvpnclientcertConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_vpngateway_sslvpnclientcert.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{

			"ids": []string{"${alibabacloudstack_vpngateway_sslvpnclientcert.default.id}-fake"},
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
			"ids":        []string{"${alibabacloudstack_vpngateway_sslvpnclientcert.default.id}"},
			"name_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_vpngateway_sslvpnclientcert.default.id}"},
			"name_regex": name + "_fake",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"ids.0":                                           CHECKSET,
			"ssl_vpn_client_certs.#":                          "1",
			"ssl_vpn_client_certs.0.ca_cert":                  CHECKSET,
			"ssl_vpn_client_certs.0.client_cert":              CHECKSET,
			"ssl_vpn_client_certs.0.client_config":            CHECKSET,
			"ssl_vpn_client_certs.0.client_key":               CHECKSET,
			"ssl_vpn_client_certs.0.create_time":              CHECKSET,
			"ssl_vpn_client_certs.0.end_time":                 CHECKSET,
			"ssl_vpn_client_certs.0.ssl_vpn_client_cert_id":   CHECKSET,
			"ssl_vpn_client_certs.0.ssl_vpn_client_cert_name": CHECKSET,
			"ssl_vpn_client_certs.0.ssl_vpn_server_id":        CHECKSET,
			"ssl_vpn_client_certs.0.status":                   CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"ssl_vpn_client_certs.#": "0",
		}
	}

	var ecsDedicatedHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	ecsDedicatedHostsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceVpngatewaySslvpnclientcertConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	%s
	
	resource "alibabacloudstack_vpngateway_ssl_vpnserver" "default" {
		client_ip_pool =  "10.8.0.0/24"
	
		local_subnet =  "192.168.1.0/24"
	
		ssl_vpn_server_name =  "${var.name}"
		vpn_gateway_id =     "${alibabacloudstack_vpn_gateway.default.id}"
		}

resource "alibabacloudstack_vpngateway_sslvpnclientcert" "default" {
	ssl_vpn_client_cert_name= "${var.name}"

	ssl_vpn_server_id= "${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}"

	vpn_gateway_id= "${alibabacloudstack_vpn_gateway.default.id}"
}

`, name, VpnGatewayCommonTestCase)
}
