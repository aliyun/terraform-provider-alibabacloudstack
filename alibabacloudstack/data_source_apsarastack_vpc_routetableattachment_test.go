package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackVpcRouteTableAttachmentsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_table_attachments.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"]`,
		}),
	}

	route_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_route_table_attachments.default.id}"]`,
			"route_table_id": `"${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"]`,
			"route_table_id": `"${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}_fake"`,
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_route_table_attachments.default.id}"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"]`,
			"vswitch_id": `"${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_table_attachments.default.id}"]`,

			"route_table_id": `"${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"]`,

			"route_table_id": `"${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}_fake"`,
			"vswitch_id":     `"${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}_fake"`}),
	}

	AlibabacloudstackVpcRouteTableAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, route_table_idConf, vswitch_idConf, allConf)
}

var existAlibabacloudstackVpcRouteTableAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#":    "1",
		"attachments.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcRouteTableAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var AlibabacloudstackVpcRouteTableAttachmentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_route_table_attachments.default",
	existMapFunc: existAlibabacloudstackVpcRouteTableAttachmentsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcRouteTableAttachmentsMapFunc,
}

func testAccCheckAlibabacloudstackVpcRouteTableAttachmentsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcRouteTableAttachments%d"
}






data "alibabacloudstack_vpc_route_table_attachments" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
