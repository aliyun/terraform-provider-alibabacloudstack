package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPrometheusV2Instance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_prometheus_v2_instance.default"
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc_prometheus%d", rand)
	ra := resourceAttrInit(resourceId, map[string]string{
		"cluster_name":     name,
		"cluster_id":       CHECKSET,
		"http_api":         CHECKSET,
		"remote_write_url": CHECKSET,
		"push_gateway_url": CHECKSET,
		"objid":            CHECKSET,
	})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PrometheusService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribePrometheusV2Instance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePrometheusV2InstanceDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_name": "${var.name}",
					"tags":         []string{"test1", "test2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_name": name,
						"tags.#":       "2",
						"tags.0":       "test1",
						"tags.1":       "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": []string{"test1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.#": "1",
						"tags.0": "test1",
						"tags.1": REMOVEKEY,
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

func resourcePrometheusV2InstanceDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
 `, name)
}
