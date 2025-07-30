package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenInstance0(t *testing.T) {
	var v *CbnDescribecensResponse

	resourceId := "alibabacloudstack_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenInstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribecensRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sceninstance%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%sceninstancemodify%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenInstanceBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"cen_instance_name": name,

					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cen_instance_name": name,

						"description": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"cen_instance_name": modify_name,
					"description":       modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cen_instance_name": modify_name,
						"description":       modify_name,
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

var AlibabacloudTestAccCenInstanceCheckmap = map[string]string{

	"status": CHECKSET,

	"protection_level": CHECKSET,

	"create_time": CHECKSET,

	"cen_id":            CHECKSET,
	"cen_instance_name": CHECKSET,
	"description":       CHECKSET,
}

func AlibabacloudTestAccCenInstanceBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
