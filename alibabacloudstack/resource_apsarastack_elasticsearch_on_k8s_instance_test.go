package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const DataK8sNodeSpec = "8C 32Gi"
const DataK8sNodeAmount = "2"
const DataK8sNodeDisk = "20"
const DataK8sNodeDiskType = "cloud_ssd"

const DataK8sNodeSpecForUpdate = "8C 32Gi"
const DataK8sNodeAmountForUpdate = "3"
const DataK8sNodeDiskForUpdate = "30"

const DataK8sNodeAmountForMultiZone = "4"
const DefaultZoneK8sAmount = "1"

const KibanaK8sNodeSpec = "2C 4Gi"

const MasterK8sNodeSpec = "4C 16Gi"
const MasterK8sNodeSpecForUpdate = "elasticsearch.sn2ne.xlarge"

const ClientK8sNodeSpec = "4C 16Gi"
const ClientK8sNodeAmount = "2"

const ClientK8sNodeSpecForUpdate = "elasticsearch.sn2ne.xlarge"
const ClientK8sNodeAmountForUpdate = "3"

func init() {
	resource.AddTestSweepers("alibabacloudstack_elasticsearch_k8s_instance", &resource.Sweeper{
		Name: "alibabacloudstack_elasticsearch_k8s_instance",
		F:    testSweepElasticsearch,
	})
}

//func testSweepElasticsearch(region string) error {
//	rawClient, err := sharedClientForRegion(region)
//	if err != nil {
//		return fmt.Errorf("Error getting AlibabacloudStack client: %s", err)
//	}
//
//	client := rawClient.(*connectivity.AlibabacloudStackClient)
//	prefixes := []string{
//		"",
//		"tf-testAcc",
//		"tf_testAcc",
//	}
//
//	var instances []elasticsearch.Instance
//	req := elasticsearch.CreateListInstanceRequest()
//	req.RegionId = client.RegionId
//	req.Page = requests.NewInteger(1)
//	req.Size = requests.NewInteger(PageSizeLarge)
//
//	for {
//		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
//			return elasticsearchClient.ListInstance(req)
//		})
//
//		if err != nil {
//			log.Printf("[ERROR] %s", errmsgs.WrapError(fmt.Errorf("Error listing Elasticsearch instances: %s", err)))
//			break
//		}
//
//		resp, _ := raw.(*elasticsearch.ListInstanceResponse)
//		if resp == nil || len(resp.Result) < 1 {
//			break
//		}
//
//		instances = append(instances, resp.Result...)
//
//		if len(resp.Result) < PageSizeLarge {
//			break
//		}
//
//		page, err := getNextpageNumber(req.Page)
//		if err != nil {
//			return err
//		}
//		req.Page = page
//	}
//
//	sweeped := false
//	service := VpcService{client}
//	for _, v := range instances {
//		description := v.Description
//		id := v.InstanceId
//		skip := true
//
//		for _, prefix := range prefixes {
//			if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
//				skip = false
//				break
//			}
//		}
//		// If a ES description is not set successfully, it should be fetched by vswitch name and deleted.
//		if skip {
//			if need, err := service.needSweepVpc(v.NetworkConfig.VpcId, v.NetworkConfig.VswitchId); err == nil {
//				skip = !need
//			}
//		}
//		if skip {
//			log.Printf("[INFO] Skipping Elasticsearch Instance: %s (%s)", description, id)
//			continue
//		}
//
//		log.Printf("[INFO] Deleting Elasticsearch Instance: %s (%s)", description, id)
//		req := elasticsearch.CreateDeleteInstanceRequest()
//		req.InstanceId = id
//		_, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
//			return elasticsearchClient.DeleteInstance(req)
//		})
//		if err != nil {
//			log.Printf("[ERROR] Failed to delete Elasticsearch Instance (%s (%s)): %s", description, id, err)
//		} else {
//			sweeped = true
//		}
//	}
//
//	if sweeped {
//		// Waiting 30 seconds to eusure these instances have been deleted.
//		time.Sleep(30 * time.Second)
//	}
//
//	return nil
//}

func TestAccAlibabacloudStackElasticsearchK8sInstance_multi(t *testing.T) {
	var instance *elasticsearch.DescribeInstanceResponse

	resourceId := "alibabacloudstack_elasticsearch_k8s_instance.default.1"
	ra := resourceAttrInit(resourceId, elasticsearchk8sMap)

	serviceFunc := func() interface{} {
		return &ElasticsearchService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &instance, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccES%s%d", defaultRegionToTest, rand)
	if len(name) > 30 {
		name = name[:30]
	}
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceElasticsearchK8sInstanceConfigDependence_multi)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          name,
					"vswitch_id":           "${alibabacloudstack_vswitch.default.id}",
					"version":              "7.10.0_ali1.6.0",
					"password":             "Admin123@ascm",
					"data_node_spec":       DataK8sNodeSpec,
					"data_node_amount":     DataK8sNodeAmountForMultiZone,
					"data_node_disk_size":  DataK8sNodeDisk,
					"data_node_disk_type":  DataK8sNodeDiskType,
					"instance_charge_type": string(PostPaid),
					"master_node_spec":     MasterK8sNodeSpec,
					"kibana_node_spec":     KibanaK8sNodeSpec,
					"zone_count":           1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

var elasticsearchk8sMap = map[string]string{
	"description":                   CHECKSET,
	"data_node_spec":                DataNodeSpec,
	"data_node_amount":              DataNodeAmount,
	"data_node_disk_size":           DataNodeDisk,
	"data_node_disk_type":           DataNodeDiskType,
	"instance_charge_type":          string(PostPaid),
	"status":                        "active",
	"private_whitelist.#":           "0",
	"public_whitelist.#":            "0",
	"enable_public":                 "false",
	"kibana_whitelist.#":            "0",
	"enable_kibana_public_network":  "true",
	"kibana_private_whitelist.#":    "0",
	"enable_kibana_private_network": "false",
	"master_node_spec":              "",
	"id":                            CHECKSET,
	"domain":                        CHECKSET,
	"port":                          CHECKSET,
	"kibana_domain":                 CHECKSET,
	"kibana_port":                   CHECKSET,
	"vswitch_id":                    CHECKSET,
}

var AlibabacloudStackElasticsearchK8sMap = map[string]string{
	"id":                   CHECKSET,
	"domain":               CHECKSET,
	"port":                 CHECKSET,
	"kibana_domain":        CHECKSET,
	"kibana_port":          CHECKSET,
	"vswitch_id":           CHECKSET,
	"description":          CHECKSET,
	"instance_charge_type": string(PostPaid),
}

func resourceElasticsearchK8sInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}
	data "alibabacloudstack_zones" "default" {
			available_resource_creation= "VSwitch"
		}	
	resource "alibabacloudstack_vpc" "default" {
 	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
		
	`, ElasticsearchInstanceCommonTestCase, name)
}

func resourceElasticsearchK8sInstanceConfigDependence_multi(name string) string {
	return fmt.Sprintf(`
    %s
	variable "name" {
		default = "%s"
	}
	data "alibabacloudstack_zones" "default" {
			
		}	
	resource "alibabacloudstack_vpc" "default" {
  	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}
		
	`, ElasticsearchInstanceCommonTestCase, name)
}
