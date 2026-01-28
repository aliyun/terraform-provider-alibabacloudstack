package alibabacloudstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEdasK8sSerice_basic(t *testing.T) {
	var v *EdasK8sService
	resourceId := "alibabacloudstack_edas_k8s_service.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeEdasK8sService")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sSericeDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckEdasK8sServicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alibabacloudstack_edas_k8s_application.default.id}",
					"service_name": "${var.name}",
					"type":         "ClusterIP",
					"port_mappings": []map[string]interface{}{
						{
							"service_port": "80",
							"target_port":  "8080",
							"protocol":     "TCP",
						},
					},
					"labels": map[string]string{
						"Creator": "TF",
					},
					"annotations": map[string]string{
						"Creator": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name": name,
						"type":         "ClusterIP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"port_mappings": []map[string]interface{}{
						{
							"service_port": "8000",
							"target_port":  "8800",
							"protocol":     "TCP",
						},
					},
					"labels": map[string]string{
						"Creator": "TF",
						"Foo":     "Bar",
					},
					"annotations": map[string]string{
						"Creator": "TF",
						"Foo":     "Bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port_mappings.0.service_port": "8000",
						"port_mappings.0.target_port":  "8800",
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
					"labels":      REMOVEKEY,
					"annotations": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackEdasK8sSerice_NodePort(t *testing.T) {
	var v *EdasK8sService
	resourceId := "alibabacloudstack_edas_k8s_service.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeEdasK8sService")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(0, 1000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tftestacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sSericeDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckEdasK8sServicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id":       "${alibabacloudstack_edas_k8s_application.default.id}",
					"service_name": "${var.name}",
					"type":         "NodePort",
					"port_mappings": []map[string]interface{}{
						{
							"service_port": "80",
							"target_port":  "8080",
							"protocol":     "TCP",
						},
					},
					"external_traffic_policy": "Local",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":            name,
						"type":                    "NodePort",
						"external_traffic_policy": "Local",
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

var testAccCheckEdasK8sServicCheckMap = map[string]string{
	"service_name": CHECKSET,
	"type":         CHECKSET,
}

func testAccCheckEdasK8sServicDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasK8sSericeDependence(name string) string {

	return fmt.Sprintf(`
		variable "name" {
			default = "%v"
		}
		
		%s

		resource "alibabacloudstack_edas_k8s_application" "default" {
			package_type			= "FatJar"
			application_name		= "terraform-test-fatjar"
			application_description	= "This is description of application"
			cluster_id				= local.edas_cluster_id
			replicas				= 1
			package_url				= "http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar"
			package_version			= "2025-02-21 18:46:19"
			jdk						= "Open JDK 8"
			internet_target_port	= 18082
			internet_slb_port		= 8080
			internet_slb_protocol	= "TCP"
			limit_mem				= 1024
			limit_m_cpu				= 1000
			requests_mem			= 1024
			requests_m_cpu			= 1000
			command					= "/bin/sh"
			command_args			= ["-c", "sleep 1001", ]
			pre_stop				= "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
			post_start				= "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
			namespace				= "default"
		}
		`, name, EdasClusterCommonTestCase(), os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN"))
}
