package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
)

func TestAccApsaraStackEdasApplicationPackageAttachment_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "apsarastack_application_deployment.default"
	ra := resourceAttrInit(resourceId, edasAPAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc-edasdeploymentbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasAPAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":   "${apsarastack_edas_application.default.id}",
					"group_id": "all",
					"war_url":  "http://edas-sz.oss-cn-shenzhen.aliyuncs.com/prod/demo/SPRING_CLOUD_CONSUMER.jar",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testEdasCheckDeploymentDestroy(s *terraform.State) error {
	return nil
}

var edasAPAttachmentBasicMap = map[string]string{
	"app_id":   CHECKSET,
	"group_id": CHECKSET,
	"war_url":  CHECKSET,
}

func resourceEdasAPAttachmentDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		data "apsarastack_zones" "default" {
			available_resource_creation= "VSwitch"
		}	
		resource "apsarastack_vpc" "default" {
		  name = "${var.name}"
		  cidr_block = "10.1.0.0/21"
		}
		
		resource "apsarastack_vswitch" "default" {
		  vpc_id = "${apsarastack_vpc.default.id}"
		  cidr_block = "10.1.1.0/24"
		 // availability_zone = "cn-neimeng-env30-amtest30001-a"
	      availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
		}
		
		resource "apsarastack_security_group" "default" {
		  name = "${var.name}"
		  description = "New security group"
		  vpc_id = "${apsarastack_vpc.default.id}"
		}
		
		resource "apsarastack_instance" "default" {
		  vswitch_id = "${apsarastack_vswitch.default.id}"
		  image_id = "centos_7_7_x64_20G_alibase_20211028.vhd"
		  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  system_disk_category = "cloud_ssd"
		  system_disk_size ="60"
		  instance_type = "ecs.xn4.small"
		
		  security_groups = ["${apsarastack_security_group.default.id}"]
		  instance_name = "${var.name}"
		  tags = {
			Name = "TerraformTest-instance"
		  }
		}

		resource "apsarastack_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${apsarastack_vpc.default.id}"
		  //region_id    = "cn-neimeng-env30-d01"
		}
		
		resource "apsarastack_edas_instance_cluster_attachment" "default" {
		  cluster_id = "${apsarastack_edas_cluster.default.id}"
		  instance_ids = ["${apsarastack_instance.default.id}"]
		  pass_word = "Li65272237###"
		}
		
		resource "apsarastack_edas_application" "default" {
		  application_name = "${var.name}"
		  cluster_id = "${apsarastack_edas_cluster.default.id}"
		  package_type = "JAR"
		  ecu_info = ["${apsarastack_edas_instance_cluster_attachment.default.ecu_map[apsarastack_instance.default.id]}"]
		  build_pack_id = "15"
		}
		`, name)
}
