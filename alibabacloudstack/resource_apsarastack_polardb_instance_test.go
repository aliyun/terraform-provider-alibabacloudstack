package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackPolardbInstanceMysql(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse
	var ips []map[string]interface{}

	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_polardb_mysql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence("MySQL"))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_type":                 "${local.polardb_instance_type_0.cpu_type}",
					"engine":                   "${local.polardb_instance_type_0.engine}",
					"engine_version":           "${local.polardb_instance_type_0.engine_version}",
					"db_instance_class":        "${local.polardb_instance_type_0.id}",
					"db_instance_storage":      "${local.polardb_instance_type_0.storage_min}",
					"param_group_id":           "${alibabacloudstack_polardb_parameter_group.default.id}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            name,
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.default.id}",
					"parameters": []map[string]interface{}{{
						"name":  "show_old_temporals",
						"value": "ON",
					}},
					"security_ips": []string{"10.168.1.1", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":              "MySQL",
						"instance_name":       name,
						"db_instance_class":   CHECKSET,
						"db_instance_storage": CHECKSET,
						"security_ips.#":      "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,     // Resource address
						"parameters.*", // TypeSet attribute path (wildcard `*` represents any element in the collection)
						map[string]string{
							"name":  "show_old_temporals",
							"value": "ON",
						},
					),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "encryption", "period", "auto_renew", "param_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time":     "22:00Z-02:00Z",
					"monitoring_period": 300,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time":     "22:00Z-02:00Z",
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": TfRawString("local.polardb_instance_type_0.storage_min+10"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_class": "${local.polardb_instance_type_1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testPolardbAccCheckSecurityIpExists("alibabacloudstack_polardb_dbinstance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"created": "tf",
						"for":     "test acc",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.created": "tf",
						"tags.for":     "test acc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.created": REMOVEKEY,
						"tags.for":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"encryption_key": "${alibabacloudstack_kms_key.key.id}",
					"tde_status":     true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "true",
					}),
				),
			},
		},
	})
}

func testPolardbAccCheckSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		PolardbService := PolardbService{client}
		resp, err := PolardbService.DescribeDBSecurityIps(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s security ip %#v", rs.Primary.ID, resp)

		if err != nil {
			return err
		}

		if len(resp.Items.DBInstanceIPArray) < 1 {
			return fmt.Errorf("DB security ip not found")
		}

		ips = PolardbService.flattenDBSecurityIPs(resp)
		return nil
	}
}

func TestAccAlibabacloudStackPolardbInstanceTDESSL(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse

	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-instance_ted%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence("MySQL"))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_type":                 "${local.polardb_instance_type_0.cpu_type}",
					"engine":                   "${local.polardb_instance_type_0.engine}",
					"engine_version":           "${local.polardb_instance_type_0.engine_version}",
					"db_instance_class":        "${local.polardb_instance_type_0.id}",
					"db_instance_storage":      "${local.polardb_instance_type_0.storage_min}",
					"encryption":               "true",
					"encryption_key":           "${alibabacloudstack_kms_key.key.id}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            name,
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"tde_status":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "MySQL",
						"instance_name":  name,
						"encryption":     "true",
						"tde_status":     "true",
						"encryption_key": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
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
		},
	})
}

func TestAccAlibabacloudStackPolardbInstancePGSql(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse
	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-instance_pgsql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence("PolarDB_PG"))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cpu_type":                 "${local.polardb_instance_type_0.cpu_type}",
					"engine":                   "${local.polardb_instance_type_0.engine}",
					"engine_version":           "${local.polardb_instance_type_0.engine_version}",
					"db_instance_class":        "${local.polardb_instance_type_0.id}",
					"db_instance_storage":      "${local.polardb_instance_type_0.storage_min}",
					"encryption":               "true",
					"encryption_key":           "${alibabacloudstack_kms_key.key.id}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"enable_ssl":               "true",
					"acl":                      "require",
					"tde_status":               "true",
					"parameters": []map[string]interface{}{{
						"name":  "polar_px_interconnect_transmit_timeout",
						"value": "4000",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":        "PolarDB_PG",
						"instance_name": name,
						"acl":           "require",
						"encryption":    "true",
						"tde_status":    "true",
						"enable_ssl":    "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"parameters.*",
						map[string]string{
							"name":  "polar_px_interconnect_transmit_timeout",
							"value": "4000",
						},
					),
				),
			},
		},
	})
}

func resourcePolardbInstanceClassicConfigDependence(engine string) func(string) string {
	var parameterGroup string
	if engine  == "MySQL" {
		parameterGroup = `resource "alibabacloudstack_polardb_parameter_group" "default" {
		  engine = "mysql"
		  engine_version = "8.0"
		  parameter_group_name = var.name
		  parameter_group_desc = var.name
		  parameters = {
		    loose_multi_blocks_ddl_count = "10"
		  }
		}`
	}
	return func(name string) string {
		return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

variable "engine" {
	default = "%s"
}

%s

data "alibabacloudstack_polardb_instance_types" "intel" {
	engine = var.engine
	cpu_type = "intel"
	sorted_by = "CPU"
}

data "alibabacloudstack_polardb_instance_types" "anyone" {
	engine = var.engine
	sorted_by = "CPU"
}

locals {
	polardb_instance_type_0 = length(data.alibabacloudstack_polardb_instance_types.intel.instance_types) > 1 ? data.alibabacloudstack_polardb_instance_types.intel.instance_types.0 : data.alibabacloudstack_polardb_instance_types.anyone.instance_types.0
	polardb_instance_type_1 = length(data.alibabacloudstack_polardb_instance_types.intel.instance_types) > 1 ? data.alibabacloudstack_polardb_instance_types.intel.instance_types.1 : data.alibabacloudstack_polardb_instance_types.anyone.instance_types.1
}

%s

%s

`, name, engine, SecurityGroupCommonTestCase, KeyCommonTestCase, parameterGroup)
	}
}

var PolardbinstanceBasicMap = map[string]string{
	"engine_version":      CHECKSET,
	"db_instance_class":   CHECKSET,
	"db_instance_storage": CHECKSET,
	"instance_name":       CHECKSET,
	"zone_id":             CHECKSET,
	"connection_string":   CHECKSET,
	"port":                CHECKSET,
}
