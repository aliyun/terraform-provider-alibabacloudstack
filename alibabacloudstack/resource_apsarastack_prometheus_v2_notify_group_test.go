package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPrometheusV2NotifyGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_notify_group.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePrometheusV2NotifyGroup")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-notifygroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePrometheusV2NotifyGroupConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"type":        "WEBHOOK",
					"description": "${var.name}",
					"webhook_url": "xxxxxxxxxxxxx",
					"webhook_header_params": []map[string]interface{}{
						{
							"key":   "aaaa",
							"value": "1111",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                          name,
						"type":                          "WEBHOOK",
						"description":                   name,
						"webhook_url":                   "xxxxxxxxxxxxx",
						"webhook_header_params.#":       "1",
						"webhook_header_params.0.key":   "aaaa",
						"webhook_header_params.0.value": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
					"webhook_url": "yyyyyyyyyyyyyyyyy",
					"webhook_header_params": []map[string]interface{}{
						{
							"key":   "xxxx",
							"value": "2222",
						},
						{
							"key":   "yyyy",
							"value": "3333",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                   name + "_update",
						"webhook_url":                   "yyyyyyyyyyyyyyyyy",
						"webhook_header_params.#":       "2",
						"webhook_header_params.0.key":   "xxxx",
						"webhook_header_params.0.value": "2222",
						"webhook_header_params.1.key":   "yyyy",
						"webhook_header_params.1.value": "3333",
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

func TestAccAlibabacloudStackPrometheusV2NotifyGroup_DINGDING(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_notify_group.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePrometheusV2NotifyGroup")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-notifygroup-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePrometheusV2NotifyGroupConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":        "${var.name}",
					"type":        "DINGDING",
					"description": "${var.name}",
					"im":          "xxxxxxxxxxxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"type":        "DINGDING",
						"description": name,
						"im":          "xxxxxxxxxxxxx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
					"im":          "yyyyyyyyyyyyyyyyy",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
						"im":          "yyyyyyyyyyyyyyyyy",
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

func resourcePrometheusV2NotifyGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
