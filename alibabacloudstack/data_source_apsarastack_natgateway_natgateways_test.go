package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNatGatewaysDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_nat_gateways.default"
	name := fmt.Sprintf("tf-testacc-natgateways%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceNatGatewaysConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_nat_gateway.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nat_gateway.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_nat_gateway.default.id}_fake"},
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "${alibabacloudstack_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_id": "fake-vpc-id",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_nat_gateway.default.name}",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
			"ids":        []string{"${alibabacloudstack_nat_gateway.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_nat_gateway.default.name}",
			"vpc_id":     "fake-vpc-id",
			"ids":        []string{"${alibabacloudstack_nat_gateway.default.id}"},
		}),
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
			"gateways.0.name":             name,
			"gateways.0.description":      name+"_decription",
			"gateways.0.vpc_id":           CHECKSET,
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
		resourceId:   resourceId,
		existMapFunc: existNatGatewaysMapFunc,
		fakeMapFunc:  fakeNatGatewaysMapFunc,
	}
	natGatewaysCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, vpcIdConf, allConf)
}

func dataSourceNatGatewaysConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = alibabacloudstack_vpc.default.id
  specification = "Small"
  name          = var.name
  description   = "${var.name}_decription"
}

`, name, DataZoneCommonTestCase)
}
