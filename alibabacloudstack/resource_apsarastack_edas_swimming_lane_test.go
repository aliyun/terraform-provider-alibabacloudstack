package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEdasSwimmingLane_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_edas_swimming_lane.default"

	ra := resourceAttrInit(resourceId, map[string]string{})
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeEdasSwimmingLane")
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf_swimming_lane%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasSwimmingLaneDependence)

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
					"name":      "${var.name}",
					"group_id":  "${alibabacloudstack_edas_swimming_lane_group.default.group_id}",
					"apps":      []string{"${alibabacloudstack_edas_k8s_application.default[1].id}", "${alibabacloudstack_edas_k8s_application.default[0].id}"},
					"priority":  "10",
					"path":      "/root/test",
					"condition": "ADD",
					"rest_items": []map[string]interface{}{
						{
							"name":  "test1",
							"value": "100",
							"type":  "cookie",
							"cond":  "==",
						}, {
							"name":     "test2",
							"value":    "50",
							"type":     "header",
							"cond":     ">=",
							"operator": "mod",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"apps.#":                "2",
						"logical_region_id":     CHECKSET,
						"condition":             "ADD",
						"priority":              "10",
						"path":                  "/root/test",
						"rest_items.#":          "2",
						"rest_items.0.name":     "test1",
						"rest_items.0.value":    "100",
						"rest_items.0.type":     "cookie",
						"rest_items.0.cond":     "==",
						"rest_items.0.operator": "rawvalue",
						"rest_items.1.name":     "test2",
						"rest_items.1.value":    "50",
						"rest_items.1.type":     "header",
						"rest_items.1.cond":     ">=",
						"rest_items.1.operator": "mod",
					}),
				),
			},
			// strategy_type not in readfunc result
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"name":      "${var.name}_update",
					"condition": "OR",
					"apps":      []string{"${alibabacloudstack_edas_k8s_application.default[0].id}"},
					"path":      "/root/test_update",
					"enabled":   true,
					"rest_items": []map[string]interface{}{
						{
							"name":  "test1",
							"value": "90",
							"type":  "cookie",
							"cond":  ">=",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  fmt.Sprintf("%s_update", name),
						"apps.#":                "1",
						"condition":             "OR",
						"path":                  "/root/test_update",
						"enabled":               "true",
						"rest_items.#":          "1",
						"rest_items.0.name":     "test1",
						"rest_items.0.value":    "",
						"rest_items.0.type":     "cookie",
						"rest_items.0.cond":     ">=",
						"rest_items.0.operator": "rawvalue",
						"rest_items.1.name":     REMOVEKEY,
						"rest_items.1.value":    REMOVEKEY,
						"rest_items.1.type":     REMOVEKEY,
						"rest_items.1.cond":     REMOVEKEY,
						"rest_items.1.operator": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func resourceEdasSwimmingLaneDependence(name string) string {
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
  replicas                	= 1
  package_type 				= "FatJar"
  package_url     			= "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
  package_version 			= var.package_version
  jdk             			= "Open JDK 8"
  limit_mem             	= 1024
  requests_mem          	= 1024
  requests_m_cpu        	= 300
  limit_m_cpu           	= 300
}

resource "alibabacloudstack_edas_swimming_lane_group" "default" { 
	name 		  = "${var.name}"
	entry_app_id  = "${alibabacloudstack_edas_k8s_application.default[1].id}"
	apps 		  = ["${alibabacloudstack_edas_k8s_application.default[1].id}", "${alibabacloudstack_edas_k8s_application.default[0].id}"]
	strategy_type = "CONTENT"
}


`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
