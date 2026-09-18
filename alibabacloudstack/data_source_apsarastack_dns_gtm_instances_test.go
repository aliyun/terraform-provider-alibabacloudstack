package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsGtmInstancesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_gtm_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourceDnsGtmInstancesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                 "1",
				"instances.#":           "1",
				"instances.0.name":      fmt.Sprintf("tfacc%d", rand),
				"instances.0.prefix":    fmt.Sprintf("tfacc%d", rand),
				"instances.0.ttl":       "300",
				"instances.0.zone_name": fmt.Sprintf("tfacc%d.testtf.", rand),
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":       "0",
				"instances.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_gtm_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^fake-name$",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_gtm_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_gtm_instance.default.name}",
			"ids":        []string{"${alibabacloudstack_dns_gtm_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_gtm_instance.default.name}",
			"ids":        []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceDnsGtmInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_vpc_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
	name   = "${var.name}.testtf."
	remark = "Created by Terraform for DNS GTM instance test"
	vpc_ids = [
		"${alibabacloudstack_vpc_vpc.default.id}"
	]
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
	name    = "${var.name}"
	prefix  = "${var.name}"
	zone_id = "${alibabacloudstack_dns_private_domain.default.id}"
	ttl     = 300
}
`, name)
}
