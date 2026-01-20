package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_gpdb_instance", &resource.Sweeper{
		Name: "alibabacloudstack_gpdb_instance",
		F:    testSweepGpdbInstances,
	})
}

func testSweepGpdbInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var instances []gpdb.DBInstance
	request := gpdb.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DescribeDBInstances(request)
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "testSweepGpdbInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		response, _ := raw.(*gpdb.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)

		if response == nil || len(response.Items.DBInstance) < 1 {
			break
		}
		instances = append(instances, response.Items.DBInstance...)

		if len(response.Items.DBInstance) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeper := false
	service := VpcService{client}
	for _, v := range instances {
		id := v.DBInstanceId
		description := v.DBInstanceDescription
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(description), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If description is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			if need, err := service.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping GPDB instance: %s (%s)\n", description, id)
			continue
		}

		// Delete Instance
		request := gpdb.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "gpdb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		raw, err := client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
			return gpdbClient.DeleteDBInstance(request)
		})
		if err != nil {
			log.Printf("[error] Failed to delete GPDB instance, ID:%v(%v)\n", id, request.GetActionName())
		} else {
			sweeper = true
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeper {
		// Waiting 30 seconds to ensure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}

	return nil
}

func TestAccAlibabacloudStackGpdbInstance_classic(t *testing.T) {
	var v DBInstanceAttribute
	resourceId := "alibabacloudstack_gpdb_instance.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-gpdb-storage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceGpdbClassicConfigDependence)

	ResourceTest(t, resource.TestCase{

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"availability_zone":        "${data.alibabacloudstack_zones.gpdb.zones.0.id}",
					"engine":                   "gpdb",
					"engine_version":           "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version}",
					"instance_class":           "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id}",
					"description":              name,
					"db_instance_mode":         "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.db_instance_mode}",
					"db_instance_storage_type": "local_ssd",
					"seg_node_num":             "2",
					"cpu_type":                 "Intel",
					"security_ip_list": []string{"10.168.1.11","10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":              name,
						"db_instance_storage_type": "local_ssd",
						"seg_node_num":             "2",
						"engine_version":           CHECKSET,
						"engine":                   "gpdb",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_class"},
			},
			// change description
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			// change security ips
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12","10.168.1.13"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "2",
					}),
				),
			},
		}})
}

func TestAccAlibabacloudStackGpdbInstance_vpc(t *testing.T) {
	var v DBInstanceAttribute
	resourceId := "alibabacloudstack_gpdb_instance.default"
	serverFunc := func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeGpdbInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-gpdb-vpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceGpdbVpcConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"availability_zone":        "${data.alibabacloudstack_zones.gpdb.zones.0.id}",
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.default.id}",
					"engine":                   "gpdb",
					"engine_version":           "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version}",
					"db_instance_mode":         "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.db_instance_mode}",
					"instance_class":           "${data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id}",
					"db_instance_storage_type": "local_ssd",
					"instance_group_count":     "2",
					"seg_node_num":             "2",
					"description":              name,
					"cpu_type":                 "Intel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			// change description
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			// change security ips
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
		}})
}

func resourceGpdbClassicConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
        data "alibabacloudstack_zones" "gpdb" {
            available_resource_creation = "Gpdb"
        }
		data "alibabacloudstack_gpdb_instance_types" "default" {
		  engine_version = "6.0"
		  status = "Available"
		}`, name)
}

func resourceGpdbVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
%s
		`, resourceGpdbClassicConfigDependence(name), VSwitchCommonTestCase)
}
