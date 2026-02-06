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
	name := fmt.Sprintf("tf-testacc-polardb-instance_mysql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceConfigDependence)
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
					"engine":                   "MySQL",
					"engine_version":           "5.7",
					"db_instance_class":        "rds.mysql.t1.small",
					"db_instance_storage":      "5",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.default.id}",
					"db_instance_storage_type": "local_ssd",
					"parameters": []map[string]interface{}{{
						"name":  "show_old_temporals",
						"value": "ON",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":       name,
						"engine":              "MySQL",
						"engine_version":      "5.7",
						"db_instance_class":   CHECKSET,
						"db_instance_storage": CHECKSET,
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
				ImportStateVerifyIgnore: []string{"force_restart", "encryption", "period", "auto_renew"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "22:00Z-02:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "22:00Z-02:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testaccdbinstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testaccdbinstance_instance_name",
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
		},
	})
}

func TestAccAlibabacloudStackPolardbInstanceMysqlParamsGroup(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse

	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-instance_mysql%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceConfigDependence)
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
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"param_group_id":           "${alibabacloudstack_polardb_parameter_group.default.id}",
					"db_instance_class":        "rds.mysql.t1.small",
					"db_instance_storage":      "5",
					"instance_name":            "${var.name}",
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.default.id}",
					"db_instance_storage_type": "local_ssd",
					"parameters": []map[string]interface{}{{
						"name":  "show_old_temporals",
						"value": "ON",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":       name,
						"engine":              "MySQL",
						"engine_version":      "8.0",
						"db_instance_class":   CHECKSET,
						"db_instance_storage": CHECKSET,
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

func resourcePolardbInstanceConfigDependence(name string) string {
	param_group_name := generateRandomString(10)
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_polardb_parameter_group" "default" {
  engine = "mysql"
  engine_version = "8.0"
  parameter_group_name = "%s"
  parameter_group_desc = var.name
  parameters = {
    loose_multi_blocks_ddl_count = "1"
  }
}

`, VSwitchCommonTestCase, name, param_group_name)
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
	name := fmt.Sprintf("tf-testacc-polardb-instance_mysql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence)
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
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"db_instance_class":        "rds.mysql.t1.small",
					"db_instance_storage":      "10",
					"encryption":               "true",
					"encryption_key":           "${alibabacloudstack_kms_key.key.id}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"encryption":    "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
					"tde_status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status":     "true",
						"enable_ssl":     "true",
						"encryption_key": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackPolardbInstancePGSql(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse
	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstancePGSqlMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-polardb-instance_pgsql%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence)
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
					"engine":                   "PolarDB_PG",
					"engine_version":           "14",
					"db_instance_class":        "polardb.x4.medium.2",
					"db_instance_storage":      "10",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
					"cpu_type":                 "intel",
					"enable_ssl":               "true",
					"acl":                      "require",
					"tde_status":               "true",
					"encryption":               "true",
					"encryption_key":           "${alibabacloudstack_kms_key.key.id}",
					"parameters": []map[string]interface{}{{
						"name":  "polar_px_interconnect_transmit_timeout",
						"value": "4000",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"cpu_type":      "intel",
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

func resourcePolardbInstanceClassicConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}
%s
resource "alibabacloudstack_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
}

`, name, DataZoneCommonTestCase)
}

var PolardbinstanceBasicMap = map[string]string{
	"engine":              "MySQL",
	"engine_version":      "8.0",
	"db_instance_class":   CHECKSET,
	"db_instance_storage": "10",
	"instance_name":       CHECKSET,
	"zone_id":             CHECKSET,
	"connection_string":   CHECKSET,
	"port":                CHECKSET,
}

var PolardbinstancePGSqlMap = map[string]string{
	"engine":              "PolarDB_PG",
	"engine_version":      "14",
	"db_instance_class":   CHECKSET,
	"db_instance_storage": "10",
	"instance_name":       CHECKSET,
	"zone_id":             CHECKSET,
	"connection_string":   CHECKSET,
	"port":                CHECKSET,
}
