package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectVbrPconnDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_expressconnect_vbr_pconn_associations.default"
	testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
	testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD")

	name := fmt.Sprintf("tf-testAcc-ExpressconnectVbrPconnDataSource-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressconnectVbrPconnDependence(rand))

	vbrIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vbr_id": "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vbr_id": "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}_",
		}),
	}

	physicalConnectionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vbr_id":                 "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"physical_connection_id": getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vbr_id":                 "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"physical_connection_id": getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"),
		}),
	}

	vlanIdIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vbr_id":  "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"vlan_id": fmt.Sprintf("%d", rand+3),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vbr_id":  "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"vlan_id": fmt.Sprintf("%d", rand+4),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vbr_id":                 "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"physical_connection_id": getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD"),
			"vlan_id":                fmt.Sprintf("%d", rand+3),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vbr_id":                 "${alibabacloudstack_expressconnect_vbr_pconn_association.default.vbr_id}",
			"physical_connection_id": getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"),
			"vlan_id":                fmt.Sprintf("%d", rand+4),
		}),
	}

	var existExpressconnectVbrPconnMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"ids.0":                            CHECKSET,
			"vbr_pconn_associations.#":         "1",
			"vbr_pconn_associations.0.vbr_id":  CHECKSET,
			"vbr_pconn_associations.0.vlan_id": CHECKSET,
			"vbr_pconn_associations.0.physical_connection_id":   CHECKSET,
			"vbr_pconn_associations.0.local_gateway_ip":         CHECKSET,
			"vbr_pconn_associations.0.peer_gateway_ip":          CHECKSET,
			"vbr_pconn_associations.0.peering_subnet_mask":      CHECKSET,
			"vbr_pconn_associations.0.local_ipv6_gateway_ip":    CHECKSET,
			"vbr_pconn_associations.0.peer_ipv6_gateway_ip":     CHECKSET,
			"vbr_pconn_associations.0.peering_ipv6_subnet_mask": CHECKSET,
		}
	}

	var fakeExpressconnectVbrPconnMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "0",
			"vbr_pconn_associations.#": "0",
		}
	}

	var ExpressconnectVbrPconnCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressconnectVbrPconnMapFunc,
		fakeMapFunc:  fakeExpressconnectVbrPconnMapFunc,
	}

	ExpressconnectVbrPconnCheckInfo.dataSourceTestCheck(t, rand, vbrIdConf, physicalConnectionIdConf, vlanIdIdConf, allConf)
}

func dataSourceExpressconnectVbrPconnDependence(vlanId int) func(name string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
		variable "name" {
		  default = "%s"
		}

		resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
			physical_connection_id =     "%s"
			vlan_id =                    %d
			local_gateway_ip =           "10.0.0.1"
			peer_gateway_ip =            "10.0.0.2"
			peering_subnet_mask =        "255.255.255.252"
			virtual_border_router_name = "${var.name}"
			enable_ipv6              = true
			local_ipv6_gateway_ip = "2408:4004:cc:400::1"
			peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
			peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
		}
		
		resource "alibabacloudstack_expressconnect_vbr_pconn_association" "default" {
			physical_connection_id=   "%s"
			vbr_id=                   "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
			vlan_id=                  "%d"
			local_gateway_ip=         "10.100.0.1"
			peer_gateway_ip=          "10.100.0.2"
			peering_subnet_mask=      "255.255.255.0"
			enable_ipv6              = true
			local_ipv6_gateway_ip=    "2408:4004:cc:500::1"
			peer_ipv6_gateway_ip=     "2408:4004:cc:500::2"
			peering_ipv6_subnet_mask= "2408:4004:cc:500::/56"
		}

		`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID_THIRD"), vlanId+3)
	}
}
