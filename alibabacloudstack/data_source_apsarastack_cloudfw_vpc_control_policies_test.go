package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCloudfirewallVpcControlPolicies(t *testing.T) {
	resourceId := "data.alibabacloudstack_cloudfw_vpc_control_policies.default"
	rand := getAccTestRandInt(10000, 20000)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-vpccontrolpolicies-%d", rand),
		dataSourceCloudfwVpcControlPoliciesConfigDependence)

	testDataSourceAttr := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"policies.#":                  "1",
				"policies.0.description":      fmt.Sprintf("tf-testacc-vpccontrolpolicies-%d", rand),
				"policies.0.application_name": "ANY",
				"policies.0.source_type":      "net",
				"policies.0.dest_port":        "33/33",
				"policies.0.acl_action":       "log",
				"policies.0.proto":            "UDP",
				"policies.0.source":           "0.0.0.0/16",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"policies.#": "0",
			}
		},
		Providers: testYunDunProviders(),
	}

	testDataSourceAttr.dataSourceTestCheck(t, rand,
		// Test query by ids
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_cloudfw_vpc_control_policy.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_cloudfw_vpc_control_policy.default.id}_fake"},
			}),
		},
		// Test query by name_regex
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_cloudfw_vpc_control_policy.default.description}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_cloudfw_vpc_control_policy.default.description}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"proto": "${alibabacloudstack_cloudfw_vpc_control_policy.default.proto}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"proto": "TCP",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"acl_action": "${alibabacloudstack_cloudfw_vpc_control_policy.default.acl_action}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"acl_action": "accept",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"destination": "${alibabacloudstack_cloudfw_vpc_control_policy.default.destination}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"destination": "0.1.0.0/16",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"source": "${alibabacloudstack_cloudfw_vpc_control_policy.default.source}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"source": "0.1.0.0/16",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_cloudfw_vpc_control_policy.default.id}"},
				"name_regex":  "${alibabacloudstack_cloudfw_vpc_control_policy.default.description}",
				"proto":       "${alibabacloudstack_cloudfw_vpc_control_policy.default.proto}",
				"acl_action":  "${alibabacloudstack_cloudfw_vpc_control_policy.default.acl_action}",
				"destination": "${alibabacloudstack_cloudfw_vpc_control_policy.default.destination}",
				"source":      "${alibabacloudstack_cloudfw_vpc_control_policy.default.source}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":         []string{"${alibabacloudstack_cloudfw_vpc_control_policy.default.id}_fake"},
				"name_regex":  "${alibabacloudstack_cloudfw_vpc_control_policy.default.description}_fake",
				"proto":       "TCP",
				"acl_action":  "accept",
				"destination": "0.1.0.0/16",
				"source":      "0.1.0.0/16",
			}),
		},
	)
}

func dataSourceCloudfwVpcControlPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_cloudfw_vpc_control_policy" "default" {
    destination = "0.0.0.0/16"
	description = "${var.name}"
	application_name = "ANY"
	source_type = "net"
	dest_port = "33/33"
	acl_action = "log"
	destination_type = "net"

	source = "0.0.0.0/16"
	dest_port_type = "port"
	proto = "UDP"
	application_id = "0"
	release = true
}
`, name)
}
