package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasAccessRule0(t *testing.T) {
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
					"user_access_type":  "no_squash",
					"source_cidr_ip":    "1.1.1.1/0",
					"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
					"rw_access_type":    "RDONLY",
					"priority":          2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type":  "no_squash",
						"source_cidr_ip":    "1.1.1.1/0",
						"access_group_name": CHECKSET,
						"rw_access_type":         "RDONLY",
						"priority":          "2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"user_access_type": "root_squash",
					"rw_access_type":   "RDWR",
					"source_cidr_ip":   "1.1.1.2/0",
					"priority":         4,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "root_squash",
						"rw_access_type":   "RDWR",
						"source_cidr_ip":   "1.1.1.2/0",
						"priority":         "4",
					}),
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

var AlibabacloudTestAccNasAccessruleCheckmap = map[string]string{
	"user_access_type": CHECKSET,
	"priority": CHECKSET,
	"access_group_name": CHECKSET,
	"source_cidr_ip": CHECKSET,
	"access_rule_id": CHECKSET,
	"rw_access_type": CHECKSET,
}

func AlibabacloudTestAccNasAccessruleBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_nas_access_group" "default" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
	description = "tf-testAccNasConfig"
}

`, name)
}
