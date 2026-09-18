package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcRouteTableAttachmentsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_vpc_route_table_attachments.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAlibabacloudstackVpcRouteTableAttachments%d", rand),
		dataSourceVpcRouteTableAttachmentsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"},
		}),
	}

	route_table_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}"},
			"route_table_id": "${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"},
			"route_table_id": "${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}_fake",
		}),
	}

	vswitch_idConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}"},
			"vswitch_id": "${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"},
			"vswitch_id": "${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}"},
			"route_table_id": "${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}",
			"vswitch_id":     "${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_vpc_route_table_attachments.default.id}_fake"},
			"route_table_id": "${alibabacloudstack_vpc_route_table_attachments.default.RouteTableId}_fake",
			"vswitch_id":     "${alibabacloudstack_vpc_route_table_attachments.default.VSwitchId}_fake",
		}),
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

func dataSourceVpcRouteTableAttachmentsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
