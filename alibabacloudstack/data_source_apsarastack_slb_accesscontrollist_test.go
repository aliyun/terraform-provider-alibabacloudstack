package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackSlbAccessControlListsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_access_control_lists.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_access_control_lists.default.id}_fake"]`,
		}),
	}

	acl_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_slb_access_control_lists.default.id}"]`,
			"acl_name": `"${alibabacloudstack_slb_access_control_lists.default.AclName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":      `["${alibabacloudstack_slb_access_control_lists.default.id}_fake"]`,
			"acl_name": `"${alibabacloudstack_slb_access_control_lists.default.AclName}_fake"`,
		}),
	}

	address_ip_versionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_access_control_lists.default.id}"]`,
			"address_ip_version": `"${alibabacloudstack_slb_access_control_lists.default.AddressIpVersion}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_slb_access_control_lists.default.id}_fake"]`,
			"address_ip_version": `"${alibabacloudstack_slb_access_control_lists.default.AddressIpVersion}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_access_control_lists.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_slb_access_control_lists.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_slb_access_control_lists.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_slb_access_control_lists.default.ResourceGroupId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_access_control_lists.default.id}"]`,

			"acl_name":           `"${alibabacloudstack_slb_access_control_lists.default.AclName}"`,
			"address_ip_version": `"${alibabacloudstack_slb_access_control_lists.default.AddressIpVersion}"`,
			"resource_group_id":  `"${alibabacloudstack_slb_access_control_lists.default.ResourceGroupId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_slb_access_control_lists.default.id}_fake"]`,

			"acl_name":           `"${alibabacloudstack_slb_access_control_lists.default.AclName}_fake"`,
			"address_ip_version": `"${alibabacloudstack_slb_access_control_lists.default.AddressIpVersion}_fake"`,
			"resource_group_id":  `"${alibabacloudstack_slb_access_control_lists.default.ResourceGroupId}_fake"`}),
	}

	AlibabacloudstackSlbAccessControlListsCheckInfo.dataSourceTestCheck(t, rand, idsConf, acl_nameConf, address_ip_versionConf, resource_group_idConf, allConf)
}

var existAlibabacloudstackSlbAccessControlListsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"lists.#":    "1",
		"lists.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackSlbAccessControlListsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"lists.#": "0",
	}
}

var AlibabacloudstackSlbAccessControlListsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_slb_access_control_lists.default",
	existMapFunc: existAlibabacloudstackSlbAccessControlListsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackSlbAccessControlListsMapFunc,
}

func testAccCheckAlibabacloudstackSlbAccessControlListsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackSlbAccessControlLists%d"
}






data "alibabacloudstack_slb_access_control_lists" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
