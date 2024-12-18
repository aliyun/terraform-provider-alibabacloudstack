package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasMounttarget0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_nas_mounttarget.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasMounttargetCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribemounttargetsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasmount_target%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasMounttargetBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultiUcGvP.VpcId)}}",

					"network_type": "Vpc",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultulTodx.VSwitchId)}}",

					"file_system_id": "${{ref(resource, NAS::FileSystem::5.0.0.10.pre::defaultjVQQ4E.FileSystemId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultiUcGvP.VpcId)}}",

						"network_type": "Vpc",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultulTodx.VSwitchId)}}",

						"file_system_id": "${{ref(resource, NAS::FileSystem::5.0.0.10.pre::defaultjVQQ4E.FileSystemId)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": "Inactive",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultiUcGvP.VpcId)}}",

					"network_type": "Vpc",

					"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultulTodx.VSwitchId)}}",

					"file_system_id": "${{ref(resource, NAS::FileSystem::5.0.0.10.pre::defaultjVQQ4E.FileSystemId)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "${{ref(resource, VPC::VPC::7.0.0.1.pre::defaultiUcGvP.VpcId)}}",

						"network_type": "Vpc",

						"vswitch_id": "${{ref(resource, VPC::VSwitch::4.0.5.4.pre::defaultulTodx.VSwitchId)}}",

						"file_system_id": "${{ref(resource, NAS::FileSystem::5.0.0.10.pre::defaultjVQQ4E.FileSystemId)}}",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNasMounttargetCheckmap = map[string]string{

	"status": CHECKSET,

	"access_group_name": CHECKSET,

	"vswitch_id": CHECKSET,

	"mount_target_extra": CHECKSET,

	"vpc_id": CHECKSET,

	"mount_target_domain": CHECKSET,

	"network_type": CHECKSET,

	"file_system_id": CHECKSET,
}

func AlibabacloudTestAccNasMounttargetBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


variable "azone" {
    default = cn-hangzhou-k
}

variable "region_id" {
    default = cn-hangzhou
}




`, name)
}
