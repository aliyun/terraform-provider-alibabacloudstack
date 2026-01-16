package alibabacloudstack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const EsVersion = "7.10.0_ali1.6.0"
const EsDiskType = "yoda-lvm"
const DataNodeSpec = "1C 2Gi"
const DataNodeAmount = "3"
const DataNodeDisk = "500"

const DataNodeSpecForUpdate = "2C 4Gi"
const DataNodeAmountForUpdate = "5"
const DataNodeDiskForUpdate = "600"

const KibanaNodeSpec = "1C 2Gi"
const KibanaNodeSpecForUpdate = "2C 4Gi"

const MasterNodeSpec = "1C 2Gi"
const MasterNodeAmount = "3"
const MasterNodeSpecForUpdate = "2C 4Gi"
const MasterNodeAmountForUpdate = "3"

const ClientNodeSpec = "1C 2Gi"
const ClientNodeAmount = "2"

const ClientNodeSpecForUpdate = "2C 4Gi"
const ClientNodeAmountForUpdate = "3"

func init() {
	resource.AddTestSweepers("alibabacloudstack_elasticsearch_instance", &resource.Sweeper{
		Name: "alibabacloudstack_elasticsearch_instance",
		F:    testSweepElasticsearch,
	})
}

func testSweepElasticsearch(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting AlibabacloudStack client: %s", err)
	}

	client := rawClient.(*connectivity.AlibabacloudStackClient)
	prefixes := []string{
		"",
		"tf-testAcc",
		"tf_testAcc",
	}

	var instances []elasticsearch.Instance
	req := elasticsearch.CreateListInstanceRequest()
	req.RegionId = client.RegionId
	req.Page = requests.NewInteger(1)
	req.Size = requests.NewInteger(PageSizeLarge)

	for {
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.ListInstance(req)
		})

		if err != nil {
			log.Printf("[ERROR] %s", errmsgs.WrapError(fmt.Errorf("Error listing Elasticsearch instances: %s", err)))
			break
		}

		resp, _ := raw.(*elasticsearch.ListInstanceResponse)
		if resp == nil || len(resp.Result) < 1 {
			break
		}

		instances = append(instances, resp.Result...)

		if len(resp.Result) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.Page)
		if err != nil {
			return err
		}
		req.Page = page
	}

	sweeped := false
	service := VpcService{client}
	for _, v := range instances {
		description := v.Description
		id := v.InstanceId
		skip := true

		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.NetworkConfig.VpcId, v.NetworkConfig.VswitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Elasticsearch Instance: %s (%s)", description, id)
			continue
		}

		log.Printf("[INFO] Deleting Elasticsearch Instance: %s (%s)", description, id)
		req := elasticsearch.CreateDeleteInstanceRequest()
		req.InstanceId = id
		_, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Elasticsearch Instance (%s (%s)): %s", description, id, err)
		} else {
			sweeped = true
		}
	}

	if sweeped {
		// Waiting 30 seconds to ensure these instances have been deleted.
		time.Sleep(30 * time.Second)
	}

	return nil
}

func TestAccAlibabacloudStackElasticsearchInstance_vpc(t *testing.T) {
	var instance map[string]interface{}

	resourceId := "alibabacloudstack_elasticsearch_instance.default"
	ra := resourceAttrInit(resourceId, elasticsearchMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccES-vpc%d", rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchInstanceConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":             "${data.alibabacloudstack_zones.default.zones.0.id}",
					"cpu_type":            "Intel",
					"version":             EsVersion,
					"description":         name,
					"scene":               "normal",
					"data_node_amount":    DataNodeAmount,
					"data_node_spec":      DataNodeSpec,
					"data_node_disk_size": DataNodeDisk,
					"data_node_disk_type": EsDiskType,
					"vswitch_id":          "${alibabacloudstack_vpc_vswitch.default.id}",
					"password":            "${random_password.password.0.result}",
					"monitor_password":    "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":     EsVersion,
						"description": name,
						"scene":       "normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"setting_config": map[string]string{
						"\"action.auto_create_index\"":                      "+.*,-*",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"setting_config.action.auto_create_index":                      "+.*,-*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"setting_config": map[string]string{
						"\"action.destructive_requires_name\"": "true",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"setting_config.action.auto_create_index":                      REMOVEKEY,
						"setting_config.action.destructive_requires_name":              "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"setting_config": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"setting_config.action.destructive_requires_name": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_node_spec":      KibanaNodeSpec,
					"kibana_password":       "${random_password.password.0.result}",
					"master_node_amount":    MasterNodeAmount,
					"master_node_spec":      MasterNodeSpec,
					"master_node_disk_size": "100",
					"master_node_disk_type": EsDiskType,
					"client_node_amount":    ClientNodeAmount,
					"client_node_spec":      ClientNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slb_address":   CHECKSET,
						"domain":        CHECKSET,
						"port":          CHECKSET,
						"status":        CHECKSET,
						"kibana_domain": CHECKSET,
						"kibana_port":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kibana_password": "${random_password.password.3.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Upgrade kibana node configuration
				Config: testAccConfig(map[string]interface{}{
					"kibana_node_spec": KibanaNodeSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Downgrade kibana node configuration
				Config: testAccConfig(map[string]interface{}{
					"kibana_node_spec": KibanaNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Remove kibana node
				Config: testAccConfig(map[string]interface{}{
					"kibana_node_spec": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// FIXME: TF cannot recognize
						"kibana_node_spec": REMOVEKEY,
						"kibana_domain":    REMOVEKEY,
						"kibana_protocol":  REMOVEKEY,
						"kibana_port":      REMOVEKEY,
					}),
				),
			},
			{
				// Upgrade client node configuration
				Config: testAccConfig(map[string]interface{}{
					"client_node_amount": ClientNodeAmountForUpdate,
					"client_node_spec":   ClientNodeSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Downgrade client node configuration
				Config: testAccConfig(map[string]interface{}{
					"client_node_amount": ClientNodeAmount,
					"client_node_spec":   ClientNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Remove client node
				Config: testAccConfig(map[string]interface{}{
					"client_node_amount": REMOVEKEY,
					"client_node_spec":   REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"client_node_amount": REMOVEKEY,
						"client_node_spec":   REMOVEKEY,
					}),
				),
			},
			{
				// Upgrade data node configuration
				Config: testAccConfig(map[string]interface{}{
					"data_node_amount": DataNodeAmountForUpdate,
					"data_node_spec":   DataNodeSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Downgrade data node configuration
				Config: testAccConfig(map[string]interface{}{
					"data_node_amount": DataNodeAmount,
					"data_node_spec":   DataNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Upgrade master node configuration
				Config: testAccConfig(map[string]interface{}{
					//					"master_node_amount": MasterNodeAmountForUpdate,
					"master_node_spec": MasterNodeSpecForUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				// Downgrade master node configuration
				Config: testAccConfig(map[string]interface{}{
					//					"master_node_amount": MasterNodeAmount,
					"master_node_spec": MasterNodeSpec,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "monitor_password", "kibana_password", "scene"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_whitelist": []map[string]interface{}{
						map[string]interface{}{
							"vpc_id": "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
							"ips":    []string{"${alibabacloudstack_vpc_vswitch.default.cidr_block}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_whitelist.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_whitelist": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_whitelist.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"apack_accesslog_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"apack_accesslog_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"thread_pool_search_queue_size": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"thread_pool_search_queue_size": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_routing_allocation_disk_watermark_low": "80%",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_routing_allocation_disk_watermark_low": "80%",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.2.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitor_password": "${random_password.password.3.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTPS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTPS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"protocol": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name[:len(name)-1],
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name[:len(name)-1],
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "monitor_password", "kibana_password", "scene"},
			},
		},
	})
}

var elasticsearchMap = map[string]string{
	"description":                CHECKSET,
	"data_node_spec":             CHECKSET,
	"data_node_amount":           CHECKSET,
	"data_node_disk_size":        CHECKSET,
	"data_node_disk_type":        CHECKSET,
	"status":                     "active",
	"private_whitelist.#":        "0",
	"public_whitelist.#":         "0",
	"enable_public":              "false",
	"kibana_whitelist.#":         "0",
	"kibana_private_whitelist.#": "0",
	"id":                         CHECKSET,
	"domain":                     CHECKSET,
	"port":                       CHECKSET,
}

func resourceElasticsearchInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

%s

%s
		
	`, name, VSwitchCommonTestCase, RandomPasswordTestCase(12, 4))
}
