package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_mongodb_instance", &resource.Sweeper{
		Name: "alibabacloudstack_mongodb_instance",
		F:    testSweepMongoDBInstances,
	})
}

func testAccCheckMongoDBInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_mongodb_instance" {
			continue
		}
		_, err := ddsService.DescribeMongoDBInstance(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		return err
	}
	return nil
}

func testSweepMongoDBInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []dds.DBInstance
	request := dds.CreateDescribeDBInstancesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "testSweepMongoDBInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
		}
		response, _ := raw.(*dds.DescribeDBInstancesResponse)
		addDebug(request.GetActionName(), response)

		if response == nil || len(response.DBInstances.DBInstance) < 1 {
			break
		}
		insts = append(insts, response.DBInstances.DBInstance...)

		if len(response.DBInstances.DBInstance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
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

		if skip {
			log.Printf("[INFO] Skipping MongoDB instance: %s (%s)\n", name, id)
			continue
		}
		log.Printf("[INFO] Deleting MongoDB instance: %s (%s)\n", name, id)
		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] Failed to delete MongoDB instance,ID:%v(%v)\n", id, request.GetActionName())
		} else {
			sweeped = true
		}
		addDebug(request.GetActionName(), raw)
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlibabacloudStackMongoDBInstance_classicv3(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-accdbinstance%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testMongoDBInstanceClassicBase(false, "3.4"))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":                 "${data.alibabacloudstack_zones.default.zones[0].id}",
					"db_instance_description": "${var.name}",
					"engine_version":          "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.engine_version}",
					"db_instance_storage":     "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.storage_min}",
					"db_instance_class":       "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.id}",
					"security_ip_list":        []string{"192.168.1.1"},
					"private_connections": []map[string]interface{}{
						{
							"connect_string_prefix": name + "pr1",
							"connect_port":          3778,
						},
						{
							"connect_string_prefix": name + "pr2",
							"connect_port":          3779,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
						"storage_engine":          "WiredTiger",
						"instance_charge_type":    "PostPaid",
						"replication_factor":      "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection": true,
					"public_connections": []map[string]interface{}{
						{
							"connect_string_prefix": name + "pu1",
							"connect_port":          3788,
						},
						{
							"connect_string_prefix": name + "pu2",
							"connect_port":          3789,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_connection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_connection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Close",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackMongoDBInstance_classicv4(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-accdbinstance%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testMongoDBInstanceClassicBase(true, "4.0"))
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":                 "${data.alibabacloudstack_zones.default.zones[0].id}",
					"db_instance_description": "${var.name}",
					"engine_version":          "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.engine_version}",
					"db_instance_storage":     "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.storage_min}",
					"db_instance_class":       "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.0.id}",
					"audit_status":            "Enable",
					"audit_filter": []map[string]interface{}{{
						"role_type": "db",
						"filters":   []string{"update", "delete"},
					}},
					"security_ip_list": []string{"192.168.1.1"},
					"vswitch_id":       "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": name,
						"storage_engine":          "WiredTiger",
						"instance_charge_type":    "PostPaid",
						"replication_factor":      "3",
						"audit_status":            "Enable",
						"audit_filter.#":          "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"audit_filter.*",
						map[string]string{
							"role_type": "db",
							"filters.#": "2",
						},
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"update",
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"delete",
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_public_connection": "true",
					})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_filter": []map[string]interface{}{{
						"role_type": "db",
						"filters":   []string{"update", "delete", "admin"},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":   "Enable",
						"audit_filter.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"audit_filter.*",
						map[string]string{
							"role_type": "db",
							"filters.#": "3",
						},
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"update",
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"delete",
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"admin",
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tde_status":     "enabled",
					"encryption_key": "${alibabacloudstack_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status":     "enabled",
						"encryption_key": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status": "Disabled",
						"audit_filter": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Open",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_action": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_action": "Close",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{
			// 		"tags": map[string]string{
			// 			"Created": "TF",
			// 			"For":     "acceptance test",
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"tags.%":       "2",
			// 			"tags.Created": "TF",
			// 			"tags.For":     "acceptance test",
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_description": "tf-testAccMongoDBInstance_test",
					//					"tags":                    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_description": "tf-testAccMongoDBInstance_test",
						//						"tags.%":                  REMOVEKEY,
						//						"tags.Created":            REMOVEKEY,
						//						"tags.For":                REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_storage": "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.1.storage_min}",
					"db_instance_class":   "${data.alibabacloudstack_mongodb_instance_types.default.instance_types.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "${random_password.password.0.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
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
			// {
			// 	Config: testMongoDBInstance_classic_security_group_id,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"security_group_id": CHECKSET,
			// 		}),
			// 	),
			// },
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_period": []string{"Wednesday"},
					"backup_time":   "11:00Z-12:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "02:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "02:00Z",
						"maintain_end_time":   "03:00Z",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackMongoDBInstance_multiAZ(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	password := getAccTestPassword(12)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithMultiAZ(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_multiAZ_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "4.0",
						"db_instance_storage":  "10",
						"db_instance_class":    "dds.mongo.mid",
						"name":                 "",
						"storage_engine":       "WiredTiger",
						"instance_charge_type": "PostPaid",
						"replication_factor":   "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_action"},
			},
			{
				Config: testMongoDBInstance_multiAZ_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multiAZ_configure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multiAZ_account_password(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multiAZ_security_ip_list(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multiAZ_backup(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multiAZ_together(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                "tf-testAccMongoDBInstance_test_together",
						"account_password":    "inputYourCodeHere",
						"security_ip_list.#":  "2",
						"security_ip_list.0":  "10.168.1.12",
						"security_ip_list.1":  "10.168.1.13",
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
						"backup_period.#":     "2",
						"backup_period.0":     "Tuesday",
						"backup_period.1":     "Wednesday",
						"backup_time":         "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func testMongoDBInstanceClassicBase(enableVpc bool, engineVersion string) func(string) string {
	var vpcString string
	if enableVpc {
		vpcString = VSwitchCommonTestCase
	} else {
		vpcString = DataZoneCommonTestCase
	}
	return func(name string) string {
		return fmt.Sprintf(`
	
variable "name" {
	  default = "%s"
	}
	
	data "alibabacloudstack_mongodb_instance_types" "default" {
		db_instnace_type = "replicate"
		engine_version = "%s"
	}
%s

%s

%s

`, name, engineVersion, KeyCommonTestCase, RandomPasswordTestCase(12, 1), vpcString)
	}
}

const testMongoDBInstance_multiAZ_base = `

data "alibabacloudstack_zones" "default" {
  
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
}`

const testMongoDBInstance_multiAZ_name = `

data "alibabacloudstack_zones" "default" {
  
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  name                = "tf-testAccMongoDBInstance_test"
}`

const testMongoDBInstance_multiAZ_configure = `

data "alibabacloudstack_zones" "default" {
  
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
}`

func testMongoDBInstance_multiAZ_account_password(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {

}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
}`, password)
}

func testMongoDBInstance_multiAZ_security_ip_list(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
}`, password)
}

func testMongoDBInstance_multiAZ_backup(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  backup_period       = ["Wednesday"]
  backup_time         = "11:00Z-12:00Z"
}`, password)
}

func testMongoDBInstance_multiAZ_together(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
  multi                       = true
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multiAZ"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test_together"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12", "10.168.1.13"]
  backup_period       = ["Tuesday", "Wednesday"]
  backup_time         = "10:00Z-11:00Z"
}`, password)
}
