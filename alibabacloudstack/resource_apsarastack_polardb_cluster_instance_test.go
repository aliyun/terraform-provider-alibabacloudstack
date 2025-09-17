package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbClusterInstance_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_polardb_cluster_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbClusterInstance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardb-shared-%d", rand)
	// name := "tfacc-polardb-shared-14498"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceoPolardbClusterInstanceDependence)

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
					"db_cluster_description": "${var.name}",
					"zone_id":                "${data.alibabacloudstack_zones.default.zones.0.id}",
					"db_type":                "${var.db_type}",
					"db_version":             "${var.db_version}",
					"hot_standby_cluster":    "off",
					"storage_space":          "20",
					"vpc_id":                 "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}",
					"db_node_class":          "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.1.id}",
					"sub_category":           "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.1.sub_category}",
					"storage_type":           "ESSDPL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
						"db_type":                CHECKSET,
						"db_version":             CHECKSET,
						"hot_standby_cluster":    "off",
						"storage_space":          "20",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"db_node_class":          CHECKSET,
						"storage_type":           "ESSDPL1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters", "sub_category"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readonly_node_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readonly_node_num": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class": "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}",
					// "db_read_node_class": "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.3.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": CHECKSET,
						// "db_read_node_class": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class": "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.1.id}",
					// "db_read_node_class": "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": CHECKSET,
						// "db_read_node_class": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"readonly_node_num": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"readonly_node_num": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_standby_cluster": "standby",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_standby_cluster": "standby",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"hot_standby_cluster": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hot_standby_cluster": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "true",
					"tde_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "true",
						"tde_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips_groups": map[string]interface{}{
						"test1": "10.0.0.10,192.168.1.10/24",
						"test2": "10.0.0.1,192.168.1.1/24",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips_groups.%":     "2",
						"security_ips_groups.test1": CHECKSET,
						"security_ips_groups.test2": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						// mysql parameter
						// {
						// 	"name":  "connect_timeout",
						// 	"value": "20",
						// },
						// {
						// 	"name":  "loose_hotspot",
						// 	"value": "ON",
						// },
						// pg parameter
						{
							"name":  "auto_explain.log_analyze",
							"value": "on",
						},
						{
							"name":  "auto_explain.sample_rate",
							"value": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
					// mysql parameter
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	resourceId,
					// 	"parameters.*",
					// 	map[string]string{
					// 		"name":  "connect_timeout",
					// 		"value": "20",
					// 	},
					// ),
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	resourceId,
					// 	"parameters.*",
					// 	map[string]string{
					// 		"name":  "loose_hotspot",
					// 		"value": "ON",
					// 	},
					// ),
					// pg parameter
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"parameters.*",
						map[string]string{
							"name":  "auto_explain.log_analyze",
							"value": "on",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"parameters.*",
						map[string]string{
							"name":  "auto_explain.sample_rate",
							"value": "0",
						},
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": "${var.name}_update",
					"deletion_lock":          1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": fmt.Sprintf("%s_update", name),
						"deletion_lock":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_lock": 0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_lock": "0",
					}),
				),
			},
		},
	})
}

func resourceoPolardbClusterInstanceDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

variable "db_type" {
  //default = "MySQL"
  default = "PostgreSQL"
}

variable "db_version" {
  //default = "8.0"
  default = "14"
}

%s

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by = "CPU"
  sub_category = "normal_exclusive"
}

 `, name, VSwitchCommonTestCase)
}
