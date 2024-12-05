package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackRouteTablesDataSourceBasic(t *testing.T) {
	preCheck := func() {
		testAccPreCheck(t)
	}
	rand := acctest.RandInt()

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackRouteTablesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_route_table.default.name}"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
			"ids":        `[ "${alibabacloudstack_route_table.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackRouteTablesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_route_table.default.name}_fake"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
			"ids":        `[ "${alibabacloudstack_route_table.default.id}" ]`,
		}),
	}

	routeTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConfig)
}

func testAccCheckAlibabacloudStackRouteTablesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccRouteTablesDatasource%d"
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

data "alibabacloudstack_route_tables" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
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
