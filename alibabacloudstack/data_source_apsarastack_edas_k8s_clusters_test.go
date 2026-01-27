package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEdasK8sClustersDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_edas_k8s_clusters.default"
	name := fmt.Sprintf("testtf%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasK8sClustersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${local.k8s_cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_k8s_cluster.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_k8s_cluster.default.id}_fake"},
		}),
	}
	csidConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cs_clsuter_id": "${alibabacloudstack_edas_k8s_cluster.default.cs_cluster_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cs_clsuter_id": "${alibabacloudstack_edas_k8s_cluster.default.cs_cluster_id}_fake",
		}),
	}

	logicalRegionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_k8s_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_k8s_cluster.default.namespace_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_k8s_cluster.default.id}_fake"},
			"logical_region_id": "${alibabacloudstack_edas_k8s_cluster.default.namespace_id}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_k8s_cluster.default.id}"},
			"logical_region_id": "${alibabacloudstack_edas_k8s_cluster.default.namespace_id}",
			"name_regex":        "${local.k8s_cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alibabacloudstack_edas_k8s_cluster.default.id}_fake"},
			"logical_region_id": "${alibabacloudstack_edas_k8s_cluster.default.namespace_id}_fake",
			"name_regex":        "${local.k8s_cluster_name}_fake",
		}),
	}

	var existEdasK8sClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":                   "1",
			"clusters.0.cluster_id":        CHECKSET,
			"clusters.0.cluster_name":      CHECKSET,
			"clusters.0.cluster_type":      CHECKSET,
			"clusters.0.network_mode":      CHECKSET,
			"clusters.0.vpc_id":            CHECKSET,
			"clusters.0.logical_region_id": CHECKSET,
		}
	}

	var fakeEdasK8sClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var edasApplicationCheckInfo = dataSourceAttr{
		resourceId:        resourceId,
		existMapFunc:      existEdasK8sClustersMapFunc,
		fakeMapFunc:       fakeEdasK8sClustersMapFunc,
		ExternalProviders: testAccExternalProviders,
	}

	edasApplicationCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, csidConf, logicalRegionConf, allConf)
}

func dataSourceEdasK8sClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_account" "current" {
}

locals {
	logical_id = "${data.alibabacloudstack_account.current.region}:${var.name}"
}

%s

resource "alibabacloudstack_edas_namespace" "default" {
	description =      "${var.name}"
	namespace_name =       "${var.name}"
	namespace_logical_id = substr(local.logical_id, 0, min(length(local.logical_id), 32))
}

resource "alibabacloudstack_edas_k8s_cluster" "default" {
    cs_cluster_id = "${local.k8s_cluster_id}"
	namespace_id  = "${alibabacloudstack_edas_namespace.default.namespace_logical_id}"
}

`, name, AckK8sCommonTestCase())
}
