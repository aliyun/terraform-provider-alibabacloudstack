package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmcpClustersDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_bmcp_clusters.default"
	password := getAccTestPassword(12)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		func(name string) string {
			return dataSourceBmcpClustersConfigDependence(name, password)
		})

	// Test cluster_name_regex filter using real cluster name from data source
	clusterNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_name_regex": "${alibabacloudstack_bmcp_cluster.default.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_name_regex": "cluster-fake-12345",
		}),
	}

	// Test status_regex filter using real status from data source
	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status_regex": "${alibabacloudstack_bmcp_cluster.default.status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status_regex": "status-fake-12345",
		}),
	}

	var existClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        CHECKSET,
			"ids.0":                        CHECKSET,
			"clusters.#":                   CHECKSET,
			"clusters.0.id":                CHECKSET,
			"clusters.0.cluster_id":        CHECKSET,
			"clusters.0.cluster_name":      CHECKSET,
			"clusters.0.status":            CHECKSET,
			"clusters.0.region_id":         CHECKSET,
			"clusters.0.zone":              CHECKSET,
			"clusters.0.vpc_id":            CHECKSET,
			"clusters.0.evpc_id":           CHECKSET,
			"clusters.0.organization":      CHECKSET,
			"clusters.0.node_count":        CHECKSET,
			"clusters.0.cpu_count":         CHECKSET,
			"clusters.0.mem_count":         CHECKSET,
			"clusters.0.flops_count":       CHECKSET,
			"clusters.0.video_memory":      CHECKSET,
			"clusters.0.cluster_arch_type": CHECKSET,
			"clusters.0.switch_method":     CHECKSET,
			"clusters.0.create_time":       CHECKSET,
			"clusters.0.update_time":       CHECKSET,
		}
	}

	var fakeClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var clustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existClustersMapFunc,
		fakeMapFunc:  fakeClustersMapFunc,
	}

	// Test all filter configurations
	clustersCheckInfo.dataSourceTestCheck(t, rand, clusterNameRegexConf, statusRegexConf)
}

func dataSourceBmcpClustersConfigDependence(name string, password string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name       = "${var.name}"
	cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	name              = "${var.name}"
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "192.168.40.0/24"
	is_cgw            = true
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_vswitch" "standard" {
	name              = "${var.name}-std"
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "192.168.50.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_evpc_evpc" "default" {
	evpc_name   = "${var.name}"
	description = "${var.name}"
}

data "alibabacloudstack_bmcp_machinetypes" "all" {
    min_standard_instance_count = 1
}

resource "alibabacloudstack_bmcp_cluster" "default" {
	cluster_name        = "${var.name}"
	vpc_id              = alibabacloudstack_vpc.default.id
	evpc_id             = alibabacloudstack_evpc_evpc.default.id
	zone_id             = data.alibabacloudstack_zones.default.zones.0.id
	standard_vswitch_id = alibabacloudstack_vswitch.standard.id
	password            = "%s"
	machine_type        = data.alibabacloudstack_bmcp_machinetypes.all.machinetypes.0.name
	node_count          = 1
	vswitch_id          = alibabacloudstack_vswitch.default.id

	lifecycle {
      ignore_changes = [
		machine_type
      ]
  	}
}
`, name, password)
}
