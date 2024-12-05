package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

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

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerVpcCountConfigDependence)

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
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_instance.instance.0.id}",
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

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

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
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alibabacloudstack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_instance.instance.1.id}",
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

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	resource.Test(t, resource.TestCase{
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
							"server_id": "${alibabacloudstack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_instance.instance.1.id}",
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
							"server_id": "${alibabacloudstack_instance.instance.0.id}",
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
							"server_id": "${alibabacloudstack_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alibabacloudstack_instance.instance.1.id}",
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
%s

%s

%s

variable "name" {
  default = "tf-testAccSlbBackendServersVpc"
}
resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name              = "${var.name}"
}
resource "alibabacloudstack_security_group" "group" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}
resource "alibabacloudstack_instance" "instance" {
  image_id                   = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type              = "${local.instance_type_id}"
  instance_name              = "${var.name}"
  count                      = "2"
  security_groups            = "${alibabacloudstack_security_group.group.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
}


data "alibabacloudstack_instance_types" "new" {
 	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
	eni_amount = 2
}
resource "alibabacloudstack_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    security_groups = [ "${alibabacloudstack_security_group.group.id}" ]
}
resource "alibabacloudstack_instance" "new" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${alibabacloudstack_security_group.group.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_network_interface_attachment" "default" {
	count = 1
    instance_id = "${alibabacloudstack_instance.new.0.id}"
    network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages)
}

func resourceBackendServerConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s

variable "name" {
	default = "tf-testAccSlbBackendServersVpc"
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
resource "alibabacloudstack_instance" "instance" {
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

`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages)
}

var slb_vpc = map[string]string{
	"backend_servers.#": "2",
}
