package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPrometheusV2Contact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_contact.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePrometheusV2Contact")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return fmt.Sprintf(`
variable "name" {
  default = "%s"
}		
`, name)
	})

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"username": "${var.name}",
					"mobile":   "13812345678",
					"mail":     "test@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": name,
						"mobile":   "13812345678",
						"mail":     "test@example.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"username": "${var.name}_update",
					"mobile":   "13812345677",
					"mail":     "test1@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": fmt.Sprintf("%s_update", name),
						"mobile":   "13812345677",
						"mail":     "test1@example.com",
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
