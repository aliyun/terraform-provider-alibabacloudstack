package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackSlbServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alibabacloudstack_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupDependence)
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
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupVpcUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupVpcUpdate",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": []string{"${alibabacloudstack_network_interface.default.0.id}"},
							"port":       "70",
							"weight":     "10",
							"type":       "eni",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpc",
						"servers.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAlibabacloudStackSlbServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alibabacloudstack_slb_server_group.default.1"
	ra := resourceAttrInit(resourceId, serverGroupMultiClassicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupMultiVpcDependence)
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
					"count":            "2",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alibabacloudstack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAlibabacloudStackSlbServerGroup_classic(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alibabacloudstack_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMultiClassicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupClassic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceServerGroupClassicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alibabacloudstack_slb.default.id}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alibabacloudstack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupClassicUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassicUpdate",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alibabacloudstack_instance.default.0.id}", "${alibabacloudstack_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alibabacloudstack_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupClassic",
						"servers.#": "2",
					}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func resourceSlbServerGroupDependence(name string) string {
	return fmt.Sprintf(`
%s

%s
provider "alibabacloudstack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}

data "alibabacloudstack_instance_types" "new" {
	eni_amount = 2
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
resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]
  internet_max_bandwidth_out = 10
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_instance" "new" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${alibabacloudstack_security_group.default.*.id}"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alibabacloudstack_instance_types.new.instance_types.0.availability_zones.0}"
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_network_interface_attachment" "default" {
	count = 1
    instance_id = "${alibabacloudstack_instance.new.0.id}"
    network_interface_id = "${element(alibabacloudstack_network_interface.default.*.id, count.index)}"
}
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackImages, name)
}

func resourceServerGroupClassicDependence(name string) string {
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
resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = ["${alibabacloudstack_security_group.default.id}"]

  internet_max_bandwidth_out = "10"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages, name)

}

func resourceSlbServerGroupMultiVpcDependence(name string) string {
	return fmt.Sprintf(`
%s

%s

%s
variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
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
resource "alibabacloudstack_instance" "default" {
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
 security_groups = ["${alibabacloudstack_security_group.default.id}"]
  internet_max_bandwidth_out = "10"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  system_disk_category = "cloud_sperf"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}
`, DataAlibabacloudstackVswitchZones, DataAlibabacloudstackInstanceTypes, DataAlibabacloudstackImages)

}

var serverGroupMap = map[string]string{
	"name":      "tf-server-group",
	"servers.#": "1",
}

var serverGroupMultiClassicMap = map[string]string{
	"servers.#": "2",
}

var serversMap = []map[string]interface{}{
	{
		"server_ids": []string{"${alibabacloudstack_instance.default.0.id}"},
		"port":       "1",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alibabacloudstack_instance.default.1.id}"},
		"port":       "2",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alibabacloudstack_instance.default.2.id}"},
		"port":       "3",
		"weight":     "10",
	},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.3.id}"},
	//	"port":       "4",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.4.id}"},
	//	"port":       "5",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.5.id}"},
	//	"port":       "6",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.6.id}"},
	//	"port":       "7",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.7.id}"},
	//	"port":       "8",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.8.id}"},
	//	"port":       "9",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.9.id}"},
	//	"port":       "10",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.10.id}"},
	//	"port":       "11",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.11.id}"},
	//	"port":       "12",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.12.id}"},
	//	"port":       "13",
	//	"weight":     "10",
	//},
	//{
	//	"server_ids": []string{"${alibabacloudstack_instance.default.13.id}"},
	//	"port":       "14",
	//	"weight":     "10",
	//},
}
