package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsRecordsDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_dns_records.default"
	name := fmt.Sprintf("tf-testdnsrecordbasic-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDnsRecordsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_dns_domain.default.domain_id}",
			"ids":     []string{"${alibabacloudstack_dns_record.default.record_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_dns_domain.default.domain_id}",
			"ids":     []string{"fake-record-id"},
		}),
	}

	hostRecordRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":           "${alibabacloudstack_dns_domain.default.domain_id}",
			"host_record_regex": name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":           "${alibabacloudstack_dns_domain.default.domain_id}",
			"host_record_regex": "fake-regex",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_dns_domain.default.domain_id}",
			"type":    "A",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${alibabacloudstack_dns_domain.default.domain_id}",
			"type":    "CNAME",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":           "${alibabacloudstack_dns_domain.default.domain_id}",
			"ids":               []string{"${alibabacloudstack_dns_record.default.record_id}"},
			"host_record_regex": name,
			"type":              "A",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":           "${alibabacloudstack_dns_domain.default.domain_id}",
			"ids":               []string{"fake-record-id"},
			"host_record_regex": "fake-regex",
			"type":              "A",
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"records.#":            "1",
			"ids.#":                "1",
			"records.0.record_id":   CHECKSET,
			"records.0.zone_id":     CHECKSET,
			"records.0.name":       name,
			"records.0.type":       "A",
			"records.0.ttl":        "300",
			"records.0.rr_set.#":   "3",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"records.#": "0",
			"ids.#":     "0",
		}
	}

	var dnsRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}
	dnsRecordsCheckInfo.dataSourceTestCheck(t, rand, idsConf, hostRecordRegexConf, typeConf, allConf)
}

func dataSourceDnsRecordsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "${var.name}."
}

resource "alibabacloudstack_dns_record" "default" {
  zone_id      = alibabacloudstack_dns_domain.default.domain_id
  lba_strategy = "ALL_RR"
  name         = var.name
  type         = "A"
  ttl          = 300
  rr_set       = ["192.168.2.4", "192.168.2.7", "10.0.0.4"]
}

`, name)
}
