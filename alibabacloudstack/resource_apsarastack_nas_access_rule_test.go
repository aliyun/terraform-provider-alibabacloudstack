package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNasAccessRule_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_nas_access_rule.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackNasAccessRule0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeNasAccessRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlibabacloudStackNasAccessRule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackNasAccessRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{

			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"access_group_name": "${alibabacloudstack_nas_access_group.example.access_group_name}",
					"source_cidr_ip":    "168.1.1.0/16",
					"priority":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"source_cidr_ip":    "168.1.1.0/16",
						"priority":          "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"source_cidr_ip": "172.168.1.0/16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip": "172.168.1.0/16",
					}),
				),
			},
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"rw_access_type": "RDONLY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rw_access_type": "RDONLY",
					}),
				),
			},
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"priority": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"priority": "2",
					}),
				),
			},
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"user_access_type": "root_squash",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_access_type": "root_squash",
					}),
				),
			},
			{
				Config: providerCommon + testAccConfig(map[string]interface{}{
					"source_cidr_ip": "168.1.1.0/16",
					"priority":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_cidr_ip": "168.1.1.0/16",
						"priority":       "1",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackNasAccessRule0 = map[string]string{
	"access_rule_id": CHECKSET,
}

func AlibabacloudStackNasAccessRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
	}
`, name)
}
