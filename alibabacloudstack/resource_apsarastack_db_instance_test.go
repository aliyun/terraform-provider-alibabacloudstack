package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_db_instance", &resource.Sweeper{
		Name: "alibabacloudstack_db_instance",
		F:    testSweepDBInstances,
	})
}

func testSweepDBInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []rds.DBInstance
	req := rds.CreateDescribeDBInstancesRequest()
	req.RegionId = client.RegionId
	req.Headers = map[string]string{"RegionId": client.RegionId}
	req.QueryParams = map[string]string{"Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	req.PageSize = requests.NewInteger(PageSizeLarge)
	if strings.ToLower(client.Config.Protocol) == "https" {
		req.Scheme = "https"
	} else {
		req.Scheme = "http"
	}
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDBInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving RDS Instances: %s", err)
		}
		resp, _ := raw.(*rds.DescribeDBInstancesResponse)
		if resp == nil || len(resp.Items.DBInstance) < 1 {
			break
		}
		insts = append(insts, resp.Items.DBInstance...)

		if len(resp.Items.DBInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.DBInstanceDescription
		id := v.DBInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}

		}

		if skip {
			log.Printf("[INFO] Skipping RDS Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting RDS Instance: %s (%s)", name, id)
		if len(v.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId) > 0 {
			request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
			if strings.ToLower(client.Config.Protocol) == "https" {
				request.Scheme = "https"
			} else {
				request.Scheme = "http"
			}
			request.DBInstanceId = id
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.QueryParams = map[string]string{"Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			if _, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ReleaseReadWriteSplittingConnection(request)
			}); err != nil {
				log.Printf("[ERROR] ReleaseReadWriteSplittingConnection error: %#v", err)
			} else {
				time.Sleep(5 * time.Second)
			}
		}
		req := rds.CreateDeleteDBInstanceRequest()
		req.DBInstanceId = id
		if strings.ToLower(client.Config.Protocol) == "https" {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}
		req.Headers = map[string]string{"RegionId": client.RegionId}
		req.QueryParams = map[string]string{"Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete RDS Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlibabacloudStackDBInstanceMysql(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "alibabacloudstack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)
	resource.Test(t, resource.TestCase{
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
					"engine":           "MySQL",
					"engine_version":   "5.7",
					"instance_type":    "rds.mysql.t1.small",
					"instance_storage": "30",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${alibabacloudstack_vswitch.default.id}",
					"storage_type":     "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MySQL",
						"engine_version":   "5.7",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart", "encryption"},
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
					"instance_storage": "30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "rds.mysql.t1.small",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("alibabacloudstack_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":               "MySQL",
					"engine_version":       "5.7",
					"instance_type":        "rds.mysql.t1.small",
					"instance_storage":     "30",
					"instance_name":        "tf-testAccDBInstanceConfig",
					"instance_charge_type": "Postpaid",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":            "MySQL",
						"engine_version":    "5.7",
						"instance_type":     CHECKSET,
						"instance_storage":  "30",
						"instance_name":     "tf-testAccDBInstanceConfig",
						"zone_id":           CHECKSET,
						"connection_string": CHECKSET,
						"port":              CHECKSET,
					}),
				),
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"security_ip_mode": SafetyMode,
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"security_ip_mode": SafetyMode,
			//					}),
			//				),
			//			},
		},
	})
}

func resourceDBInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccAlibabacloudStackDBInstanceMultiInstance(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alibabacloudstack_db_instance.default.2"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)

	resource.Test(t, resource.TestCase{
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
					"count":            "3",
					"engine":           "MySQL",
					"engine_version":   "5.7",
					"instance_type":    "rds.mysql.t1.small",
					"instance_storage": "30",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${alibabacloudstack_vswitch.default.id}",
					"storage_type":     "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func TestAccAlibabacloudStackDBInstanceMultiAZ(t *testing.T) {
	var instance = &rds.DBInstanceAttribute{}
	resourceId := "alibabacloudstack_db_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBInstance")
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstance_multiAZ"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysqlAZConfigDependence)
	resource.Test(t, resource.TestCase{
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
					"engine":           "MySQL",
					"engine_version":   "5.7",
					"instance_type":    "rds.mysql.t1.small",
					"instance_storage": "30",
					"zone_id":          "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${alibabacloudstack_vswitch.default.id}",
					"storage_type":     "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_multiAZ",
					}),
				),
			},
		},
	})

}

func resourceDBInstanceMysqlAZConfigDependence(name string) string {
	return fmt.Sprintf(`
%s

variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

resource "alibabacloudstack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alibabacloudstack_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccAlibabacloudStackDBInstanceClassic(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "alibabacloudstack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceClassicConfigDependence)
	resource.Test(t, resource.TestCase{
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
					"engine":           "MySQL",
					"engine_version":   "5.7",
					"instance_type":    "rds.mysql.t1.small",
					"instance_storage": "30",
					"zone_id":          "${data.alibabacloudstack_zones.default.zones[0].id}",
					"instance_name":    "${var.name}",
					"storage_type":     "local_ssd",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func resourceDBInstanceClassicConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
  	available_resource_creation= "Rds"
}`, name)
}

func testAccCheckSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		rdsService := RdsService{client}
		resp, err := rdsService.DescribeDBSecurityIps(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s security ip %#v", rs.Primary.ID, resp)

		if err != nil {
			return err
		}

		if len(resp) < 1 {
			return fmt.Errorf("DB security ip not found")
		}

		ips = rdsService.flattenDBSecurityIPs(resp)
		return nil
	}
}

var instanceBasicMap = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "5.7",
	"instance_type":     CHECKSET,
	"instance_storage":  "30",
	"instance_name":     "tf-testAccDBInstanceConfig",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
}
