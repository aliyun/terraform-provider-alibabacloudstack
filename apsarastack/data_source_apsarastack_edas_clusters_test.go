package apsarastack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackEdasClustersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.apsarastack_edas_clusters.default"
	name := fmt.Sprintf("tf-testacc-edas-clusters%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasClustersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${apsarastack_edas_cluster.default.cluster_name}",
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "fake_tf-testacc*",
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${apsarastack_edas_cluster.default.id}"},
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${apsarastack_edas_cluster.default.id}_fake"},
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${apsarastack_edas_cluster.default.id}"},
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
			"name_regex":        "${apsarastack_edas_cluster.default.cluster_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${apsarastack_edas_cluster.default.id}_fake"},
			"logical_region_id": os.Getenv("APSARASTACK_REGION"),
			"name_regex":        "${apsarastack_edas_cluster.default.cluster_name}",
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

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceEdasClustersConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "apsarastack_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "apsarastack_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${apsarastack_vpc.default.id}"
          region_id    = "cn-neimeng-env30-d01"
		}
		`, name)
}
