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

	ResourceTest(t, resource.TestCase{
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
							"server_id":   "${alibabacloudstack_ecs_instance.default.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alibabacloudstack_ecs_instance.new.id}",
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
				// delete_protection_validation是本地属性，无法从远端加载
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

	ResourceTest(t, resource.TestCase{
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
							"server_id":   "${alibabacloudstack_ecs_instance.default.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alibabacloudstack_ecs_instance.new.id}",
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

	variable "name" {
		default = "%s"
	}

	data "alibabacloudstack_instance_types" "new" {
		availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		eni_amount = 2
	}

	%s

	resource "alibabacloudstack_ecs_instance" "new" {
		image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type        = "${data.alibabacloudstack_instance_types.new.instance_types[0].id}"
		system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
		system_disk_size     = 40
		system_disk_name     = "test_sys_diskv2"
		security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
		instance_name        = "${var.name}_ecs"
		vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
		zone_id    = data.alibabacloudstack_zones.default.zones.0.id
		is_outdated          = false
		lifecycle {
		ignore_changes = [
			instance_type
		]
		}
	}

	resource "alibabacloudstack_network_interface" "default" {
		count = 1
		name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
		security_groups = [ "${alibabacloudstack_ecs_securitygroup.default.id}" ]
	}

	resource "alibabacloudstack_network_interface_attachment" "default" {
		count = 1
		instance_id = "${alibabacloudstack_ecs_instance.new.id}"
		network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
	}

	resource "alibabacloudstack_slb" "default" {
		name = "${var.name}"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"

	}
	`, name, ECSInstanceCommonTestCase)
}

var testAccSlbMasterSlaveServerGroupCheckMap = map[string]string{
	"name":      CHECKSET,
	"servers.#": "2",
}
