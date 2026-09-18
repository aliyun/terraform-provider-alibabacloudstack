package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackHsmClustersDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_hsm_clusters.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_hsm_cluster%d", rand),
		dataSourceHsmClustersConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                          "1",
				"clusters.#":                     "1",
				"clusters.0.cluster_name":        fmt.Sprintf("tf_hsm_cluster%d", rand),
				"clusters.0.vsm_type":            "gvsm",
				"clusters.0.cluster_zones.#":     "1",
				"clusters.0.hsm_cluster_items.#": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":      "0",
				"clusters.#": "0",
			}
		},
		Providers:         testYunDunProviders(),
		ExternalProviders: testAccExternalProviders,
	}

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_hsm_cluster.default.cluster_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "fake-name-regex",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_hsm_cluster.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-instance-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"vsm_type": "${alibabacloudstack_hsm_instance.default.vsm_type}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"vsm_type": "evsm",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_hsm_cluster.default.cluster_name}",
				"ids":        []string{"${alibabacloudstack_hsm_cluster.default.id}"},
				"vsm_type":   "${alibabacloudstack_hsm_instance.default.vsm_type}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "fake-name-regex",
				"ids":        []string{"fake-instance-id"},
				"vsm_type":   "evsm",
			}),
		},
	)
}

func dataSourceHsmClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	provider                    = alibabacloudstack-common
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	provider   = alibabacloudstack-common
	vpc_name   = "${var.name}_vpc"
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	provider     = alibabacloudstack-common
	vswitch_name = "${var.name}_vsw"
	vpc_id       = alibabacloudstack_vpc.default.id
	cidr_block   = "172.16.1.0/24"
	zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
}

%s

resource "alibabacloudstack_hsm_instance" "default" {
	product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code}"
	vendor_code  = "${data.alibabacloudstack_hsm_vendors.default.vendors.0.code}"
	vsm_type     = "gvsm"
	zone_no      = "${data.alibabacloudstack_zones.default.zones.0.id}"
	remark       = "${var.name}"
}

resource "alibabacloudstack_hsm_cluster" "default" {
	cluster_name       = "${var.name}"
	master_instance_id = alibabacloudstack_hsm_instance.default.id
	vpc_id             = alibabacloudstack_vpc.default.id
	vswitch_ids        = alibabacloudstack_vswitch.default.id
	zone_nos           = data.alibabacloudstack_zones.default.zones.0.id
	ip_white_list      = "123.12.13.1/16"
	password           = random_password.password.0.result
}
`, name, RandomPasswordTestCase(12, 1))
}
