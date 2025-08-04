package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbProxy_basic(t *testing.T) {
	var v *PolardbDescribedbproxyResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbproxy-%d", rand)
	resourceId := "alibabacloudstack_polardb_proxy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbDescribedbproxyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbProxyConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":        "${alibabacloudstack_polardb_dbinstance.default.id}",
					"db_proxy_instance_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "1",
						"db_proxy_instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num": "2",
					"effective_time":        "Immediate",
					// "effective_specific_time": "${var.db_proxy_effective_time}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_persist": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_persist": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"causal_consist_read": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"causal_consist_read": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_connect_string_port": "8000",
					"db_proxy_connect_string":      "test1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_connect_string": "test1234",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				//  effective_time effective_specific_time is Execution parameters
				ImportStateVerifyIgnore: []string{"effective_time", "effective_specific_time"},
			},
		},
	})
}

func TestAccAlibabacloudStackPolardbProxy_ReadWriteSpliting(t *testing.T) {
	var v *PolardbDescribedbproxyResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbproxy-%d", rand)
	resourceId := "alibabacloudstack_polardb_proxy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbDescribedbproxyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbProxyReadWriteSplitingDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":                    "${alibabacloudstack_polardb_dbinstance.default.id}",
					"db_proxy_instance_num":             "1",
					"read_write_spliting":               "1",
					"read_only_instance_max_delay_time": "30",
					"read_only_instance_weight": []map[string]interface{}{
						{
							"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
							"weight":         "0",
						},
						{
							"db_instance_id": "${alibabacloudstack_polardb_readonly_instance.default.id}",
							"weight":         "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":                      "1",
						"db_proxy_instance_type":                     CHECKSET,
						"read_write_spliting":                        "1",
						"read_only_instance_max_delay_time":          "30",
						"read_only_instance_weight.#":                "2",
						"read_only_instance_weight.0.db_instance_id": CHECKSET,
						"read_only_instance_weight.0.weight":         "0",
						"read_only_instance_weight.1.db_instance_id": CHECKSET,
						"read_only_instance_weight.1.weight":         "100",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_weight": []map[string]interface{}{
						{
							"db_instance_id": "${alibabacloudstack_polardb_dbinstance.default.id}",
							"weight":         "50",
						},
						{
							"db_instance_id": "${alibabacloudstack_polardb_readonly_instance.default.id}",
							"weight":         "50",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_weight.0.weight": "50",
						"read_only_instance_weight.1.weight": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_write_spliting":       "0",
					"read_only_instance_weight": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_write_spliting":                        "0",
						"read_only_instance_weight.#":                REMOVEKEY,
						"read_only_instance_weight.0.db_instance_id": REMOVEKEY,
						"read_only_instance_weight.0.weight":         REMOVEKEY,
						"read_only_instance_weight.1.db_instance_id": REMOVEKEY,
						"read_only_instance_weight.1.weight":         REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				//  effective_time effective_specific_time is Execution parameters
				ImportStateVerifyIgnore: []string{"effective_time", "effective_specific_time"},
			},
		},
	})
}

func resourcePolardbProxyConfigDependence(name string) string {
	now := time.Now().UTC()
	one_minute_later := now.Add(+1 * time.Minute)
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
		
	variable "db_proxy_effective_time" {
		default = "%v"
	}

	%s
	`, name, one_minute_later.Format("2006-01-02T15:04:05Z"), PolarDBMysqlCommonTestCase(false))
}

func resourcePolardbProxyReadWriteSplitingDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}

	%s

	resource "alibabacloudstack_polardb_readonly_instance" "default" {
		master_db_instance_id = "${alibabacloudstack_polardb_dbinstance.default.id}"
		zone_id = "${alibabacloudstack_polardb_dbinstance.default.zone_id}"
		engine_version = "${alibabacloudstack_polardb_dbinstance.default.engine_version}"
		instance_type = "${alibabacloudstack_polardb_dbinstance.default.instance_type}"
		instance_storage = "${alibabacloudstack_polardb_dbinstance.default.instance_storage}"
		instance_name = "${var.name}"
		db_instance_storage_type = "${alibabacloudstack_polardb_dbinstance.default.storage_type}"
	}
	`, name, PolarDBMysqlCommonTestCase(false))
}
