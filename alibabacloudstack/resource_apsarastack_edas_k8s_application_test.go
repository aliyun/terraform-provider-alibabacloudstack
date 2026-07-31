package alibabacloudstack

import (
	"fmt"
	"log"
	"os"
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
				if errmsgs.IsExpectedErrors(err, errmsgs.ThrottlingUser) {
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
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  "${var.name}",
					"cluster_id":        "${local.edas_cluster_id}",
					"logical_region_id": "${local.edas_logical_region_id}",
					"package_type":      "FatJar",
					"package_url":       fmt.Sprintf("http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar", os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN")),
					"package_version":   "2026-07-30 17:17:18",
					"jdk":               "Open JDK 8",
					"replicas":          "2",
					"internet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "18082",
							"port":        "18082",
							"protocol":    "HTTP",
						},
					},
					"internet_external_traffic_policy": "Local",
					"internet_scheduler":               "rr",
					"custom_tolerations": []map[string]interface{}{
						{
							"key":      "node.kubernetes.io/disk-pressure",
							"operator": "Exists",
							"effect":   "NoSchedule",
						},
					},
					"custom_node_affinity_require": []map[string]interface{}{
						{
							"match_expressions": []map[string]interface{}{
								{
									"key":      "node-role.kubernetes.io/control-plane",
									"operator": "DoesNotExist",
								},
							},
						},
					},
					"custom_node_affinity_preferred": []map[string]interface{}{
						{
							"weight": "100",
							"match_expressions": []map[string]interface{}{
								{
									"key":      "kubernetes.io/os",
									"values":   []string{"linux"},
									"operator": "In",
								},
							},
						},
					},
					"custom_pod_affinity_preferred": []map[string]interface{}{
						{
							"weight":        "1",
							"k8s_namespace": []string{"default"},
							"topology_key":  "kubernetes.io/hostname",
							"match_expressions": []map[string]interface{}{
								{
									"key":      "edas.component",
									"values":   []string{"app"},
									"operator": "In",
								},
								{
									"key":      "edas.controlplane",
									"values":   []string{"edas-oam"},
									"operator": "In",
								},
							},
						},
					},
					"custom_pod_ant_affinity_preferred": []map[string]interface{}{
						{
							"weight":        "1",
							"k8s_namespace": []string{"default"},
							"topology_key":  "kubernetes.io/hostname",
							"match_expressions": []map[string]interface{}{
								{
									"key":      "edas.component",
									"values":   []string{"app"},
									"operator": "In",
								},
								{
									"key":      "edas.oam.acname",
									"values":   []string{"${var.name}"},
									"operator": "NotIn",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type":                                            "FatJar",
						"package_url":                                             CHECKSET,
						"replicas":                                                "2",
						"internet_external_traffic_policy":                        "Local",
						"internet_scheduler":                                      "rr",
						"internet_service_port_infos.#":                           "1",
						"internet_service_port_infos.0.target_port":               "18082",
						"internet_service_port_infos.0.port":                      "18082",
						"internet_service_port_infos.0.protocol":                  "HTTP",
						"custom_tolerations.#":                                    "1",
						"custom_tolerations.0.effect":                             "NoSchedule",
						"custom_node_affinity_require.#":                          "1",
						"custom_node_affinity_require.0.match_expressions.#":      "1",
						"custom_node_affinity_preferred.#":                        "1",
						"custom_node_affinity_preferred.0.match_expressions.#":    "1",
						"custom_pod_affinity_preferred.#":                         "1",
						"custom_pod_affinity_preferred.0.match_expressions.#":     "2",
						"custom_pod_affinity_preferred.0.topology_key":            "kubernetes.io/hostname",
						"custom_pod_affinity_preferred.0.weight":                  "1",
						"custom_pod_ant_affinity_preferred.#":                     "1",
						"custom_pod_ant_affinity_preferred.0.match_expressions.#": "2",
						"custom_pod_ant_affinity_preferred.0.topology_key":        "kubernetes.io/hostname",
						"custom_pod_ant_affinity_preferred.0.weight":              "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_aliases": []map[string]interface{}{
						{
							"ip":        "127.0.0.1",
							"hostnames": []string{"alics.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_aliases.#":             "1",
						"host_aliases.0.ip":          "127.0.0.1",
						"host_aliases.0.hostnames.#": "1",
						"host_aliases.0.hostnames.0": "alics.com",
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
					"limit_mem":      "256",
					"requests_m_cpu": "500",
					"requests_mem":   "256",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"limit_m_cpu":    "500",
						"limit_mem":      "256",
						"requests_m_cpu": "500",
						"requests_mem":   "256",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "18082",
							"port":        "18082",
							"protocol":    "HTTP",
						},
					},
					"internet_external_traffic_policy": "Local",
					"internet_scheduler":               "rr",
					"intranet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "8000",
							"port":        "8000",
							"protocol":    "TCP",
						},
					},
					"intranet_external_traffic_policy": "Local",
					"intranet_scheduler":               "rr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_service_port_infos.#":             "1",
						"internet_service_port_infos.0.target_port": "18082",
						"internet_service_port_infos.0.port":        "18082",
						"internet_service_port_infos.1.target_port": REMOVEKEY,
						"internet_service_port_infos.1.port":        REMOVEKEY,
						"internet_service_port_infos.1.protocol":    REMOVEKEY,
						"intranet_service_port_infos.#":             "1",
						"intranet_service_port_infos.0.target_port": "8000",
						"intranet_service_port_infos.0.port":        "8000",
						"intranet_external_traffic_policy":          "Local",
						"intranet_scheduler":                        "rr",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// "intranet_scheduler", "internet_scheduler" cannot be read back
				// affinity/tolerations are not returned by GetK8sApplication, cannot be read back
				ImportStateVerifyIgnore: []string{
					"intranet_scheduler", "internet_scheduler", "logical_region_id",
					"custom_node_affinity_require",
					"custom_node_affinity_preferred",
					"custom_pod_affinity_preferred",
					"custom_pod_ant_affinity_preferred",
					"custom_tolerations",
				},
			},
		},
	})
}

func TestAccAlibabacloudStackEdasK8sApplicationJar_slbbind(t *testing.T) {
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
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  "${var.name}",
					"cluster_id":        "${local.edas_cluster_id}",
					"logical_region_id": "${local.edas_logical_region_id}",
					"package_type":      "FatJar",
					"package_url":       fmt.Sprintf("http://fileserver.edas.%s//prod/demo/SPRING_CLOUD_PROVIDER.jar", os.Getenv("ALIBABACLOUDSTACK_POPGW_DOMAIN")),
					"package_version":   "2026-07-30 17:17:18",
					"jdk":               "Open JDK 8",
					"cr_instance_id":    "cri-private",
					"replicas":          "1",
					"internet_slb_id":   "${alibabacloudstack_slb_loadbalancer.default.id}",
					"internet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "18082",
							"port":        "18082",
							"protocol":    "HTTP",
						},
					},
					"intranet_slb_id": "${alibabacloudstack_slb_loadbalancer.default1.id}",
					"intranet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "8080",
							"port":        "8080",
							"protocol":    "TCP",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_service_port_infos.#":             "1",
						"internet_service_port_infos.0.target_port": "18082",
						"internet_service_port_infos.0.port":        "18082",
						"intranet_service_port_infos.#":             "1",
						"intranet_service_port_infos.0.target_port": "8080",
						"intranet_service_port_infos.0.port":        "8080",
						"intranet_scheduler":                        "rr",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// "intranet_scheduler", "internet_scheduler" cannot be read back
				ImportStateVerifyIgnore: []string{"intranet_scheduler", "cr_instance_id", "internet_scheduler", "logical_region_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "18000",
							"port":        "18000",
							"protocol":    "HTTP",
						},
					},
					"intranet_service_port_infos": []map[string]interface{}{
						{
							"target_port": "8000",
							"port":        "8000",
							"protocol":    "TCP",
						},
					},
					"intranet_scheduler": "rr",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_service_port_infos.#":             "1",
						"internet_service_port_infos.0.target_port": "18000",
						"internet_service_port_infos.0.port":        "18000",
						"intranet_service_port_infos.#":             "1",
						"intranet_service_port_infos.0.target_port": "8000",
						"intranet_service_port_infos.0.port":        "8000",
						"intranet_scheduler":                        "rr",
					}),
				),
			},
		},
	})
}

func TestUatAlibabacloudStackEdasK8sApplication_image(t *testing.T) {
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
	name := fmt.Sprintf("tf-testacc-edask8sappimg%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasK8sApplicationConfigDependence)
	imageUrl := os.Getenv("ALIBABACLOUDSTACK_EDAS_IMAGE_URL")
	if imageUrl == "" {
		t.Skip("ALIBABACLOUDSTACK_EDAS_IMAGE_URL is not set, skipping Image type EDAS K8s application acceptance test")
	}
	updatedImageUrl := os.Getenv("ALIBABACLOUDSTACK_EDAS_IMAGE_URL_UPDATE")
	if updatedImageUrl == "" {
		updatedImageUrl = fmt.Sprintf("%s-updated", imageUrl)
	}
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckEdasK8sApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name":  "${var.name}",
					"cluster_id":        "${local.edas_cluster_id}",
					"logical_region_id": "${local.edas_logical_region_id}",
					"package_type":      "Image",
					"image_url":         imageUrl,
					"replicas":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type": "Image",
						"image_url":    imageUrl,
						"package_url":  "",
						"replicas":     "1",
					}),
				),
			},
			{
				// 验证 image_url 更新后能正确回读
				Config: testAccConfig(map[string]interface{}{
					"image_url": updatedImageUrl,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type": "Image",
						"image_url":    updatedImageUrl,
						"package_url":  "",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"logical_region_id"},
			},
		},
	})
}

var edasK8sApplicationBasicMap = map[string]string{
	// "application_name": CHECKSET,
	// "cluster_id":       CHECKSET,
	"replicas": CHECKSET,
}

func testAccCheckEdasK8sApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	edasService := EdasService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_edas_k8s_application" {
			continue
		}

		_, err := edasService.DescribeEdasK8sApplication(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		return errmsgs.WrapError(errmsgs.Error("EDAS K8s Application still exists"))
	}
	return nil
}

func resourceEdasK8sApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
		
	%s
		
	resource "alibabacloudstack_slb_loadbalancer" "default" {
		name          = "${var.name}_slb"
		vswitch_id    = "${local.k8s_vswitch_id}"
		address_type  = "internet"
		specification = "slb.s2.small"
	}

	resource "alibabacloudstack_slb_loadbalancer" "default1" {
		name          = "${var.name}_slb1"
		vswitch_id    = "${local.k8s_vswitch_id}"
		address_type  = "intranet"
		specification = "slb.s2.small"
	}

`, name, EdasClusterCommonTestCase()) // GeneratePassword(12))
}
