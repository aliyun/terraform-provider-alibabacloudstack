package alibabacloudstack

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_adb_db_instance", &resource.Sweeper{
		Name: "alibabacloudstack_adb_db_instance",
		F:    testSweepAdbDbInstances,
	})
}

func testSweepAdbDbInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Println(errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "AlibabacloudStack_adb_db_clusters", action, errmsgs.AlibabacloudStackSdkGoERROR))
			break
		}

		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			log.Println(errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.Items.DBCluster", response))
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DBClusterDescription"])
			id := fmt.Sprint(item["DBClusterId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ADB Instance: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting adb Instance: %s (%s)", name, id)
			action := "DeleteDBCluster"
			conn, err := client.NewAdsClient()
			if err != nil {
				log.Println(errmsgs.WrapError(err))
				break
			}
			request := map[string]interface{}{
				"DBClusterId": id,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if errmsgs.NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				log.Printf("[ERROR] Deleting ADB cluster failed with error: %#v", err)
				return nil
			})
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// 316 Basic Edition ClusterType:analyticdb
func TestAccAlibabacloudStackAdbDbCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackAdbDbClusterMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-adbcluster%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackAdbDbClusterBasicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_category": "${local.adb_instance_types.0.cluster_category}",
					"db_cluster_version":  "3.0",
					"db_node_class":       "${local.adb_instance_types.0.id}",
					"description":         name,
					"db_node_count":       "${local.adb_instance_types.0.node_min}",
					"db_node_storage":     "${local.adb_instance_types.0.storage_min}",
					"mode":                "${local.adb_instance_types.0.mode}",
					"vswitch_id":          "${local.vswtich_id}",
					"cluster_type":        "${local.adb_instance_types.0.cluster_type}",
					"cpu_type":            "${local.adb_instance_types.0.cpu_type}",
					"security_ips":        []string{"10.168.1.11", "10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": CHECKSET,
						"db_node_class":       CHECKSET,
						"description":         name,
						"db_node_count":       CHECKSET,
						"db_node_storage":     CHECKSET,
						"mode":                CHECKSET,
						"vswitch_id":          CHECKSET,
						"security_ips.#":      "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"db_node_class": "${local.adb_instance_types.1.id}",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": TfRawString("local.adb_instance_types.0.node_min+1"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"db_node_storage": TfRawString("local.adb_instance_types.0.storage_min+100"),
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.10", "10.168.1.11"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.11", "10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "2",
					}),
				),
			},
		},
	})
}

var AlibabacloudStackAdbDbClusterMap = map[string]string{
	"auto_renew_period": NOSET,
	//"compute_resource":  "8Core40GB",
	//"connection_string":  CHECKSET,
	"db_cluster_version": "3.0",
	//"elastic_io_resource": "0",
	"maintain_time": CHECKSET,
	//"resource_group_id": CHECKSET,
	"status": "Running",
	//"tags.%":         "0",
	"zone_id": CHECKSET,
}

func AlibabacloudStackAdbDbClusterBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "existed_vswtich_id" {
	default ="%s"
}

%s

data "alibabacloudstack_adb_cluster_types" "intel" {
	status = "Available"
	sorted_by = "CPU"
	cpu_type = "x86"
}

data "alibabacloudstack_adb_cluster_types" "hygon" {
	status = "Available"
	sorted_by = "CPU"
	cpu_type = "hygon"
}

locals {
	adb_instance_types = length(data.alibabacloudstack_adb_cluster_types.intel.ids) > 0 ? data.alibabacloudstack_adb_cluster_types.intel.instance_types : data.alibabacloudstack_adb_cluster_types.hygon.instance_types
	vswtich_id = var.existed_vswtich_id == "" ? alibabacloudstack_vpc_vswitch.default.id : var.existed_vswtich_id
}

`, name, os.Getenv("ALIBABACLOUDSTACK_TEST_EXISTED_VSWITCH_ID"), VSwitchCommonTestCase)
}
