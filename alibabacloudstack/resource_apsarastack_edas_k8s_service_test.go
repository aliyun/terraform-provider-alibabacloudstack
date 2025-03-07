package alibabacloudstack

import (
	"fmt"
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

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasK8sServicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"app_id": "${alibabacloudstack_edas_k8s_application.default.id}",
					"name":   "${var.name}",
					"type":   "ClusterIP",
					"service_ports": []map[string]interface{}{
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
						"name":                    name,
						"type":                    "ClusterIP",
						"external_traffic_policy": "Local",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_traffic_policy"},
			},
		},
	})
}

var testAccCheckEdasK8sServicCheckMap = map[string]string{
	"name":                    CHECKSET,
	"type":                    CHECKSET,
	"external_traffic_policy": CHECKSET,
}

func testAccCheckEdasK8sServicDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasK8sSericeDependence(name string) string {

	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		variable "package_url" {
		  default = "http://fileserver.edas.intra.env212.shuguang.com//prod/demo/SPRING_CLOUD_PROVIDER.jar"
		}

		resource "alibabacloudstack_cs_kubernetes" "default" {
			name = var.name
			version 					= "1.20.11-aliyun.1"
			os_type 					= "linux"
			platform 					= "AliyunLinux"
			num_of_nodes 				= "1"
			master_count				= "3"
			master_vswitch_ids   		= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
			master_instance_types 		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
			master_disk_category 		= "cloud_ssd"
			vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
			worker_instance_types 		= ["ecs.n4v2.large"]
			worker_vswitch_ids 		= ["${alibabacloudstack_vpc_vswitch.default.id}"]
			worker_disk_category 		= "cloud_ssd"
			password 					= "Test12345"
			pod_cidr 					= "172.20.0.0/16"
			service_cidr 				= "172.21.0.0/20"
			worker_disk_size 			= "40"
			master_disk_size 			= "40"
			slb_internet_enabled 		= "true"
		}
		
				
		resource "alibabacloudstack_edas_k8s_cluster" "default" {
		  cs_cluster_id = "${alibabacloudstack_cs_kubernetes.default.id}"
		}

		resource "alibabacloudstack_edas_k8s_application" "default" {
			package_type            = "FatJar"
			application_name        = "terraform-test-fatjar"
			application_description = "This is description of application"
			cluster_id              = "${alibabacloudstack_edas_k8s_cluster.default.id}"
			replicas                = 1
			package_url    		 	= "http://fileserver.edas.intra.env212.shuguang.com//prod/demo/SPRING_CLOUD_PROVIDER.jar"
			package_version 		= "2025-02-21 18:46:19"
			jdk             		= "Open JDK 8"
			internet_target_port  	= 18082
			internet_slb_port     	= 8080
			internet_slb_protocol 	= "TCP"
			limit_mem             	= 1024
			limit_m_cpu           	= 1000
			requests_mem          	= 1024
			requests_m_cpu        	= 1000
			command               	= "/bin/sh"
			command_args          	= ["-c", "sleep 1001", ]
			pre_stop              	= "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
			post_start            	= "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
			namespace             	= "default"
		}
		`, name)
}
