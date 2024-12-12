package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasAccessrule0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_nas_accessrule.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccNasAccessruleCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoNasDescribeaccessrulesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snasaccess_rule%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccNasAccessruleBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"user_access": "no_squash",

					"file_system_type": "standard",

					"source_cidr_ip": "1.1.1.1/0",

					"access_group_name": "${{ref(resource, NAS::AccessGroup::2.0.0.5.pre::defaultCWVMZb.AccessGroupName)}}",

					"rw_access": "RDONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"user_access": "no_squash",

						"file_system_type": "standard",

						"source_cidr_ip": "1.1.1.1/0",

						"access_group_name": "${{ref(resource, NAS::AccessGroup::2.0.0.5.pre::defaultCWVMZb.AccessGroupName)}}",

						"rw_access": "RDONLY",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"user_access": "root_squash",

					"rw_access": "RDWR",

					"source_cidr_ip": "1.1.1.2/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"user_access": "root_squash",

						"rw_access": "RDWR",

						"source_cidr_ip": "1.1.1.2/0",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccNasAccessruleCheckmap = map[string]string{

	"user_access": CHECKSET,

	"priority": CHECKSET,

	"access_group_name": CHECKSET,

	"file_system_type": CHECKSET,

	"source_cidr_ip": CHECKSET,

	"rw_access": CHECKSET,

	"access_rule_id": CHECKSET,
}

func AlibabacloudTestAccNasAccessruleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
