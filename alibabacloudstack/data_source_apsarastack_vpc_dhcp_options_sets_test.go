package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcDhcpOptionsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_vpc_dhcp_options_sets.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sVpcDhcpOptionsDataSource-%d", defaultRegionToTest, rand),
		dataSourceVpcDhcpOptionsDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}-fakeTestAcccc",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_dhcp_options_set.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpc_dhcp_options_set.default.id}-fakeTestAcccc"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}",
			"ids":        []string{"${alibabacloudstack_vpc_dhcp_options_set.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}-fakeTestAcccc",
			"ids":        []string{"${alibabacloudstack_vpc_dhcp_options_set.default.id}-fakeTestAcccc"},
		}),
	}

	var existVpcDhcpOptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"ids.0":               CHECKSET,
			"dhcp_options_sets.#": "1",
			"dhcp_options_sets.0.dhcp_options_set_description": fmt.Sprintf("tf-testAcc%sVpcDhcpOptionsDataSource-%d", defaultRegionToTest, rand),
			"dhcp_options_sets.0.dhcp_options_set_name":        fmt.Sprintf("tf-testAcc%sVpcDhcpOptionsDataSource-%d", defaultRegionToTest, rand),
			"dhcp_options_sets.0.dhcp_options_set_id":          CHECKSET,
			"dhcp_options_sets.0.domain_name":                  CHECKSET,
			"dhcp_options_sets.0.domain_name_servers":          CHECKSET,
			"dhcp_options_sets.0.associate_vpcs.#":             "2",
		}
	}

	var fakeVpcDhcpOptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "0",
			"dhcp_options_sets.#": "0",
		}
	}

	var VpcDhcpOptionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcDhcpOptionsMapFunc,
		fakeMapFunc:  fakeVpcDhcpOptionsMapFunc,
	}

	VpcDhcpOptionsCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, allConf)
}

func dataSourceVpcDhcpOptionsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_dhcp_options_set" "default" {
  	dhcp_options_set_name =        "${var.name}"
	dhcp_options_set_description = "${var.name}"
	domain_name =                  "aliyun.com"
	domain_name_servers =         "10.82.0.12,10.82.0.13"
	associate_vpcs = [
		"${alibabacloudstack_vpc.default0.id}",
		"${alibabacloudstack_vpc.default1.id}"
	]
}

 `, name)
}
