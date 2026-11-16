package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
			testAccPreCheck(t, resourceId)
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
			testAccPreCheck(t, resourceId)
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
					"encryption_key":           "${local.kms_key}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            name,
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"tde_status":               true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":         "MySQL",
						"instance_name":  name,
						"tde_status":     "true",
						"encryption_key": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": false,
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
			testAccPreCheck(t, resourceId)
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
					"encryption_key":           "${alibabacloudstack_kms_key.key.id}",
					"zone_id":                  "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"enable_ssl":               true,
					"acl":                      "require",
					"tde_status":               true,
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

// TestAccAlibabacloudStackPolardbInstance_multiZone creates a primary/standby instance whose
// standby node is pinned to another availability zone through zone_id_slave1.
//
// ALIBABACLOUDSTACK_POLARDB_MULTI_ZONE_ENGINE selects the engine under test and defaults to MySQL,
// because the compute clusters backing one engine can run out of free hosts in a test environment
// while another engine still has capacity.
func TestAccAlibabacloudStackPolardbInstance_multiZone(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse

	engine := os.Getenv("ALIBABACLOUDSTACK_POLARDB_MULTI_ZONE_ENGINE")
	if engine == "" {
		engine = "MySQL"
	}
	standbyZone := os.Getenv("ALIBABACLOUDSTACK_POLARDB_MULTIZONE2")

	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-polardb-mz%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceMultiZoneConfigDependence(engine))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "POST", "polardb", "2024-01-30", "CreateDBInstance", "")
			testAccPreCheckWithMultiAZ(t)
			testAccPreCheckPolardbMultiAZ(t)
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
					"db_instance_storage_type": "${local.polardb_instance_type_0.storage_type}",
					"instance_name":            name,
					"vswitch_id":               "${alibabacloudstack_vpc_vswitch.primary.id}",
					"zone_id":                  "${local.polardb_zone_id_1}",
					"zone_id_slave1":           "${local.polardb_zone_id_2}",
					"security_ips":             []string{"10.168.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":              engine,
						"instance_name":       name,
						"db_instance_class":   CHECKSET,
						"db_instance_storage": CHECKSET,
						"zone_id":             CHECKSET,
						"zone_id_slave1":      CHECKSET,
					}),
					// The state keeps the configured standby zone, so also confirm the backend
					// really placed a slave node outside the primary zone.
					testPolardbAccCheckSlaveZone(resourceId, standbyZone),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Read does not write the standby zones back into the state, so they are
				// absent right after an import even though the backend places the slave node.
				ImportStateVerifyIgnore: []string{"force_restart", "encryption", "period", "auto_renew", "param_group_id", "zone_id_slave1", "zone_id_slave2"},
			},
		},
	})
}

// testPolardbAccCheckSlaveZone verifies that the backend placed the standby node of a
// cross-zone instance outside the primary zone. standbyZoneId is the zone handed to
// zone_id_slave1 and is only compared when the test pins the zones of the instance,
// otherwise any zone other than the primary one proves the cross-zone placement.
func testPolardbAccCheckSlaveZone(n string, standbyZoneId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		PolardbService := PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
		resp, err := PolardbService.Describedbinstances(rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(resp.Items.DBInstance) < 1 {
			return fmt.Errorf("DB instance %s not found", rs.Primary.ID)
		}

		instance := resp.Items.DBInstance[0]
		slaveZoneIds := make([]string, 0)
		for _, slaveZone := range instance.SlaveZones.SlaveZone {
			slaveZoneIds = append(slaveZoneIds, slaveZone.ZoneId)
		}
		log.Printf("[DEBUG] instance %s is placed in zone %s with slave zones %v", rs.Primary.ID, instance.ZoneId, slaveZoneIds)

		if len(slaveZoneIds) == 0 {
			return fmt.Errorf("Instance %s reports no slave zone, the standby node is not placed in another availability zone", rs.Primary.ID)
		}

		if standbyZoneId != "" {
			for _, slaveZoneId := range slaveZoneIds {
				if slaveZoneId == standbyZoneId {
					return nil
				}
			}
			return fmt.Errorf("Instance %s has slave zones %v, none of them is the configured standby zone %s", rs.Primary.ID, slaveZoneIds, standbyZoneId)
		}

		for _, slaveZoneId := range slaveZoneIds {
			if slaveZoneId != instance.ZoneId {
				return nil
			}
		}
		return fmt.Errorf("Instance %s has slave zones %v, all of them are the primary zone %s", rs.Primary.ID, slaveZoneIds, instance.ZoneId)
	}
}

func resourcePolardbInstanceClassicConfigDependence(engine string) func(string) string {
	return polardbInstanceConfigDependence(engine, true)
}

// polardbInstanceConfigDependence renders the scaffolding shared by the primary/standby instance
// tests. withKmsKey adds the KMS key consumed by the encryption related steps; tests that do not
// encrypt skip it.
func polardbInstanceConfigDependence(engine string, withKmsKey bool) func(string) string {
	var parameterGroup string
	// Skip parameter group creation due to service instability in test environments
	// if engine  == "MySQL" {
	// 	parameterGroup = `resource "alibabacloudstack_polardb_parameter_group" "default" {
	// 	  engine = "mysql"
	// 	  engine_version = "8.0"
	// 	  parameter_group_name = var.name
	// 	  parameter_group_desc = var.name
	// 	  parameters = {
	// 	    loose_multi_blocks_ddl_count = "10"
	// 	  }
	// 	}`
	// }
	return func(name string) string {
		kmsKey := ""
		if withKmsKey {
			kmsKey = GetKeyFromEnvTestCase()
		}

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

`, name, engine, SecurityGroupCommonTestCase, kmsKey, parameterGroup)
	}
}

// testAccPreCheckPolardbMultiAZ skips a test when the PolarDB service of the current region does
// not expose at least two availability zones, the minimum for placing a standby node in another
// zone than the primary node.
func testAccPreCheckPolardbMultiAZ(t *testing.T) {
	rawClient, err := sharedClientForRegion(os.Getenv("ALIBABACLOUDSTACK_REGION"))
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	request := client.NewCommonRequest("GET", "polardb", "2024-01-30", "DescribeRegions", "")
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		t.Skipf("Skipping the test case: DescribeRegions of polardb is unavailable, err: %s", err)
	}

	regions := PolardbDescriberegionsResponse{}
	if err = json.Unmarshal(bresponse.GetHttpContentBytes(), &regions); err != nil {
		t.Skipf("Skipping the test case: DescribeRegions of polardb is unavailable, err: %s", err)
	}

	zoneIds := make(map[string]bool)
	for _, region := range regions.Regions.RDSRegion {
		if region.ZoneId != "" {
			zoneIds[region.ZoneId] = true
		}
	}
	if len(zoneIds) < 2 {
		t.Skipf("Skipping the test case: polardb provides %d availability zone(s) in this region, cross-zone placement needs at least 2", len(zoneIds))
	}
}

// resourcePolardbInstanceMultiZoneConfigDependence extends the instance scaffolding with the two
// availability zones used for cross-zone placement. ALIBABACLOUDSTACK_POLARDB_MULTIZONE1 /
// ALIBABACLOUDSTACK_POLARDB_MULTIZONE2 can pin an environment specific zone id; without them the
// first two zones returned by the zones data source are used.
func resourcePolardbInstanceMultiZoneConfigDependence(engine string) func(string) string {
	base := polardbInstanceConfigDependence(engine, false)
	return func(name string) string {
		zone1 := os.Getenv("ALIBABACLOUDSTACK_POLARDB_MULTIZONE1")
		zone2 := os.Getenv("ALIBABACLOUDSTACK_POLARDB_MULTIZONE2")

		zoneExpr1 := "data.alibabacloudstack_zones.default.zones[0].id"
		zoneExpr2 := "data.alibabacloudstack_zones.default.zones[1].id"
		if zone1 != "" {
			zoneExpr1 = fmt.Sprintf("%q", zone1)
		}
		if zone2 != "" {
			zoneExpr2 = fmt.Sprintf("%q", zone2)
		}

		return base(name) + fmt.Sprintf(`
locals {
	polardb_zone_id_1 = %s
	polardb_zone_id_2 = %s
}

resource "alibabacloudstack_vpc_vswitch" "primary" {
	vswitch_name = "${var.name}_vsw_primary"
	vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
	cidr_block   = "172.16.101.0/24"
	zone_id      = "${local.polardb_zone_id_1}"
}
`, zoneExpr1, zoneExpr2)
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
