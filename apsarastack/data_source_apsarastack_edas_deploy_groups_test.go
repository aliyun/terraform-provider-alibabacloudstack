package apsarastack

import (
	"fmt"
	"testing"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackEdasDeployGroupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.apsarastack_edas_deploy_groups.default"
	name := fmt.Sprintf("tf-testacc-edas-deploy-groups%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasDeployGroupConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${apsarastack_edas_deploy_group.default.group_name}",
			"app_id":     "${apsarastack_edas_application.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
			"app_id":     "${apsarastack_edas_application.default.id}",
		}),
	}
	var existEdasDeployGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":             "1",
			"groups.0.group_name":  name,
			"groups.0.app_id":      CHECKSET,
			"groups.0.group_type":  CHECKSET,
			"groups.0.cluster_id":  CHECKSET,
			"groups.0.create_time": CHECKSET,
			"groups.0.update_time": CHECKSET,
			"groups.0.group_id":    CHECKSET,
		}
	}

	var fakeEdasDeployGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
		}
	}

	var edasApplicationCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasDeployGroupsMapFunc,
		fakeMapFunc:  fakeEdasDeployGroupsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
	}

	edasApplicationCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceEdasDeployGroupConfigDependence(name string) string {
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

		resource "apsarastack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${apsarastack_edas_cluster.default.id}"
		  package_type = "JAR"
		  build_pack_id = "15"
		}
		
		resource "apsarastack_edas_deploy_group" "default" {
		  app_id = apsarastack_edas_application.default.id
		  group_name = "${var.name}"
		}		
		`, name)
}
