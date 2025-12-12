package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackHsmInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_hsm_instance.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &HsmService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeHsmInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_hsm_instance%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceHsmInstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "jnta.SJJ1528",
					"vendor_code":  "jnta",
					"vsm_type":     "gvsm",
					"zone_no":      "${data.alibabacloudstack_zones.default.zones.0.id}",
					"remark":       "test-tf-hsm-instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": "jnta.SJJ1528",
						"vendor_code":  "jnta",
						"vsm_type":     "gvsm",
						"zone_no":      CHECKSET,
						"remark":       "test-tf-hsm-instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "updated_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "updated_remark",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_id":     "vpc-u4j8obg11ilr3btib78n5",
					"vswitch_id": "vsw-u4jeh1248tpm5mc78g9fd",
					"ip":         "172.16.1.100",
					"white_list": "192.168.1.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":     CHECKSET,
						"vswitch_id": CHECKSET,
						"ip":         "172.16.1.100",
						"white_list": "192.168.1.0/24",
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

func resourceHsmInstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

`, name)
}
