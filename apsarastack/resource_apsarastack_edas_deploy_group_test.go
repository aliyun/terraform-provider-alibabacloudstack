package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackEdasDeployGroup_basic(t *testing.T) {
	var v *edas.DeployGroup
	resourceId := "apsarastack_edas_deploy_group.default"

	ra := resourceAttrInit(resourceId, edasDeployGroupBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasdeploygroupbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasDeployGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasDeployGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":     "${apsarastack_edas_application.default.id}",
					"group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": fmt.Sprintf("tf-testacc-edasdeploygroupchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": fmt.Sprintf("tf-testacc-edasdeploygroupchange%v", rand)}),
				),
			},
		},
	})
}

func TestAccApsaraStackEdasDeployGroup_multi(t *testing.T) {
	var v *edas.DeployGroup
	resourceId := "apsarastack_edas_deploy_group.default.1"

	ra := resourceAttrInit(resourceId, edasDeployGroupBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}

	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasdeploygroupmulti%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasDeployGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasDeployGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":      "2",
					"app_id":     "${apsarastack_edas_application.default.id}",
					"group_name": "${var.name}-${count.index}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEdasDeployGroupDestroy(s *terraform.State) error {
	return nil
}

var edasDeployGroupBasicMap = map[string]string{
	"app_id":     CHECKSET,
	"group_name": CHECKSET,
	"group_type": CHECKSET,
}

func resourceEdasDeployGroupConfigDependence(name string) string {
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
          //region_id    = "cn-neimeng-env30-d01"
		}

		resource "apsarastack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${apsarastack_edas_cluster.default.id}"
		  package_type = "JAR"
		  build_pack_id = "15"
		}
		`, name)
}
