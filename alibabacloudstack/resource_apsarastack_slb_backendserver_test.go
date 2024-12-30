package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbBackendServers_vpc(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerVpcCountConfigDependence)

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
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
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

func TestAccAlibabacloudStackSlbBackendServers_multi_vpc(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb_backend_server.default.1"
	ra := resourceAttrInit(resourceId, slb_vpc)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

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
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_instance.new.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAlibabacloudStackSlbBackendServers_classic(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alibabacloudstack_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(), //testAccCheckSlbBackendServersDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_instance.new.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
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
			{
				//Config: testAccSlbBackendServersClassicUpdateServer,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_ecs_instance.default.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_instance.new.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
		},
	})
}

func buildBackendServersMap(count int) []map[string]interface{} {
	var result []map[string]interface{}

	str := `${alibabacloudstack_instance.instance.%d.id}`
	for i := 0; i < count; i++ {
		tmp := make(map[string]interface{}, 2)
		tmp["server_id"] = fmt.Sprintf(str, i)
		tmp["weight"] = "10"
		result = append(result, tmp)
	}
	return result
}

func resourceBackendServerVpcCountConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
	}

%s

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${alibabacloudstack_vpc_vswitch.default.id}"
}


data "alibabacloudstack_instance_types" "new" {
 	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	eni_amount = 2
}

resource "alibabacloudstack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
    security_groups = [alibabacloudstack_ecs_securitygroup.default.id]
}
resource "alibabacloudstack_instance" "new" {
	image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type        = "${data.alibabacloudstack_instance_types.new.instance_types[0].id}"
	system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	system_disk_size     = 40
	system_disk_name     = "test_sys_disk"
	security_groups      =	[alibabacloudstack_ecs_securitygroup.default.id]
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
resource "alibabacloudstack_network_interface_attachment" "default" {
	count = 1
    instance_id = "${alibabacloudstack_instance.new.id}"
    network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}
`, name, ECSInstanceCommonTestCase)
}

func resourceBackendServerConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_instance" "new" {
	image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
	instance_type        = "${local.default_instance_type_id}"
	system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	system_disk_size     = 40
	system_disk_name     = "test_sys_diskv2"
	security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
	instance_name        = "${var.name}_ecsv2"
	vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
	zone_id    = data.alibabacloudstack_zones.default.zones.0.id
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

var slb_vpc = map[string]string{
	"backend_servers.#": "2",
}
