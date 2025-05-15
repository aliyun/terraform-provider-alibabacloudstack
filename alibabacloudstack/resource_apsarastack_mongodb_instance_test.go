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

func TestAccAlibabacloudStackMongoDBInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_classic_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
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
				Config: testMongoDBInstance_classic_ssl_action,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_status": "Open",
					}),
				),
			},
			{
				Config: testMongoDBInstance_classic_ssl_action_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_status": "Closed",
					}),
				),
			},
			//			{
			//				Config: testMongoDBInstance_classic_tags,
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"tags.%":       "2",
			//						"tags.Created": "TF",
			//						"tags.For":     "acceptance test",
			//					}),
			//				),
			//			},
			{
				Config: testMongoDBInstance_classic_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBInstance_test",
						//						"tags.%":       REMOVEKEY,
						//						"tags.Created": REMOVEKEY,
						//						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testMongoDBInstance_classic_configure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			// {
			// 	Config: testMongoDBInstance_classic_account_password,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"account_password": "inputYourCodeHere",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testMongoDBInstance_classic_security_ip_list,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"security_ip_list.#": "1",
			// 			"security_ip_list.0": "10.168.1.12",
			// 		}),
			// 	),
			// },
			//			{
			//				Config: testMongoDBInstance_classic_security_group_id,
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"security_group_id": CHECKSET,
			//					}),
			//				),
			//			},
			// {
			// 	Config: testMongoDBInstance_classic_backup,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"backup_period.#": "1",
			// 			"backup_period.0": "Wednesday",
			// 			"backup_time":     "11:00Z-12:00Z",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testMongoDBInstance_classic_maintain_time,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"maintain_start_time": "02:00Z",
			// 			"maintain_end_time":   "03:00Z",
			// 		}),
			// 	),
			// },
			// {
			// 	Config: testMongoDBInstance_classic_together,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{
			// 			"name":                "tf-testAccMongoDBInstance_test_together",
			// 			"account_password":    "inputYourCodeHere",
			// 			"security_ip_list.#":  "2",
			// 			"security_ip_list.0":  "10.168.1.12",
			// 			"security_ip_list.1":  "10.168.1.13",
			// 			"db_instance_storage": "30",
			// 			"db_instance_class":   "dds.mongo.standard",
			// 			"backup_period.#":     "2",
			// 			"backup_period.0":     "Tuesday",
			// 			"backup_period.1":     "Wednesday",
			// 			"backup_time":         "10:00Z-11:00Z",
			// 			"maintain_start_time": REMOVEKEY,
			// 			"maintain_end_time":   REMOVEKEY,
			// 			"ssl_status":          "Open",
			// 		}),
			// 	),
			// },
		},
	})
}

func TestAccAlibabacloudStackMongoDBInstance_Version4(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_classic_base4,
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
				Config: testMongoDBInstance_classic_tde,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackMongoDBInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	password := GeneratePassword()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_vpc_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
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
				Config: testMongoDBInstance_vpc_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBInstance_vpc_configure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testMongoDBInstance_vpc_account_password(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBInstance_vpc_security_ip_list(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
			{
				Config: testMongoDBInstance_vpc_backup(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBInstance_vpc_together(password),
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

func TestAccAlibabacloudStackMongoDBInstance_multiAZ(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	password := GeneratePassword()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_multiAZ_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
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

func TestAccAlibabacloudStackMongoDBInstance_multi_instance(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_instance.default.2"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	password := GeneratePassword()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBInstance_multi_instance_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version":       "3.4",
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
				Config: testMongoDBInstance_multi_instance_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multi_instance_configure,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_storage": "30",
						"db_instance_class":   "dds.mongo.standard",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multi_instance_account_password(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multi_instance_security_ip_list(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "1",
						"security_ip_list.0": "10.168.1.12",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multi_instance_backup(password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBInstance_multi_instance_together(password),
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

const testMongoDBInstance_classic_base = `

data "alibabacloudstack_zones" "default" { 
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
}`

const testMongoDBInstance_classic_ssl_action = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  ssl_action          = "Open"
}`

const testMongoDBInstance_classic_ssl_action_update = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  ssl_action          = "Close"
}`

const testMongoDBInstance_classic_base4 = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
}`

const testMongoDBInstance_classic_tags = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  ssl_action          = "Close"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}`

const testMongoDBInstance_classic_name = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  name                = "tf-testAccMongoDBInstance_test"
  ssl_action          = "Close"
}`

const testMongoDBInstance_classic_configure = `

data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  ssl_action          = "Close"
}`

func testMongoDBInstance_classic_account_password(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  ssl_action          = "Close"
}`, password)
}

const testMongoDBInstance_classic_tde = `
data "alibabacloudstack_zones" "default" {
  
}

resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "4.0"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  tde_status    = "enabled"
}`

func testMongoDBInstance_classic_security_ip_list(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  ssl_action          = "Close"
}`, password)
}

func testMongoDBInstance_classic_security_group_id(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

data "alibabacloudstack_security_groups" "default" {
}
resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_group_id    = "${data.alibabacloudstack_security_groups.default.groups.0.id}"
  ssl_action          = "Close"
}`, password)
}

func testMongoDBInstance_classic_backup(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  backup_period       = ["Wednesday"]
  backup_time         = "11:00Z-12:00Z"
  ssl_action          = "Close"
}`, password)
}

func testMongoDBInstance_classic_maintain_time(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  backup_period       = ["Wednesday"]
  backup_time         = "11:00Z-12:00Z"
  maintain_start_time = "02:00Z"
  maintain_end_time   = "03:00Z"
  ssl_action          = "Close"
}`, password)
}

func testMongoDBInstance_classic_together(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

resource "alibabacloudstack_mongodb_instance" "default" {
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test_together"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12", "10.168.1.13"]
  backup_period       = ["Tuesday", "Wednesday"]
  backup_time         = "10:00Z-11:00Z"
  ssl_action          = "Open"
}`, password)
}

const testMongoDBInstance_vpc_base = `
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
}`

const testMongoDBInstance_vpc_name = `
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  name                = "tf-testAccMongoDBInstance_test"
}`

const testMongoDBInstance_vpc_configure = `
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
}`

func testMongoDBInstance_vpc_account_password(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
}`, password)
}

func testMongoDBInstance_vpc_security_ip_list(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
}`, password)
}

func testMongoDBInstance_vpc_backup(password string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default" {
  
}

variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  backup_period       = ["Wednesday"]
  backup_time         = "11:00Z-12:00Z"
}`, password)
}

func testMongoDBInstance_vpc_together(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  
}
variable "name" {
  default = "tf-testAccMongoDBInstance_vpc"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test_together"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12", "10.168.1.13"]
  backup_period       = ["Tuesday", "Wednesday"]
  backup_time         = "10:00Z-11:00Z"
}`, password)
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
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
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test_together"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12", "10.168.1.13"]
  backup_period       = ["Tuesday", "Wednesday"]
  backup_time         = "10:00Z-11:00Z"
}`, password)
}

const testMongoDBInstance_multi_instance_base = `

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
}`

const testMongoDBInstance_multi_instance_name = `

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 10
  db_instance_class   = "dds.mongo.mid"
  name                = "tf-testAccMongoDBInstance_test"
}`

const testMongoDBInstance_multi_instance_configure = `

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
}`

func testMongoDBInstance_multi_instance_account_password(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
}`, password)
}

func testMongoDBInstance_multi_instance_security_ip_list(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
}`, password)
}

func testMongoDBInstance_multi_instance_backup(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12"]
  backup_period       = ["Wednesday"]
  backup_time         = "11:00Z-12:00Z"
}`, password)
}

func testMongoDBInstance_multi_instance_together(password string) string {
	return fmt.Sprintf(`

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBInstance_multi_instance"
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
  count               = 3
  engine_version      = "3.4"
  db_instance_storage = 30
  db_instance_class   = "dds.mongo.standard"
  name                = "tf-testAccMongoDBInstance_test_together"
  account_password    = "%s"
  security_ip_list    = ["10.168.1.12", "10.168.1.13"]
  backup_period       = ["Tuesday", "Wednesday"]
  backup_time         = "10:00Z-11:00Z"
}`, password)
}
