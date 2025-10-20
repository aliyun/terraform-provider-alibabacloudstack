package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEdasSwimmingLaneGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_edas_swimming_lane_group.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasSwimmingLaneGroup")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_swimming_lane_group%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSwimmingLaneGroupDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":          "${var.name}",
					"entry_app_id":  "${alibabacloudstack_edas_k8s_application.default[1].id}",
					"apps":          []string{"${alibabacloudstack_edas_k8s_application.default[1].id}", "${alibabacloudstack_edas_k8s_application.default[0].id}"},
					"strategy_type": "CONTENT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":              name,
						"entry_app_id":      CHECKSET,
						"apps.#":            "2",
						"logical_region_id": CHECKSET,
					}),
				),
			},
			// strategy_type not in readfunc result
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"strategy_type"},
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":         "${var.name}_update",
					"entry_app_id": "${alibabacloudstack_edas_k8s_application.default[0].id}",
					"apps":         []string{"${alibabacloudstack_edas_k8s_application.default[0].id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":   fmt.Sprintf("%s_update", name),
						"apps.#": "1",
					}),
				),
			},
		},
	})
}

func resourceEdasSwimmingLaneGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "package_version" {	
	default = "2025-10-17 17:17:18"
}

%s

resource "alibabacloudstack_edas_k8s_application" "default" {
  count                   	= 2
  application_name        	= "testapp${count.index}"
  application_description 	= "This is description of application"
  cluster_id              	= local.edas_cluster_id
  replicas                	= 2
  package_type 				= "FatJar"
  package_url     			= "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
  package_version 			= var.package_version
  jdk             			= "Open JDK 8"
  limit_mem             	= 1024
  requests_mem          	= 1024
  requests_m_cpu        	= 300
  limit_m_cpu           	= 300
}

`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
