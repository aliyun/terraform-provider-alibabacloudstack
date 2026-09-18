package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsPrivateRecordsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_private_records.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc%d", rand),
		dataSourceDnsPrivateRecordsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"records.#":                "1",
				"records.0.name":           fmt.Sprintf("tf-testacc%d", rand),
				"records.0.type":           "A",
				"records.0.ttl":            "300",
				"records.0.lba_strategy":   "ALL_RR",
				"records.0.rdatas.#":       "2",
				"records.0.rdatas.0.value": "192.168.1.1",
				"records.0.rdatas.1.value": "127.0.0.1",
				"records.0.line_ids.#":     "1",
				"records.0.line_ids.0":     "default",
				"ids.#":                    "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"records.#": "0",
				"ids.#":     "0",
			}
		},
	}

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"zone_id":    "${alibabacloudstack_dns_private_record.default.zone_id}",
				"name_regex": "tf-testacc[0-9]+",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"zone_id":    "${alibabacloudstack_dns_private_record.default.zone_id}",
				"name_regex": "nonexistent-record",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"zone_id": "${alibabacloudstack_dns_private_record.default.zone_id}",
				"ids":     []string{"${alibabacloudstack_dns_private_record.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"zone_id": "${alibabacloudstack_dns_private_record.default.zone_id}",
				"ids":     []string{"nonexistent-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"zone_id":    "${alibabacloudstack_dns_private_record.default.zone_id}",
				"ids":        []string{"${alibabacloudstack_dns_private_record.default.id}"},
				"name_regex": "tf-testacc[0-9]+",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"zone_id":    "${alibabacloudstack_dns_private_record.default.zone_id}",
				"ids":        []string{"nonexistent-id"},
				"name_regex": "nonexistent-record",
			}),
		},
	)
}

func dataSourceDnsPrivateRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}

resource "alibabacloudstack_dns_private_record" "default" {
    zone_id      = "${alibabacloudstack_dns_private_domain.default.id}"
    name         = "${var.name}"
    type         = "A"
    ttl          = 300
    lba_strategy = "ALL_RR"
    line_ids     = ["default"]
    rdatas {
        value = "192.168.1.1"
    }
    rdatas {
        value = "127.0.0.1"
    }
}
`, name)
}
