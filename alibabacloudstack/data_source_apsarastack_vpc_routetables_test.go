package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcRoutetablesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_route_tables.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccRouteTablesDatasource%d", rand),
		dataSourceVpcRoutetablesConfigDependence)

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_route_table.default.name}",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
			"ids":        []string{"${alibabacloudstack_route_table.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_route_table.default.name}_fake",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
			"ids":        []string{"${alibabacloudstack_route_table.default.id}"},
		}),
	}

	routeTablesCheckInfo.dataSourceTestCheck(t, rand, allConfig)
}

func dataSourceVpcRoutetablesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alibabacloudstack_route_table" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  name = "${var.name}"
  description = "${var.name}_description"
}
`, name)
}

var existRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                CHECKSET,
		"names.#":              CHECKSET,
		"tables.#":             CHECKSET,
		"tables.0.id":          CHECKSET,
		"tables.0.name":        fmt.Sprintf("tf-testAccRouteTablesDatasource%d", rand),
		"tables.0.description": fmt.Sprintf("tf-testAccRouteTablesDatasource%d_description", rand),
	}
}

var fakeRouteTablesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":    "0",
		"names.#":  "0",
		"tables.#": "0",
	}
}

var routeTablesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_route_tables.default",
	existMapFunc: existRouteTablesMapFunc,
	fakeMapFunc:  fakeRouteTablesMapFunc,
}
