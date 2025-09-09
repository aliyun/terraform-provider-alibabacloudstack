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

		// module name
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
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				// delete_protection_validation is a local attribute and cannot be loaded from remote
				// load_balancer_id cannot be read back on the proprietary cloud side temporarily
				ImportStateVerifyIgnore: []string{"delete_protection_validation", "load_balancer_id"},
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
			},
		},
	})
}

func resourceMasterSlaveServerGroupConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}

	%s

	resource "alibabacloudstack_ecs_instance" "new" {
	  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	  instance_type        = "${local.default_instance_type_id}"
	  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	  system_disk_size     = 20
	  system_disk_name     = "test_sys_disk"
	  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
	  instance_name        = "${var.name}_ecs"
	  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
	  is_outdated          = false
	  lifecycle {
	    ignore_changes = [
	      instance_type
	    ]
	  }
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
