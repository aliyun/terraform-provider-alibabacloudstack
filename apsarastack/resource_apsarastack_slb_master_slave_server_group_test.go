package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackSlbMasterSlaveServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "apsarastack_slb_master_slave_server_group.default"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		//module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${apsarastack_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${apsarastack_instance.default.1.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Slave",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccApsaraStackSlbMasterSlaveServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "apsarastack_slb_master_slave_server_group.default.1"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"load_balancer_id": "${apsarastack_slb.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${apsarastack_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${apsarastack_instance.default.1.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Slave",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func resourceMasterSlaveServerGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s
provider "apsarastack" {
	assume_role {}
}
variable "name" {
    default = "%s"
}

resource "apsarastack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
    vpc_id = "${apsarastack_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = data.apsarastack_zones.default.zones.0.id
    name = "${var.name}"
}
resource "apsarastack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"
    security_groups = [ "${apsarastack_security_group.default.id}" ]
}

resource "apsarastack_network_interface_attachment" "default" {
    count = 1
    instance_id = "${apsarastack_instance.default.0.id}"
    network_interface_id = "${element(apsarastack_network_interface.default.*.id, count.index)}"
}
resource "apsarastack_instance" "default" {
    image_id = "${data.apsarastack_images.default.images.0.id}"
    instance_type = "${local.instance_type_id}"
    instance_name = "${var.name}"
    count = "2"
    security_groups = "${apsarastack_security_group.default.*.id}"
    internet_max_bandwidth_out = "10"
    availability_zone = data.apsarastack_zones.default.zones.0.id
    system_disk_category = "cloud_efficiency"
    vswitch_id = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_slb" "default" {
    name = "${var.name}"
    vswitch_id = "${apsarastack_vswitch.default.id}"

}
`, DataApsarastackVswitchZones, DataApsarastackImages, DataApsarastackInstanceTypes_Eni2, name)
}

var testAccSlbMasterSlaveServerGroupCheckMap = map[string]string{
	"name":      CHECKSET,
	"servers.#": "2",
}
