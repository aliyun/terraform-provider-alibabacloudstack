package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApsaraStackEdasSlbAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "apsarastack_edas_slb_attachment.default"

	ra := resourceAttrInit(resourceId, edasSLBAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasslbattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSLBAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckSLBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"app_id": "${apsarastack_edas_application.default.id}",
					"app_id": "22ba7083-ac60-4884-9d28-b0eaf71ac427",
					"slb_id": "${apsarastack_slb.default.id}",
					"slb_ip": "${apsarastack_slb.default.address}",
					"type":   "${apsarastack_slb.default.address_type}",
					"listener_port": "22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasSLBAttachmentMap = map[string]string{
	"app_id": CHECKSET,
	"slb_id": CHECKSET,
	"slb_ip": CHECKSET,
	"type":   CHECKSET,
}

func testEdasCheckSLBAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasSLBAttachmentDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}
		data "apsarastack_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		//data "apsarastack_instance_types" "default" {
		// cpu_core_count    = 1
		// memory_size       = 2
		//}
		//
		resource "apsarastack_vpc" "default" {
		cidr_block = "172.16.0.0/12"
		name       = "${var.name}"
		}
		
		resource "apsarastack_vswitch" "default" {
		vpc_id            = "${apsarastack_vpc.default.id}"
		cidr_block        = "172.16.0.0/24"
		availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		name = "${var.name}"
		}
		//
		//resource "apsarastack_security_group" "default" {
		// name = "${var.name}"
		// description = "New security group"
		// vpc_id = "${apsarastack_vpc.default.id}"
		//}
		//
		//resource "apsarastack_instance" "default" {
		// vswitch_id = "${apsarastack_vswitch.default.id}"
		// image_id = "centos_7_7_x64_20G_alibase_20211028.vhd"
		// availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		// system_disk_category = "cloud_ssd"
		// system_disk_size ="60"
		// instance_type = "ecs.xn4.small"
		//
		// security_groups = ["${apsarastack_security_group.default.id}"]
		// instance_name = "${var.name}"
		// tags = {
		//	Name = "TerraformTest-instance"
		// }
		//}
		//
		//resource "apsarastack_edas_cluster" "default" {
		// cluster_name = "${var.name}"
		// cluster_type = 2
		// network_mode = 2
		// vpc_id       = "${apsarastack_vpc.default.id}"
		//}
		//
		//resource "apsarastack_edas_instance_cluster_attachment" "default" {
		// cluster_id = "${apsarastack_edas_cluster.default.id}"
		// instance_ids = ["${apsarastack_instance.default.id}"]
		// pass_word = "Li65272237###"
		//}
		//
		//resource "apsarastack_edas_application" "default" {
		// application_name = "${var.name}"
		// cluster_id = "${apsarastack_edas_cluster.default.id}"
		// package_type = "JAR"
		// //ecu_info = ["${apsarastack_edas_instance_cluster_attachment.default.ecu_map[apsarastack_instance.default.id]}"]
		// ecu_info = ["${apsarastack_edas_instance_cluster_attachment.default.ecu_map[apsarastack_instance.default.id]}"]
		//}

		resource "apsarastack_slb" "default" {
		  name          = "${var.name}"
		  vswitch_id    = "${apsarastack_vswitch.default.id}"
          address_type  = "internet"
		  specification = "slb.s1.small"
		}
		`, name)
}
