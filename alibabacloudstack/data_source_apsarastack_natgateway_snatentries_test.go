package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNatgatewaySnatentriesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_snat_entries.default"
	name := fmt.Sprintf("tf-testAccForSnatEntriesDatasource%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSnatEntriesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"ids":           []string{"${alibabacloudstack_snat_entry.default.snat_entry_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"ids":           []string{"fake-snat-entry-id"},
		}),
	}

	sourceCidrConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"source_cidr":   "${alibabacloudstack_vswitch.default.cidr_block}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"source_cidr":   "192.168.0.0/24",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"ids":           []string{"${alibabacloudstack_snat_entry.default.snat_entry_id}"},
			"source_cidr":   "${alibabacloudstack_vswitch.default.cidr_block}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"snat_table_id": "${alibabacloudstack_nat_gateway.default.snat_table_ids}",
			"ids":           []string{"fake-snat-entry-id"},
			"source_cidr":   "192.168.0.0/24",
		}),
	}

	var existSnatEntriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "1",
			"entries.#":  "1",
			"entries.0.id":       CHECKSET,
			"entries.0.snat_ip":  CHECKSET,
			"entries.0.source_cidr": "172.16.0.0/21",
			"entries.0.status":   CHECKSET,
		}
	}

	var fakeSnatEntriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"entries.#": "0",
		}
	}

	var snatEntriesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSnatEntriesMapFunc,
		fakeMapFunc:  fakeSnatEntriesMapFunc,
	}
	snatEntriesCheckInfo.dataSourceTestCheck(t, rand, idsConf, sourceCidrConf, allConf)
}

func dataSourceSnatEntriesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = alibabacloudstack_vpc.default.id
  specification = "Small"
  name          = var.name
}

resource "alibabacloudstack_eip" "default" {
  name = var.name
}

resource "alibabacloudstack_eip_association" "default" {
  allocation_id = alibabacloudstack_eip.default.id
  instance_id   = alibabacloudstack_nat_gateway.default.id
}

resource "alibabacloudstack_snat_entry" "default" {
  snat_table_id     = alibabacloudstack_nat_gateway.default.snat_table_ids
  source_vswitch_id = alibabacloudstack_vswitch.default.id
  snat_ip           = alibabacloudstack_eip.default.ip_address
}

`, name, DataZoneCommonTestCase)
}
