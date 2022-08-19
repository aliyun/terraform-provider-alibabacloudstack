package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccApsaraStackNetworkInterfacesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackNetworkInterfacesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_network_interfaces.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.vpc_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.public_ip"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.security_groups"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.creation_time"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.tags"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.instance_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.zone_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.private_ip"),
					resource.TestCheckNoResourceAttr("data.apsarastack_network_interfaces.default", "interfaces.status"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackNetworkInterfacesDataSourceConfig = DataApsarastackVswitchZones + DataApsarastackImages + `


variable "name" {
  default = "networkInterfaceDatasource"
}


resource "apsarastack_vpc" "vpc" {
  name       = var.name
  cidr_block = "192.168.0.0/24"
}

resource "apsarastack_vswitch" "vswitch" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.apsarastack_zones.default.zones[0].id
  vpc_id            = apsarastack_vpc.vpc.id
}

resource "apsarastack_security_group" "group" {
  name   = var.name
  vpc_id = apsarastack_vpc.vpc.id
}

data "apsarastack_instance_types" "instance_type" {
  availability_zone = data.apsarastack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

resource "apsarastack_instance" "instance" {
  availability_zone = data.apsarastack_zones.default.zones[0].id
  security_groups   = [apsarastack_security_group.group.id]
  instance_type              = data.apsarastack_instance_types.instance_type.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.apsarastack_images.default.images[0].id
  instance_name              = var.name
  vswitch_id                 = apsarastack_vswitch.vswitch.id

}
//
resource "apsarastack_network_interface" "interface" {
  name            = var.name
  vswitch_id      = apsarastack_vswitch.vswitch.id
  security_groups = [apsarastack_security_group.group.id]
}

resource "apsarastack_network_interface_attachment" "attachment" {
  instance_id          = apsarastack_instance.instance.id
  network_interface_id = apsarastack_network_interface.interface.id
}
data "apsarastack_network_interfaces" "default"  {
  ids = [apsarastack_network_interface_attachment.attachment.id]
  instance_id = apsarastack_instance.instance.id
}`
