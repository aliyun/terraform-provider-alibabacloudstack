package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCSAutoscalingConfig_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cs_autoscaling_config.default"
	ra := resourceAttrInit(resourceId, csAutoscalingConfigBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "GetAutoscalingConfig")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccAutoscalingConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCSAutoscalingConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: Create autoscaling config with default values
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                "${local.k8s_cluster_id}",
					"cool_down_duration":        "10m",
					"unneeded_duration":         "10m",
					"utilization_threshold":     "0.5",
					"gpu_utilization_threshold": "0.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                CHECKSET,
						"cool_down_duration":        "10m",
						"unneeded_duration":         "10m",
						"utilization_threshold":     "0.5",
						"gpu_utilization_threshold": "0.5",
					}),
				),
			},
			// Step 2: Update autoscaling config parameters
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":                "${local.k8s_cluster_id}",
					"cool_down_duration":        "5m",
					"unneeded_duration":         "5m",
					"utilization_threshold":     "0.6",
					"gpu_utilization_threshold": "0.7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                CHECKSET,
						"cool_down_duration":        "5m",
						"unneeded_duration":         "5m",
						"utilization_threshold":     "0.6",
						"gpu_utilization_threshold": "0.7",
					}),
				),
			},
			// Step 3: Import verification
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var csAutoscalingConfigBasicMap = map[string]string{
	"cool_down_duration":        "10m",
	"unneeded_duration":         "10m",
	"utilization_threshold":     "0.5",
	"gpu_utilization_threshold": "0.5",
}

func resourceCSAutoscalingConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

`, name, AckK8sCommonTestCase())
}
