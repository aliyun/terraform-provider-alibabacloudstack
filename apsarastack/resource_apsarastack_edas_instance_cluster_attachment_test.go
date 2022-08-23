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

func TestAccApsaraStackEdasinstanceClusterAttachment_basic(t *testing.T) {
	var v *edas.Cluster
	resourceId := "apsarastack_edas_instance_cluster_attachment.default"
	ra := resourceAttrInit(resourceId, edasICAttachmentMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasicattachment%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasICAttachmentDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testEdasCheckICAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":   "${apsarastack_edas_cluster.default.id}",
					"instance_ids": []string{"${apsarastack_instance.default.id}"},
					"pass_word":    "Li65272237###",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasICAttachmentMap = map[string]string{
	"cluster_id": CHECKSET,
}

func testEdasCheckICAttachmentDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasICAttachmentDependence(name string) string {
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
		  availability_zone = "cn-qingdao-apsara-amtest3001-a"
		  //availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  name = "${var.name}"
		}
		
		resource "apsarastack_security_group" "default" {
		  name = "${var.name}"
		  description = "New security group"
		  vpc_id = "${apsarastack_vpc.default.id}"
		}
		
		resource "apsarastack_instance" "default" {
		  vswitch_id = "${apsarastack_vswitch.default.id}"
		  image_id = "centos_7_7_x64_20G_alibase_20200426.vhd"
		  availability_zone = "cn-qingdao-apsara-amtest3001-a"
		  //availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		  system_disk_category = "cloud_efficiency"
		  system_disk_size ="60"
		  instance_type = "ecs.n4v2.xlarge"
		
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
		  //region_id    = "cn-qingdao-apsara-d01"
		}
		`, name)
}
