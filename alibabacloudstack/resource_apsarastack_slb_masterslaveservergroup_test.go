package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbMasterSlaveServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alibabacloudstack_slb_master_slave_server_group.default"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
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
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alibabacloudstack_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alibabacloudstack_instance.default.1.id}",
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

func TestAccAlibabacloudStackSlbMasterSlaveServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alibabacloudstack_slb_master_slave_server_group.default.1"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
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
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alibabacloudstack_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alibabacloudstack_instance.default.1.id}",
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

variable "name" {
    default = "%s"
}

resource "alibabacloudstack_vpc" "default" {
    name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
    vpc_id = "${alibabacloudstack_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = data.alibabacloudstack_zones.default.zones.0.id
    name = "${var.name}"
}
resource "alibabacloudstack_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.default.id}" ]
}

resource "alibabacloudstack_network_interface_attachment" "default" {
    count = 1
    instance_id = "${alibabacloudstack_instance.default.0.id}"
    network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}
resource "alibabacloudstack_instance" "default" {
    image_id = "${data.alibabacloudstack_images.default.images.0.id}"
    instance_type = "${local.instance_type_id}"
    instance_name = "${var.name}"
    count = "2"
    security_groups = "${alibabacloudstack_security_group.default.*.id}"
    internet_max_bandwidth_out = "10"
    availability_zone = data.alibabacloudstack_zones.default.zones.0.id
    system_disk_category = "cloud_efficiency"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb" "default" {
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"

}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackImages, DataAlibabacloudstackInstanceTypes_Eni2, name)
}

var testAccSlbMasterSlaveServerGroupCheckMap = map[string]string{
	"name":      CHECKSET,
	"servers.#": "2",
}
