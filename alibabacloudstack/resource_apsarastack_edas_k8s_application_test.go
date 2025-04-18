package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_edas_k8s_application", &resource.Sweeper{
		Name: "alibabacloudstack_edas_k8s_application",
		F:    testSweepEdasK8sApplication,
	})
}

func testSweepEdasK8sApplication(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	applicationListRq := edas.CreateListApplicationRequest()
	applicationListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(applicationListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas k8s application in service list: %s", err)
	}

	listApplicationResponse, _ := raw.(*edas.ListApplicationResponse)
	if listApplicationResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas k8s application in service list: %s", listApplicationResponse.Message)
		return errmsgs.WrapError(errmsgs.Error(listApplicationResponse.Message))
	}

	for _, v := range listApplicationResponse.ApplicationList.Application {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas application: %s", name)
			continue
		}
		log.Printf("[INFO] delete edas application: %s", name)
		// stop it before delete
		stopAppRequest := edas.CreateStopApplicationRequest()
		stopAppRequest.RegionId = region
		stopAppRequest.AppId = v.AppId

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.StopApplication(stopAppRequest)
		})
		if err != nil {
			return err
		}
		addDebug(stopAppRequest.GetActionName(), raw, stopAppRequest.RoaRequest, stopAppRequest)
		stopAppResponse, _ := raw.(*edas.StopApplicationResponse)
		changeOrderId := stopAppResponse.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, 5*time.Minute, 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return err
			}
		}

		deleteApplicationRequest := edas.CreateDeleteApplicationRequest()
		deleteApplicationRequest.RegionId = region
		deleteApplicationRequest.AppId = v.AppId

		wait := incrementalWait(1*time.Second, 2*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteApplication(deleteApplicationRequest)
			})
			if err != nil {
				if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteApplicationRequest.GetActionName(), raw, deleteApplicationRequest.RoaRequest, deleteApplicationRequest)
			rsp := raw.(*edas.DeleteApplicationResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = errmsgs.Error("Operation cannot be processed because there are running instances.")
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackEdasK8sApplication_basic(t *testing.T) {
	var v *EdasK8sApplcation
	resourceId := "alibabacloudstack_edas_k8s_application.default"
	ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edask8sappb%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
	// region := os.Getenv("ALIBABACLOUDSTACK_REGION")
	image := fmt.Sprintf("cr.registry.inter.env212.shuguang.com/edas-1780333199011301/outlier-sc-p-exception:outlier-sc-p-excepti-1734082849")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "${var.name}",
					"cluster_id":       "29a23803-50a5-4e08-879f-acff2b5cbee9",
					"package_type":     "Image",
					"image_url":        image,
					"replicas":         "3",
					//"internet_slb_protocol": "TCP",
					//"internet_slb_port":     "8080",
					//"internet_target_port":  "18082",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_url": image,
						"replicas":  "2",
						//"internet_slb_protocol": "TCP",
						//"internet_slb_port":     "8080",
						//"internet_target_port":  "18082",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"replicas": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"replicas": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"command":      "/bin/sh",
					"command_args": []string{"-c", "sleep 1000"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command":        "/bin/sh",
						"command_args.#": "2",
						"command_args.0": "-c",
						"command_args.1": "sleep 1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"envs": map[string]string{"a": "b"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"envs.%": "1",
						"envs.a": "b",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"limit_m_cpu":    "500",
					"limit_mem":      "1000",
					"requests_m_cpu": "100",
					"requests_mem":   "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"limit_m_cpu":    "500",
						"limit_mem":      "1000",
						"requests_m_cpu": "100",
						"requests_mem":   "100",
					}),
				),
			},
		},
	})
}

/*
	func TestAccAlibabacloudStackEdasK8sApplicationJar_basic(t *testing.T) {
		var v *edas.Applcation
		resourceId := "alibabacloudstack_edas_k8s_application.default"
		ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
		serviceFunc := func() interface{} {
			return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
		}
		rc := resourceCheckInit(resourceId, &v, serviceFunc)
		rac := resourceAttrCheckInit(rc, ra)

		rand := getAccTestRandInt(1000, 9999)
		testAccCheck := rac.resourceAttrMapUpdateSet()
		name := fmt.Sprintf("tf-testacc-edask8sappb%v", rand)
		testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
		packageUrl := "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/SPRING_CLOUD_PROVIDER.jar"
		updateUrl := "http://edas-bj.oss-cn-beijing.aliyuncs.com/prod/demo/DUBBO_PROVIDER.jar"
		ResourceTest(t, resource.TestCase{
			PreCheck: func() {

				testAccPreCheck(t)
			},

			IDRefreshName: resourceId,
			Providers:     testAccProviders,
			CheckDestroy:  testAccCheckEdasK8sApplicationDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConfig(map[string]interface{}{
						"application_name": "${var.name}",
						"cluster_id":       "${alibabacloudstack_edas_k8s_cluster.default.id}",
						"package_type":     "FatJar",
						"package_url":      packageUrl,
						"jdk":              "Open JDK 8",
						"replicas":         "1",
						"readiness":        `{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"host\":\"\", \"port\":18081}}`,
						"liveness":         `{\"failureThreshold\": 3,\"initialDelaySeconds\": 5,\"successThreshold\": 1,\"timeoutSeconds\": 1,\"tcpSocket\":{\"host\":\"\", \"port\":18081}}`,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"package_type": "FatJar",
							"package_url":  packageUrl,
							"replicas":     "1",
							"jdk":          "Open JDK 8",
							"readiness":    CHECKSET,
							"liveness":     CHECKSET,
						}),
					),
				},

				{
					ResourceName:            resourceId,
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"package_url", "package_version"},
				},

				{
					Config: testAccConfig(map[string]interface{}{
						"readiness": "{}",
						"liveness":  "{}",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"readiness": "{}",
							"liveness":  "{}",
						}),
					),
				},

				{
					Config: testAccConfig(map[string]interface{}{
						"package_url": updateUrl,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"package_url": updateUrl,
						}),
					),
				},

				{
					Config: testAccConfig(map[string]interface{}{
						"replicas": "2",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"replicas": "2",
						}),
					),
				},

				{
					Config: testAccConfig(map[string]interface{}{
						"jdk": "Dragonwell JDK 8",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"jdk": "Dragonwell JDK 8",
						}),
					),
				},

				{
					Config: testAccConfig(map[string]interface{}{
						"package_url": updateUrl,
						"replicas":    "2",
						"jdk":         "Dragonwell JDK 8",
						"readiness":   "{}",
						"liveness":    "{}",
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(map[string]string{
							"package_url": updateUrl,
							"replicas":    "2",
							"jdk":         "Dragonwell JDK 8",
							"readiness":   "{}",
							"liveness":    "{}",
						}),
					),
				},
			},
		})
	}

	func TestAccAlibabacloudStackEdasK8sApplication_multi(t *testing.T) {
		var v *edas.Applcation
		resourceId := "alibabacloudstack_edas_k8s_application.default.1"
		ra := resourceAttrInit(resourceId, edasK8sApplicationBasicMap)
		serviceFunc := func() interface{} {
			return &EdasService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
		}
		rc := resourceCheckInit(resourceId, &v, serviceFunc)
		rac := resourceAttrCheckInit(rc, ra)

		rand := getAccTestRandInt(100, 999)
		testAccCheck := rac.resourceAttrMapUpdateSet()
		name := fmt.Sprintf("tf-testacc-edask8sappm%v", rand)
		testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
		region := os.Getenv("ALIBABACLOUDSTACK_REGION")
		image := fmt.Sprintf("registry-vpc.%s.aliyuncs.com/edas-demo-image/consumer:1.0", region)
		ResourceTest(t, resource.TestCase{
			PreCheck: func() {

				testAccPreCheck(t)
			},

			IDRefreshName: resourceId,
			Providers:     testAccProviders,
			CheckDestroy:  testAccCheckEdasApplicationDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConfig(map[string]interface{}{
						"count":            "2",
						"application_name": "${var.name}-${count.index}",
						"cluster_id":       "${alibabacloudstack_edas_k8s_cluster.default.id}",
						"replicas":         "1",
						"package_type":     "Image",
						"image_url":        image,
					}),
					Check: resource.ComposeTestCheckFunc(
						testAccCheck(nil),
					),
				},
			},
		})
	}
*/
var edasK8sApplicationBasicMap = map[string]string{
	// "application_name": CHECKSET,
	// "cluster_id":       CHECKSET,
}

func testAccCheckEdasK8sApplicationDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasK8sApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		//data "alicloud_zones" default {
		//  available_resource_creation = "VSwitch"
		//}
		//
		//data "alicloud_instance_types" "default" {
		//  availability_zone = data.alicloud_zones.default.zones.0.id
		//  cpu_core_count = 2
		//  memory_size = 4
		//  kubernetes_node_role = "Worker"
		//}
		//
		//resource "alicloud_vpc" "default" {
		//  name = var.name
		//  cidr_block = "10.1.0.0/21"
		//}
		//
		//resource "alicloud_vswitch" "default" {
		//  name = var.name
		//  vpc_id = alicloud_vpc.default.id
		//  cidr_block = "10.1.1.0/24"
		//  availability_zone = data.alicloud_zones.default.zones.0.id
		//}
		//
		//resource "alicloud_cs_managed_kubernetes" "default" {
		//  worker_instance_types = [data.alicloud_instance_types.default.instance_types.0.id]
		//  name = var.name
		//  worker_vswitch_ids = [alicloud_vswitch.default.id]
		//  worker_number = 				"2"
		//  password =                    "Test12345"
		//  pod_cidr =                   	"172.20.0.0/16"
		//  service_cidr =               	"172.21.0.0/20"
		//  worker_disk_size =            "50"
		//  worker_disk_category =        "cloud_ssd"
		//  worker_data_disk_size =       "20"
		//  worker_data_disk_category =   "cloud_ssd"
		//  worker_instance_charge_type = "PostPaid"
		//  slb_internet_enabled =        "true"
		//}
		
		// resource "alibabacloudstack_edas_k8s_cluster" "default" {
		//   cs_cluster_id = "cdca4266cc6ee4e4e984355d72b34956a"
		// }
		`, name)
}
