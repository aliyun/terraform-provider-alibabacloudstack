package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestUatAlibabacloudStackBmcpNodesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_bmcp_nodes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceBmcpNodesConfigDependence)

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

	// Test out_of_band_ip_regex filter using real data
	outOfBandIpRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"out_of_band_ip_regex": `${data.alibabacloudstack_bmcp_nodes.all.nodes.0.out_of_band_ip != "" ? substr(data.alibabacloudstack_bmcp_nodes.all.nodes.0.out_of_band_ip, 0, min(3, length(data.alibabacloudstack_bmcp_nodes.all.nodes.0.out_of_band_ip))) : ""}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"out_of_band_ip_regex": "192.168.1.1",
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
			"nodes.0.out_of_band_ip":    CHECKSET,
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
	nodesCheckInfo.dataSourceTestCheck(t, rand, nodeNameRegexConf, nodeIdRegexConf, clusterIdRegexConf, clusterNameRegexConf, snRegexConf, vpcIpRegexConf, outOfBandIpRegexConf)
}

func dataSourceBmcpNodesConfigDependence(name string) string {
	return `
data "alibabacloudstack_bmcp_nodes" "all" {}
`
}
