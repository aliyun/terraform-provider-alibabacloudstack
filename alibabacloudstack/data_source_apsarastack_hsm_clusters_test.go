package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackHsmClustersDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_hsvccm_clusters.default"
	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                   "1",
				"clusters.#":              "1",
				"clusters.0.cluster_name": fmt.Sprintf("tf_hsm_cluster%d", rand),
				"clusters.0.vsm_type":     "gvsm",
				"clusters.0.cluster_size": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":      "0",
				"clusters.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_hsm_cluster.default.cluster_name}"`,
		}),
		fakeConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"name_regex": `"fake-name-regex"`,
		}),
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_hsm_cluster.default.id}"]`,
		}),
		fakeConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"ids": `"fake-instance-id"`,
		}),
	}

	vsmTypeConf := dataSourceTestAccConfig{
		existConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"vsm_type": `"${alibabacloudstack_hsm_cluster.default.vsm_type}"`,
		}),
		fakeConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"vsm_type": `"evsm"`,
		}),
	}

	zoneNoConf := dataSourceTestAccConfig{
		existConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"zone_no": `"${alibabacloudstack_hsm_cluster.default.zone_no}"`,
		}),
		fakeConfig: buildHsmClusterDependenciesNew(rand, map[string]string{
			"zone_no": `"fake-zone-no"`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, instanceIdConf, vsmTypeConf, zoneNoConf)
}

// buildHsmClusterDependenciesNew generates the necessary dependencies for testing
func buildHsmClusterDependenciesNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf_hsm_cluster%d"
}

data "alibabacloudstack_zones" "default" {
  provider = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
}

resource "alibabacloudstack_hsm_instance" "default" {
	product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}"
	vsm_type = "gvsm"
	zone_no = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark = "${var.name}"
}

resource "alibabacloudstack_hsm_cluster" "default" {
  cluster_name     = "${var.name}"
  master_instance_id = alibabacloudstack_hsm_instance.default.id
  vpc_id           = alibabacloudstack_vpc.default.id
  vswitch_ids      = [alibabacloudstack_vswitch.default.id]
  zone_nos         = [data.alibabacloudstack_zones.default.zones.0.id]
  ip_white_list    = ["123.12.13.1/16"]
}

data "alibabacloudstack_hsm_clusters" "default" {
 %s
}
`, rand, strings.Join(pairs, "\n   "))
}
