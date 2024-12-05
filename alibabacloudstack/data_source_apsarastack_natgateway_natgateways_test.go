package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackNatGatewaysDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_nat_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_nat_gateway.default.name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_nat_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_nat_gateway.default.id}_fake" ]`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alibabacloudstack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alibabacloudstack_vpc.default.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_nat_gateway.default.name}"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
			"ids":        `[ "${alibabacloudstack_nat_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_nat_gateway.default.name}"`,
			"ids":        `[ "${alibabacloudstack_nat_gateway.default.id}" ]`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}_fake"`,
		}),
	}

	natGatewaysCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, vpcIdConf, allConf)
}

func testAccCheckAlibabacloudStackNatGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNatGatewaysDatasource%d"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
    description = "${var.name}_decription"
}

data "alibabacloudstack_nat_gateways" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existNatGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                  "1",
		"ids.#":                       "1",
		"names.#":                     "1",
		"gateways.0.id":               CHECKSET,
		"gateways.0.spec":             "Small",
		"gateways.0.status":           "Available",
		"gateways.0.creation_time":    CHECKSET,
		"gateways.0.forward_table_id": CHECKSET,
		"gateways.0.snat_table_id":    CHECKSET,
		"gateways.0.name":             fmt.Sprintf("tf-testAccNatGatewaysDatasource%d", rand),
		"gateways.0.description":      fmt.Sprintf("tf-testAccNatGatewaysDatasource%d_decription", rand),
	}
}

var fakeNatGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var natGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nat_gateways.default",
	existMapFunc: existNatGatewaysMapFunc,
	fakeMapFunc:  fakeNatGatewaysMapFunc,
}
