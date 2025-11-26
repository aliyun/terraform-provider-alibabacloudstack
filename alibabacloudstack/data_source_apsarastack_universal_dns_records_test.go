package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackUniversalDnsRecordsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_universal_dns_records.default",
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
			existConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"name_regex": `"tf-testacc[0-9]+"`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"name_regex": `"nonexistent-record"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"ids": `["${alibabacloudstack_universal_dns_record.default.id}"]`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"ids": `["nonexistent-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"ids":        `["${alibabacloudstack_universal_dns_record.default.id}"]`,
				"name_regex": `"tf-testacc[0-9]+"`,
			}),
			fakeConfig: AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand, map[string]string{
				"ids":        `["nonexistent-id"]`,
				"name_regex": `"nonexistent-record"`,
			}),
		},
	)
}

// Dependency template generation method
func AlibabacloudTestAccUniversalDnsRecorddependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
    default = "tf-testacc%d"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS record test"
}

resource "alibabacloudstack_universal_dns_record" "default" {
    zone_id      = "${alibabacloudstack_universal_dns_domain.default.id}"
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

data "alibabacloudstack_universal_dns_records" "default" {
	zone_id = "${alibabacloudstack_universal_dns_record.default.zone_id}"
    %s
}
`, rand, strings.Join(pairs, "\n    "))
}
