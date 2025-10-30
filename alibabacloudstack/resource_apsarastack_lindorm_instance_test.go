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
					"zone_id":             "cn-wulan-env2g1-amtest2001-a",
					"instance_alias":      "${var.name}",
					"cpu_brand":           "Intel",
					"disk_category":       "HDD",
					"local_disk_size":     "200",
					"engine_type":         "lindorm",
					"instance_type":       "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}",
					"vpc_id":              "vpc-ud5v703k1q97fck6d3aqq",
					"vswitch_id":          "vsw-ud5mcgtcrpssn5d9xi5pm",
					"lindorm_num":         2,
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      name,
						"cpu_brand":           "Intel",
						"engine_type":         "lindorm",
						"lindorm_num":         "2",
						"instance_status":     "ACTIVATION",
						"deletion_protection": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disk_category", "local_disk_num", "local_disk_size"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_alias":      "${var.name}_update",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_alias":      fmt.Sprintf("%s_update", name),
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lindorm_num":   "3",
					"instance_type": "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[1].name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lindorm_num":   "3",
						"instance_type": "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[1].name}",
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

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lindorm"
}

`, name)
}
