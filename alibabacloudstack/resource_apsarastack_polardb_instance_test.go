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
	name := "tf-testaccdbinstanceconfig"
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":              "MySQL",
						"engine_version":      "5.7",
						"db_instance_class":   CHECKSET,
						"db_instance_storage": CHECKSET,
					}),
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
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}
`, VSwitchCommonTestCase, name)
}

func resourcePolardbInstanceMysqlAZConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}
resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}
`, VSwitchCommonTestCase, name)
}

func TestAccAlibabacloudStackPolardbInstanceClassic(t *testing.T) {
	var instance *PolardbDescribedbinstancesResponse

	resourceId := "alibabacloudstack_polardb_dbinstance.default"
	ra := resourceAttrInit(resourceId, PolardbinstanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "Describedbinstances")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testaccdbinstanceconfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbInstanceClassicConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":                   "MySQL",
					"engine_version":           "8.0",
					"db_instance_class":        "rds.mysql.t1.small",
					"db_instance_storage":      "10",
					"zone_id":                  "cn-wulan-env205-amtest205001-a",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl":               "true",
					"tde_status":               "true",
					"encryption":               "true",
					"encryption_key":           "ae14de55-fdf8-4ea9-b0ec-5b05ff4d5340",
					"zone_id":                  "cn-wulan-env205-amtest205001-a",
					"instance_name":            "${var.name}",
					"db_instance_storage_type": "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"encryption": "true",
					}),
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

`, name, VSwitchCommonTestCase)
}

var PolardbinstanceBasicMap = map[string]string{
	"engine":              "MySQL",
	"engine_version":      "8.0",
	"db_instance_class":   CHECKSET,
	"db_instance_storage": "10",
	"instance_name":       "tf-testaccdbinstanceconfig",
	"zone_id":             CHECKSET,
	"connection_string":   CHECKSET,
	"port":                CHECKSET,
}
