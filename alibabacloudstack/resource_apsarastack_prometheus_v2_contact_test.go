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
	name := fmt.Sprintf("tfacc-contact-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, nil)

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
						"username": "tf-testAccPrometheusV2Contact",
						"mobile":   "13812345678",
						"mail":     "test@example.com",
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
					"username": "tf-testAccPrometheusV2Contact-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": "tf-testAccPrometheusV2Contact-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mobile": "13987654321",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mobile": "13987654321",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mail": "update@example.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mail": "update@example.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"username":  "${var.name}",
					"mobile":    "13812345678",
					"mail":      "test@example.com",
					"group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username": "tf-testAccPrometheusV2Contact",
						"mobile":   "13812345678",
						"mail":     "test@example.com",
					}),
				),
			},
		},
	})
}
