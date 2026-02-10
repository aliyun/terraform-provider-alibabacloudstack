package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCloudfirewallControlpoliciesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_cloud_firewall_control_policies.default"
	name := fmt.Sprintf("tf-cfwpolicies%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudFirewallControlPoliciesConfigDependence)

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"description": "${alibabacloudstack_cloud_firewall_control_policy.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"description": "fake_description",
		}),
	}

	aclUuidConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"acl_uuid": "${alibabacloudstack_cloud_firewall_control_policy.default.acl_uuid}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"acl_uuid": "fake-acl-uuid",
		}),
	}

	sourceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"source": "${alibabacloudstack_cloud_firewall_control_policy.default.source}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"source": "1.1.1.1/32",
		}),
	}

	destinationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"destination": "${alibabacloudstack_cloud_firewall_control_policy.default.destination}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"destination": "2.2.2.2/32",
		}),
	}

	combinedConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"acl_action": "${alibabacloudstack_cloud_firewall_control_policy.default.acl_action}",
			"proto":      "${alibabacloudstack_cloud_firewall_control_policy.default.proto}",
			"description": "${alibabacloudstack_cloud_firewall_control_policy.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"direction" : "${alibabacloudstack_cloud_firewall_control_policy.default.direction}",
			"acl_action": "drop",
			"proto":      "TCP",
			"description": "non-existent",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"policies.#":                  "1",
			"policies.0.description":      CHECKSET,
			"policies.0.application_name": CHECKSET,
			"policies.0.acl_action":       CHECKSET,
			"policies.0.destination_type": CHECKSET,
			"policies.0.destination":      CHECKSET,
			"policies.0.direction":        CHECKSET,
			"policies.0.proto":            CHECKSET,
			"policies.0.source":           CHECKSET,
			"policies.0.source_type":      CHECKSET,
			"policies.0.acl_uuid":         CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}

	testAttr := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
		Providers:    testYunDunProviders(),
	}

	testAttr.dataSourceTestCheck(t, rand,
		descriptionConf,
		aclUuidConf,
		sourceConf,
		destinationConf,
		combinedConf,
	)
}

func dataSourceCloudFirewallControlPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}	

resource "alibabacloudstack_cloudfw_address_book" "default" {
  group_type    = "ip"
  group_name    = var.name
  address_list  = ["100.100.100.100/30"]
  description   = "test address book"
}

resource "alibabacloudstack_cloudfw_address_book" "port" {
  group_type    = "port"
  group_name    = var.name
  address_list  = ["8888", "9999"]
  description   = "test port book"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
  application_name = "ANY"
  acl_action       = "accept"
  description      = var.name
  destination_type = "net"
  destination      = "114.2.3.0/24"
  direction        = "in"
  proto            = "ANY"
  source           = "192.1.1.0/24"
  source_type      = "net"
  dest_port        = "8080/8080"
  dest_port_type   = "port"
  release          = true
}

`, name  )
}
