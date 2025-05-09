package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_alikafka_instance", &resource.Sweeper{
		Name: "alibabacloudstack_alikafka_instance",
		F:    testSweepAlikafkaInstance,
		// When implemented, these should be removed firstly
		Dependencies: []string{
			"alibabacloudstack_alikafka_sasl_acl",
			"alibabacloudstack_alikafka_topic",
			"alibabacloudstack_alikafka_sasl_user",
		},
	})
}

func testSweepAlikafkaInstance(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "error getting AlibabaCloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = region

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*GetInstanceListResponse)
	service := VpcService{client}
	for _, v := range instanceListResp.InstanceList {

		name := v.Name
		skip := true
		for _, prefix := range prefixes {

			// ServiceStatus equals 5 means the instance is in running status.
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping alikafka instance: %s ", name)
			continue
		}
		if v.ServiceStatus != 10 {
			log.Printf("[INFO] release alikafka instance: %s ", name)

			request := alikafka.CreateReleaseInstanceRequest()
			request.InstanceId = v.InstanceId
			request.ForceDeleteInstance = requests.NewBoolean(true)
			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.ReleaseInstance(request)
			})

			if err != nil {
				log.Printf("[ERROR] Failed to release alikafka instance (%s): %s", name, err)
			}
		}

		log.Printf("[INFO] Delete alikafka instance: %s ", name)
		request2 := alikafka.CreateDeleteInstanceRequest()
		request2.InstanceId = v.InstanceId
		_, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteInstance(request2)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to delete alikafka instance (%s): %s", name, err)
		}
	}

	return nil
}

func TestAccAlibabacloudStackAlikafkaInstance_AnyTunnel(t *testing.T) {
	var v *InstanceVO
	resourceId := "alibabacloudstack_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAlikafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceSimpleDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rc.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":      name,
					"zone_id":   "${data.alibabacloudstack_zones.default.zones.0.id}", //${data.alibabacloudstack_zones.default.zones.0.id}",
					"sasl":      "true",
					"plaintext": "true",
					"spec":      "Broker4C16G",
					"cup_type":  "Intel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"sasl":                      "true",
						"spec":                      "Broker4C16G",
						"vpc_id":                    CHECKSET,
						"vip_type":                  CHECKSET,
						"status":                    CHECKSET,
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
					"num_partitions": 5,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"num_partitions": "5",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackAlikafkaInstance_SingleTunnel(t *testing.T) {
	var v *InstanceVO
	resourceId := "alibabacloudstack_alikafka_instance.default"
	ra := resourceAttrInit(resourceId, alikafkaInstanceBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeAlikafkaInstance")
	rac := resourceAttrCheckInit(rc, ra)
	rand := getAccTestRandInt(10000, 20000)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkainstancebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":       name,
					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
					"zone_id":    "cn-wulan-env149-amtest149001-a", //${data.alibabacloudstack_zones.default.zones.0.id}",
					"sasl":       "true",
					"plaintext":  "true",
					"spec":       "Broker4C16G",
					"cup_type":   "Intel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                      name,
						"sasl":                      "true",
						"plaintext":                 "true",
						"spec":                      "Broker4C16G",
						"vpc_id":                    CHECKSET,
						"vip_type":                  CHECKSET,
						"status":                    CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 			{
			// 				Config: testAccConfig(map[string]interface{}{
			// 					"auto_create_topics_enable": "true",
			// 				}),
			// 				Check: resource.ComposeTestCheckFunc(
			// 					testAccCheck(map[string]string{
			// 						"auto_create_topics_enable": "true",
			// 					}),
			// 				),
			// 			},
		},
	})
}

func resourceAlikafkaInstanceSimpleDependence(name string) string {
	return DataZoneCommonTestCase
}

func resourceAlikafkaInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "192.168.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "192.168.0.0/16"
  zone_id = "cn-wulan-env149-amtest149001-a"
}
`, name)
}

var alikafkaInstanceBasicMap = map[string]string{
	"cup_type":                   CHECKSET,
	"spec":                       CHECKSET,
	"replicas":                   CHECKSET,
	"disk_num":                   CHECKSET,
	"sasl":                       CHECKSET,
	"plaintext":                  CHECKSET,
	"message_max_bytes":          "10000000",
	"num_partitions":             "3",
	"auto_create_topics_enable":  "false",
	"num_io_threads":             "16",
	"queued_max_requests":        "80",
	"replica_fetch_wait_max_ms":  "500",
	"replica_lag_time_max_ms":    "30000",
	"num_network_threads":        "3",
	"log_retention_bytes":        "-1",
	"replica_fetch_max_bytes":    "10000000",
	"num_replica_fetchers":       "4",
	"default_replication_factor": "3",
	"offsets_retention_minutes":  "10080",
	"background_threads":         "10",
}
