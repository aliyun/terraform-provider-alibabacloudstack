package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackLindormInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_lindorm_instance.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeLindormInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-lindorm-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, LindormInstanceCommonTestCase)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":        "cn-wulan-env82-amtest82001-a",
					"instance_alias": "${var.name}",
					"cpu_brand":      "Intel",
					"disk_category":  "LV",
					// "core_spec":      "lindorm.g1.2c2g",
					"engine_type":    "lindorm",
					"engine":         "lindorm.g1.2c8g",
					"vpc_id":         "vpc-z0ze770ob8q2ka1whuzb6",
					"vswitch_id":     "vsw-z0zgufp5ygz5hj9564wka",
					"lindorm_num":    3,
					"core_disk_size": 100,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": name,
						"cpu_brand":      "Intel",
						"disk_category":  "LV",
						// "core_spec":       "lindorm.g1.2c2g",
						"engine_type":     "lindorm",
						"engine":          "lindorm.g1.2c8g",
						"lindorm_num":     "3",
						"core_disk_size":  "100",
						"instance_status": "ACTIVATION",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias": "${var.name}_update",
					}),
				),
			},
		},
	})
}

func LindormInstanceCommonTestCase(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
