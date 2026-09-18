package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCspprivateHsmGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_cspprivate_hsm_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_hsm_group%v", rand),
		dataSourceCspprivateHsmGroupsConfigDependence)

	datasourceAttr := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":               "1",
				"groups.#":            "1",
				"groups.0.group_name": CHECKSET,
				"groups.0.status":     CHECKSET,
				"groups.0.hsm_count":  "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#": "0",
			}
		},
		Providers:         testYunDunProviders(),
		ExternalProviders: testAccExternalProviders,
	}

	// Test for ids filter
	testAccConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cspprivate_hsm_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cspprivate_hsm_group.default.id}_fake"},
		}),
	}
	datasourceAttr.dataSourceTestCheck(t, rand, testAccConf)
}

func dataSourceCspprivateHsmGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	provider = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_vpc" "vpc" {	
	provider = alibabacloudstack-common
	vpc_name = "${var.name}"
	cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
	provider = alibabacloudstack-common
	vpc_id = alibabacloudstack_vpc.vpc.id
	cidr_block = "192.168.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}


resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
	product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	alias_name = "${var.name}0"
	vpc_id = "${alibabacloudstack_vpc.vpc.id}"
	vpc_cidr_block = "${alibabacloudstack_vpc.vpc.cidr_block}"
	vswitch_id = "${alibabacloudstack_vswitch.vsw.id}"
	ip = "192.168.0.100"
}

%s

resource "alibabacloudstack_cspprivate_hsm_group" "default" {
	hsm_list = ["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]
	zone_ids = ["${data.alibabacloudstack_zones.default.zones.0.id}"]
	password = "${random_password.password.0.result}"
	hsm_count = 1
}
`, name, RandomPasswordTestCase(12, 1))
}
