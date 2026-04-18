package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestUatAlibabacloudStackBmcpNodesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_bmcp_nodes.default"

	password := getAccTestPassword(12)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		func(name string) string { return dataSourceBmcpNodesConfigDependence(name, password) })

	// Test node_name_regex filter using real data
	nodeNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"node_name_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_name != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_name, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_name))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"node_name_regex": "node-fake-12345",
		}),
	}

	// Test node_id_regex filter using real data
	nodeIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"node_id_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_id != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_id, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.node_id))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"node_id_regex": "node-fake-12345",
		}),
	}

	// Test cluster_id_regex filter using real data
	clusterIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_id_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_id != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_id, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_id))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_id_regex": "cluster-fake-12345",
		}),
	}

	// Test cluster_name_regex filter using real data
	clusterNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_name_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_name != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_name, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.cluster_name))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_name_regex": "cluster-fake-12345",
		}),
	}

	// Test sn_regex filter using real data
	snRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sn_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.sn != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.sn, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.sn))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"sn_regex": "SN-fake-12345",
		}),
	}

	// Test vpc_ip_regex filter using real data
	vpcIpRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_ip_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.vpc_ip != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.vpc_ip, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.vpc_ip))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_ip_regex": "192.168.1.1",
		}),
	}

	var existNodesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     CHECKSET,
			"ids.0":                     CHECKSET,
			"nodes.#":                   CHECKSET,
			"nodes.0.id":                CHECKSET,
			"nodes.0.node_id":           CHECKSET,
			"nodes.0.node_name":         CHECKSET,
			"nodes.0.cluster_id":        CHECKSET,
			"nodes.0.cluster_name":      CHECKSET,
			"nodes.0.sn":                CHECKSET,
			"nodes.0.vpc_ip":            CHECKSET,
			"nodes.0.status":            CHECKSET,
			"nodes.0.machine_type":      CHECKSET,
			"nodes.0.machine_type_name": CHECKSET,
			"nodes.0.cpu_arch":          CHECKSET,
			"nodes.0.cpu_number":        CHECKSET,
			"nodes.0.memory":            CHECKSET,
			"nodes.0.disk":              CHECKSET,
			"nodes.0.gpu_num":           CHECKSET,
			"nodes.0.gpu_model":         CHECKSET,
			"nodes.0.region_id":         CHECKSET,
			"nodes.0.create_time":       CHECKSET,
		}
	}

	var fakeNodesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"nodes.#": "0",
		}
	}

	var nodesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existNodesMapFunc,
		fakeMapFunc:  fakeNodesMapFunc,
	}

	// Test all filter configurations
	nodesCheckInfo.dataSourceTestCheck(t, rand, nodeNameRegexConf, nodeIdRegexConf, clusterIdRegexConf, clusterNameRegexConf, snRegexConf, vpcIpRegexConf)
}

func dataSourceBmcpNodesConfigDependence(name string, password string) string {
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

data "alibabacloudstack_bmcp_nodes" "all" {
	depends_on = [alibabacloudstack_bmcp_cluster.default]
}

`, name, password)
}
