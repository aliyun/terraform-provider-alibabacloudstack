package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)


func init() {
	resource.AddTestSweepers("alibabacloudstack_polardbx_instance", &resource.Sweeper{
		Name: "alibabacloudstack_polardbx_instance",
		F: func(region string) error { return SweeperFunc(
			sweeperEntry{fn: testSweepPolardbxInstances, precise: []string{"tf-testAcc", "tf_testAcc", "tfacc-polardbx"}, broad: broadPrefixes, fatal: true, waitSec: 30},
		) },
	})
}

// testSweepPolardbxInstances cleans up all PolarDB-X instances matching the given prefixes.
// For each matching instance it:
//  1. Waits for Running state (best-effort)
//  2. Deletes all databases
//  3. Deletes all accounts (handling rights separation if enabled)
//  4. Deletes the instance itself
func testSweepPolardbxInstances(prefixes []string) error {
	rawClient, err := sharedClientForRegion("")
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	polardbxService := PolardbXService{client}

	// List all PolarDB-X instances (paginated)
	type instanceInfo struct {
		Id          string
		Description string
		Status      string
	}
	var instances []instanceInfo
	page := 1
	pageSize := 50
	for {
		request := client.NewCommonRequest("GET", "polardbx", "2020-02-02", "DescribeDBInstances", "")
		request.QueryParams["PageNumber"] = fmt.Sprint(page)
		request.QueryParams["PageSize"] = fmt.Sprint(pageSize)
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			return fmt.Errorf("error listing PolarDB-X instances: %s", err)
		}
		var resp PolardbxDescribedbinstancesResponse
		if err := json.Unmarshal(bresponse.GetHttpContentBytes(), &resp); err != nil {
			return fmt.Errorf("error parsing PolarDB-X instances response: %s", err)
		}
		for _, inst := range resp.DBInstances {
			matched := false
			for _, prefix := range prefixes {
				if strings.HasPrefix(inst.Description, prefix) || strings.HasPrefix(inst.Id, prefix) {
					matched = true
					break
				}
			}
			if matched {
				instances = append(instances, instanceInfo{
					Id:          inst.Id,
					Description: inst.Description,
					Status:      inst.Status,
				})
			}
		}
		if page*pageSize >= resp.TotalNumber {
			break
		}
		page++
	}

	if len(instances) == 0 {
		log.Printf("[INFO] No PolarDB-X instances to sweep")
		return nil
	}
	log.Printf("[INFO] Found %d PolarDB-X instances to sweep", len(instances))

	for _, inst := range instances {
		log.Printf("[INFO] Sweeping PolarDB-X instance %s (%s), status=%s", inst.Id, inst.Description, inst.Status)

		// Step 1: Wait for Running state (best-effort, timeout 5 min)
		if inst.Status != "Running" {
			log.Printf("[INFO] Instance %s is %s, waiting for Running...", inst.Id, inst.Status)
			stateConf := BuildStateConf(
				[]string{"Creating", "NodeCreating", "NodeDeleting", "ClassChanging", "CONFIG_SWITCHING"},
				[]string{"Running"},
				5*time.Minute, 10*time.Second,
				polardbxService.PolardbxDescribedbinstanceStateRefreshFunc(inst.Id, []string{}),
			)
			if _, err := stateConf.WaitForState(); err != nil {
				log.Printf("[WARN] Instance %s did not reach Running state: %s, attempting delete anyway", inst.Id, err)
			}
		}

		// Step 2: Delete all databases
		if err := sweepPolardbxDatabases(client, inst.Id); err != nil {
			log.Printf("[WARN] Error sweeping databases for %s: %s", inst.Id, err)
		}

		// Step 3: Delete all accounts (handles rights separation)
		if err := sweepPolardbxAccounts(client, inst.Id); err != nil {
			log.Printf("[WARN] Error sweeping accounts for %s: %s", inst.Id, err)
		}

		// Step 4: Delete the instance
		time.Sleep(5 * time.Second)
		request := client.NewCommonRequest("POST", "polardbx", "2020-02-02", "DeleteDBInstance", "")
		request.QueryParams["DBInstanceName"] = inst.Id
		if _, err := client.ProcessCommonRequest(request); err != nil {
			log.Printf("[WARN] Error deleting PolarDB-X instance %s: %s", inst.Id, err)
		} else {
			log.Printf("[INFO] Deleted PolarDB-X instance %s", inst.Id)
		}

		// Wait for instance to be deleted before moving to next
		time.Sleep(15 * time.Second)
	}
	return nil
}


func TestAccAlibabacloudStackPolardbxInstance_basic0(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{"private_connection_string": CHECKSET})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardbx-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "testtf",
					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version": "5.7",
					"storage":        "50",
					"vswitch_id":     "${alibabacloudstack_vpc_vswitch.default.id}",
					"gms_node_class": "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}",
					"cn_node_class":  "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}",
					"cn_node_count":  "2",
					"dn_node_class":  "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}",
					"dn_node_count":  "2",
					"compute_parameters": []map[string]interface{}{
						{
							"name":  "CONN_POOL_BLOCK_TIMEOUT",
							"value": "40000",
						},
						{
							"name":  "CONN_POOL_IDLE_TIMEOUT",
							"value": "35",
						},
					},
					"storage_parameters": []map[string]interface{}{
						{
							"name":  "innodb_stats_sample_pages",
							"value": "10",
						},
						{
							"name":  "ft_query_expansion_limit",
							"value": "35",
						},
					},
					"security_groups": []map[string]interface{}{
						{
							"group_name": "test123",
							"ips":        "10.0.0.1,10.0.0.2",
						},
					},
					"private_connection_string_prefix": name + "-pr",
					"private_connection_port":          4001,
					"enable_public_connection":         true,
					"public_connection_string_prefix":  name + "-pu",
					"public_connection_port":           4002,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// "description":                  name,
						"cn_node_count":                "2",
						"dn_node_count":                "2",
						"compute_parameters.#":         "2",
						"compute_parameters.0.name":    CHECKSET,
						"storage_parameters.#":         "2",
						"storage_parameters.0.value":   CHECKSET,
						"security_groups.#":            "1",
						"security_groups.0.group_name": "test123",
						"security_groups.0.ips":        "10.0.0.1,10.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection": false,
					"public_connection_string_prefix": REMOVEKEY,
					"public_connection_port":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_connection_string_prefix": REMOVEKEY,
						"public_connection_port":          REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []map[string]interface{}{
						{
							"group_name": "test123",
							"ips":        "10.0.0.1,192.168.1.1/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#":            "1",
						"security_groups.0.group_name": "test123",
						"security_groups.0.ips":        "10.0.0.1,192.168.1.1/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
					"enable_tde": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
						"enable_tde": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_count": "3",
					"dn_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "3",
						"dn_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dn_node_count": "2",
					"cn_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dn_node_count": "2",
						"cn_node_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_class": "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.1.id}",
					"dn_node_class": "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// "topology_type" does not support readback, "compute_parameters", "storage_parameters" only support reading back locally updated parameters.
				ImportStateVerifyIgnore: []string{"topology_type", "compute_parameters", "storage_parameters"},
			},
		},
	})
}

func resourcePolardbxDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

data "alibabacloudstack_polardbx_instance_types" "cn" {
	sorted_by = "CPU"
	spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
	sorted_by = "CPU"
	spec_type = "DN"
}
	
%s
 `, name, VSwitchCommonTestCase)
}
