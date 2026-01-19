package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcNetworkAclsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_network_acls.default"
	name := fmt.Sprintf("tf-testAccNetworkAcl-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackNetworkAclsDataSourceName)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_network_acl.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_network_acl.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_network_acl.default.network_acl_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_network_acl.default.network_acl_name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_network_acl.default.id}"},
			"status": "${alibabacloudstack_network_acl.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alibabacloudstack_network_acl.default.id}"},
			"status": "Modifying",
		}),
	}
	vpcConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "fake_vpc_id",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_network_acl.default.id}"},
			"name_regex":       "${alibabacloudstack_network_acl.default.network_acl_name}",
			"network_acl_name": "${alibabacloudstack_network_acl.default.network_acl_name}",
			"status":           "${alibabacloudstack_network_acl.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alibabacloudstack_network_acl.default.id}"},
			"name_regex":       "${alibabacloudstack_network_acl.default.network_acl_name}_fake",
			"network_acl_name": "${alibabacloudstack_network_acl.default.network_acl_name}_fake",
			"status":           "Modifying",
		}),
	}
	var existAlibabacloudStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"acls.#":                       "1",
			"acls.0.description":           fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.egress_acl_entries.#":  "1",
			"acls.0.ingress_acl_entries.#": "1",
			"acls.0.network_acl_name":      fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.vpc_id":                CHECKSET,
			"acls.0.status":                "Available",
		}
	}
	var fakeAlibabacloudStackNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackNetworkAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_network_acls.default",
		existMapFunc: existAlibabacloudStackNetworkAclsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackNetworkAclsDataSourceNameMapFunc,
	}
	alibabacloudstackNetworkAclsCheckInfo.dataSourceTestCheck(t, rand, idsConf, vpcConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlibabacloudStackNetworkAclsDataSourceName(name string) string {
	return fmt.Sprintf(`

variable "name" {	
	default = "%s"
}

%s

resource "alibabacloudstack_network_acl" "default" {
	description = "${var.name}"
	network_acl_name = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	resources {
		resource_id=   "${alibabacloudstack_vpc_vswitch.default.id}"
		resource_type= "VSwitch"
	}
}
`, name, VSwitchCommonTestCase)
}
