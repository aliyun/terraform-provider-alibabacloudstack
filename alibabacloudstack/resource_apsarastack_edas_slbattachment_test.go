package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEdasSlbAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alibabacloudstack_edas_slb_attachment.default"

	ra := resourceAttrInit(resourceId, edasSLBAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasslbattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSLBAttachmentDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckSLBAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"app_id": "${alibabacloudstack_edas_application.default.id}",
					"app_id":        "22ba7083-ac60-4884-9d28-b0eaf71ac427",
					"slb_id":        "${alibabacloudstack_slb.default.id}",
					"slb_ip":        "${alibabacloudstack_slb.default.address}",
					"type":          "${alibabacloudstack_slb.default.address_type}",
					"listener_port": "22",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
		variable "password" {
		}
		data "alibabacloudstack_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		//data "alibabacloudstack_instance_types" "default" {
		// cpu_core_count    = 1
		// memory_size       = 2
		//}
		//
		resource "alibabacloudstack_vpc" "default" {
		cidr_block = "172.16.0.0/12"
		name       = "${var.name}"
		}
		
		resource "alibabacloudstack_vswitch" "default" {
		vpc_id            = "${alibabacloudstack_vpc.default.id}"
		cidr_block        = "172.16.0.0/24"
		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		name = "${var.name}"
		}
		
		resource "alibabacloudstack_security_group" "default" {
		name = "${var.name}"
		description = "New security group"
		vpc_id = "${alibabacloudstack_vpc.default.id}"
		}
		
		resource "alibabacloudstack_instance" "default" {
		vswitch_id = "${alibabacloudstack_vswitch.default.id}"
		image_id = "centos_7_7_x64_20G_alibase_20211028.vhd"
		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		system_disk_category = "cloud_ssd"
		system_disk_size ="60"
		instance_type = "ecs.xn4.small"
		
		security_groups = ["${alibabacloudstack_security_group.default.id}"]
		instance_name = "${var.name}"
		tags = {
			Name = "TerraformTest-instance"
		}
		}
		
		resource "alibabacloudstack_edas_cluster" "default" {
		cluster_name = "${var.name}"
		cluster_type = 2
		network_mode = 2
		vpc_id       = "${alibabacloudstack_vpc.default.id}"
		}
		
		resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
		cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
		instance_ids = ["${alibabacloudstack_instance.default.id}"]
		pass_word = var.password
		}
		
		resource "alibabacloudstack_edas_application" "default" {
		application_name = "${var.name}"
		cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
		package_type = "JAR"
		//ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
		ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
		}

		resource "alibabacloudstack_slb" "default" {
		  name          = "${var.name}"
		  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
          address_type  = "internet"
		  specification = "slb.s1.small"
		}
		`, name)
}
