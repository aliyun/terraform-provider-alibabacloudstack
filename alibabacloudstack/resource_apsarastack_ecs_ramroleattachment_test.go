package alibabacloudstack

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackRamRoleAttachment_basic(t *testing.T) {
	var v *ecs.DescribeInstanceRamRoleResponse
	resourceId := "alibabacloudstack_ram_role_attachment.default"
	ra := resourceAttrInit(resourceId, ramRoleAttachmentMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckRamRoleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAscm_RamRoleAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckRamRoleAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ascmService := RamService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alibabacloudstack_ram_role_attachment" || rs.Type != "alibabacloudstack_ram_role_attachment" {
			continue
		}
		ascm, err := ascmService.DescribeRamRoleAttachment(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		if ascm.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName != "" {
			return errmsgs.WrapError(errmsgs.Error("resource  still exist"))
		}
	}

	return nil
}

const testAccCheckAscm_RamRoleAttachment = DataAlibabacloudstackVswitchZones + DataAlibabacloudstackInstanceTypes + DataAlibabacloudstackImages + `
variable "name" {
  default = "Test_ram_role_attachment"
}

resource "alibabacloudstack_vpc" "default" {
  name = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  cidr_block = "192.168.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name = var.name
}
resource "alibabacloudstack_security_group" "default" {
  name = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}
resource "alibabacloudstack_instance" "default" {
  image_id = data.alibabacloudstack_images.default.images.0.id
  instance_type = local.instance_type_id
  instance_name = var.name
  security_groups = [alibabacloudstack_security_group.default.id]
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_pperf"
  system_disk_size = 100
  vswitch_id = alibabacloudstack_vswitch.default.id
}

data "alibabacloudstack_ascm_ram_service_roles" "role" {
  product = "ecs"
}
resource "alibabacloudstack_ram_role_attachment" "default" {
   role_name    = data.alibabacloudstack_ascm_ram_service_roles.role.roles.0.name
   instance_ids = [alibabacloudstack_instance.default.id]
}
`

var ramRoleAttachmentMap = map[string]string{
	"role_name": CHECKSET,
}
