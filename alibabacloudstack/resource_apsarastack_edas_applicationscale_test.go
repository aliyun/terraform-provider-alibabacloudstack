package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEdasInstanceApplicationAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "alibabacloudstack_edas_application_scale.default"

	ra := resourceAttrInit(resourceId, edasIAAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasiaattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasIAAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckIAAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alibabacloudstack_edas_application.default.id}",
					"deploy_group": "${data.alibabacloudstack_edas_deploy_groups.default.groups.0.group_id}",
					"ecu_info":     []string{"${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasIAAttachmentMap = map[string]string{
	"app_id":       CHECKSET,
	"deploy_group": CHECKSET,
}

func testEdasCheckIAAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasIAAttachmentDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "alibabacloudstack_zones" "default" {
			available_resource_creation= "VSwitch"
		}	
		resource "alibabacloudstack_vpc" "default" {
		  name = "${var.name}"
		  cidr_block = "10.1.0.0/21"
		}
		
		resource "alibabacloudstack_vswitch" "default" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "10.1.1.0/24"
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
		  //image_id = "centos_7_7_x64_20G_alibase_20200426.vhd"
			image_id="centos_7_7_x64_20G_alibase_20211028.vhd"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  system_disk_category = "cloud_efficiency"
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
		  region_id    = "cn-neimeng-env30-d01"
		}
		
		resource "alibabacloudstack_edas_instance_cluster_attachment" "default" {
		  cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
		  instance_ids = ["${alibabacloudstack_instance.default.id}"]
		  pass_word = "Li65272237###"
		}
	
		resource "alibabacloudstack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${alibabacloudstack_edas_cluster.default.id}"
		  package_type = "JAR"
		  ecu_info = ["${alibabacloudstack_edas_instance_cluster_attachment.default.ecu_map[alibabacloudstack_instance.default.id]}"]
		  build_pack_id = "15"
		}

		data "alibabacloudstack_edas_deploy_groups" "default" {
		  app_id = "${alibabacloudstack_edas_application.default.id}"
		}
		`, name)
}
