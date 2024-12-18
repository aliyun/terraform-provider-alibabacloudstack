package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAlibabacloudStackNetworkInterfacesDataSourceBasic(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlibabacloudStackNetworkInterfacesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_network_interfaces.default"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.vswitch_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.vpc_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.public_ip"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.security_groups"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.description"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.creation_time"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.tags"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.instance_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.zone_id"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.private_ip"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_network_interfaces.default", "interfaces.status"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackNetworkInterfacesDataSourceConfig = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackImages + `


variable "name" {
  default = "networkInterfaceDatasource"
}


resource "alibabacloudstack_vpc" "vpc" {
  name       = var.name
  cidr_block = "192.168.0.0/24"
}

resource "alibabacloudstack_vswitch" "vswitch" {
  name              = var.name
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  vpc_id            = alibabacloudstack_vpc.vpc.id
}

resource "alibabacloudstack_security_group" "group" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.vpc.id
}

data "alibabacloudstack_instance_types" "instance_type" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  eni_amount        = 2
  sorted_by         = "Memory"
}

resource "alibabacloudstack_instance" "instance" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  security_groups   = [alibabacloudstack_security_group.group.id]
  instance_type              = data.alibabacloudstack_instance_types.instance_type.instance_types[0].id
  system_disk_category       = "cloud_efficiency"
  image_id                   = data.alibabacloudstack_images.default.images[0].id
  instance_name              = var.name
  vswitch_id                 = alibabacloudstack_vswitch.vswitch.id

}
//
resource "alibabacloudstack_network_interface" "interface" {
  name            = var.name
  vswitch_id      = alibabacloudstack_vswitch.vswitch.id
  security_groups = [alibabacloudstack_security_group.group.id]
}

resource "alibabacloudstack_network_interface_attachment" "attachment" {
  instance_id          = alibabacloudstack_instance.instance.id
  network_interface_id = alibabacloudstack_network_interface.interface.id
}
data "alibabacloudstack_network_interfaces" "default"  {
  ids = [alibabacloudstack_network_interface_attachment.attachment.id]
  instance_id = alibabacloudstack_instance.instance.id
}`
