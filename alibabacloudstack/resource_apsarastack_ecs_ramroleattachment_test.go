package alibabacloudstack

import (
	"fmt"
	"os"
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
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeRamRoleAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsRamRole%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsRamRoleAttachmentdependence)
	ResourceTest(t, resource.TestCase{
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
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${local.local_role_name}",
					"instance_ids": []string{"${alibabacloudstack_ecs_instance.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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

func AlibabacloudTestAccEcsRamRoleAttachmentdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

data "alibabacloudstack_ascm_ram_service_roles" "role" {
   product = "ECS"
}

locals {
	create_count = length(data.alibabacloudstack_ascm_ram_service_roles.role.roles) > 0 ? 0 : 1
}

data "alibabacloudstack_ascm_resource_groups" "group" { 
    name_regex = "%s"
}

resource "alibabacloudstack_ascm_ram_service_role" "default" {
  count = local.create_count
  organization_id = "${data.alibabacloudstack_ascm_resource_groups.group.groups.0.organization_id}"
  product_name = "ECS"
}

locals {
	local_role_name = length(data.alibabacloudstack_ascm_ram_service_roles.role.roles) > 0 ? data.alibabacloudstack_ascm_ram_service_roles.role.roles.0.name : alibabacloudstack_ascm_ram_service_role.default[0].ram_roles.0.role_name
}

`, name, ECSInstanceCommonTestCase, os.Getenv("ALIBABACLOUDSTACK_RESOURCE_GROUP_SET"))
}

var ramRoleAttachmentMap = map[string]string{
	"role_name": CHECKSET,
}
