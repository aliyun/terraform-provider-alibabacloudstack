package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEdasClustersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_edas_clusters.default"
	name := fmt.Sprintf("tf%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasClustersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_edas_cluster.default.cluster_name}",
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "fake_tf-testacc*",
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}_fake"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
		}),
	}

	logicalRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}",
			"name_regex":        "${alibabacloudstack_edas_cluster.default.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_cluster.default.id}_fake"},
			"logical_region_id": "${alibabacloudstack_edas_cluster.default.logical_region_id}_fake",
			"name_regex":        "${alibabacloudstack_edas_cluster.default.cluster_name}_fake",
		}),
	}

	var existEdasClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":              "1",
			"clusters.0.cluster_id":   CHECKSET,
			"clusters.0.cluster_name": name,
			"clusters.0.cluster_type": "2",
			"clusters.0.network_mode": "2",
			"clusters.0.vpc_id":       CHECKSET,
			"clusters.0.region_id":    CHECKSET,
		}
	}

	var fakeEdasClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var edasApplicationCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasClustersMapFunc,
		fakeMapFunc:  fakeEdasClustersMapFunc,
	}

	edasApplicationCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, logicalRegionIdConf, allConf)
}

func dataSourceEdasClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

%s

variable "logical_id" {
  default = "%s:%s"
}

resource "alibabacloudstack_edas_namespace" "default" {
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  logical_region_id = "${alibabacloudstack_edas_namespace.default.namespace_logical_id}"
  network_mode = "2"
  cluster_type = "2"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

`, name, VpcCommonTestCase, defaultRegionToTest, name)
}
