package alibabacloudstack

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("alibabacloudstack_mongodb_sharding_instance", &resource.Sweeper{
		Name: "alibabacloudstack_mongodb_sharding_instance",
		F:    testSweepMongoDBShardingInstances,
	})
}
func testSweepMongoDBShardingInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
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
	request.DBInstanceType = "sharding"
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "testSweepMongoDBShardingInstances", request.GetActionName(), AlibabacloudStackSdkGoERROR)
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
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	sweeped := false
	service := VpcService{client}
	ddsService := MongoDBService{client}
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
		// If a mongoDB name is not set successfully, it should be fetched by vpc name and deleted.
		if skip {
			instance, err := ddsService.DescribeMongoDBInstance(id)
			if err != nil {
				if NotFoundError(err) {
					continue
				}
				log.Printf("[INFO] Describe MongoDB sharding instance: %s (%s) got an error: %#v\n", name, id, err)
			}
			if need, err := service.needSweepVpc(instance.VPCId, instance.VSwitchId); err == nil {
				skip = !need
			}
		}
		if skip {
			log.Printf("[INFO] Skipping MongoDB sharding instance: %s (%s)\n", name, id)
			continue
		}
		log.Printf("[INFO] Deleting MongoDB sharding instance: %s (%s)\n", name, id)

		request := dds.CreateDeleteDBInstanceRequest()
		request.DBInstanceId = id
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DeleteDBInstance(request)
		})

		if err != nil {
			log.Printf("[error] Failed to delete MongoDB sharding instance,ID:%v(%v)\n", id, request.GetActionName())
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

func testAccCheckMongoDBShardingInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	ddsService := MongoDBService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_mongodb_sharding_instance" {
			continue
		}
		_, err := ddsService.DescribeMongoDBInstance(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return err
	}
	return nil
}

func TestAccAlibabacloudStackMongoDBShardingInstance_classic(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_classic_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testMongoDBShardingInstance_classic_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_classic_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":   "inputYourCodeHere",
						"security_ip_list.#": "2",
						"security_ip_list.0": "10.168.1.12",
						"security_ip_list.1": "10.168.1.13",
						"backup_period.#":    "2",
						"backup_period.0":    "Tuesday",
						"backup_period.1":    "Wednesday",
						"backup_time":        "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func TestAccAlibabacloudStackMongoDBShardingInstance_classicVersion4(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.MongoDBClassicNoSupportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_classic_base4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "4.0",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testMongoDBShardingInstance_classic_tde,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			}},
		//			{
		//				Config: testMongoDBShardingInstance_classic_security_group_id,
		//				Check: resource.ComposeTestCheckFunc(
		//					testAccCheck(map[string]string{
		//						"security_group_id": CHECKSET,
		//					}),
		//				),
		//			}},
	})
}

func TestAccAlibabacloudStackMongoDBShardingInstance_vpc(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithNoDefaultVpc(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_vpc_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_id":                CHECKSET,
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testMongoDBShardingInstance_vpc_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "1",
						"backup_period.0": "Wednesday",
						"backup_time":     "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_vpc_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":   "inputYourCodeHere",
						"security_ip_list.#": "2",
						"security_ip_list.0": "10.168.1.12",
						"security_ip_list.1": "10.168.1.13",
						"backup_period.#":    "2",
						"backup_period.0":    "Tuesday",
						"backup_period.1":    "Wednesday",
						"backup_time":        "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

func TestAccAlibabacloudStackMongoDBShardingInstance_multi_instance(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_sharding_instance.default.2"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_multi_instance_base,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":                   CHECKSET,
						"engine_version":            "3.4",
						"shard_list.#":              "2",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"mongo_list.#":              "2",
						"mongo_list.0.node_class":   "dds.mongos.mid",
						"mongo_list.1.node_class":   "dds.mongos.mid",
						"name":                      "",
						"storage_engine":            "WiredTiger",
						"instance_charge_type":      "PostPaid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccMongoDBShardingInstance_test",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_account_password,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_mongos,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#":            "3",
						"mongo_list.0.node_class": "dds.mongos.mid",
						"mongo_list.1.node_class": "dds.mongos.mid",
						"mongo_list.2.node_class": "dds.mongos.mid",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_shard,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#":              "3",
						"shard_list.0.node_class":   "dds.shard.mid",
						"shard_list.0.node_storage": "10",
						"shard_list.1.node_class":   "dds.shard.standard",
						"shard_list.1.node_storage": "20",
						"shard_list.2.node_class":   "dds.shard.standard",
						"shard_list.2.node_storage": "20",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_backup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":   "inputYourCodeHere",
						"security_ip_list.#": "2",
						"security_ip_list.0": "10.168.1.12",
						"security_ip_list.1": "10.168.1.13",
						"backup_period.#":    "2",
						"backup_period.0":    "Tuesday",
						"backup_period.1":    "Wednesday",
						"backup_time":        "10:00Z-11:00Z",
					}),
				),
			},
			{
				Config: testMongoDBShardingInstance_multi_instance_together,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               "tf-testAccMongoDBShardingInstance_test_together",
						"account_password":   "inputYourCodeHere",
						"security_ip_list.#": "2",
						"security_ip_list.0": "10.168.1.12",
						"security_ip_list.1": "10.168.1.13",
						"backup_period.#":    "2",
						"backup_period.0":    "Tuesday",
						"backup_period.1":    "Wednesday",
						"backup_time":        "10:00Z-11:00Z",
					}),
				),
			}},
	})
}

const testMongoDBShardingInstance_classic_base = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_classic_base4 = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_classic_tde = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  tde_status    = "enabled"
}`

const testMongoDBShardingInstance_classic_security_group_id = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
data "alibabacloudstack_security_groups" "default" {
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "4.0"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  tde_status    = "enabled"
  security_group_id    = "${data.alibabacloudstack_security_groups.default.groups.0.id}"
}`

const testMongoDBShardingInstance_classic_name = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_classic_account_password = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_classic_mongos = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_classic_shard = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_classic_backup = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_classic_together = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "inputYourCodeHere"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

const testMongoDBShardingInstance_vpc_base = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_vpc_name = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_vpc_account_password = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_vpc_mongos = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_vpc_shard = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_vpc_backup = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_vpc_together = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_vpc"
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

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  vswitch_id          = alibabacloudstack_vswitch.default.id
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "inputYourCodeHere"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`

const testMongoDBShardingInstance_multi_instance_base = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
}`

const testMongoDBShardingInstance_multi_instance_name = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name = "tf-testAccMongoDBShardingInstance_test"
}`

const testMongoDBShardingInstance_multi_instance_account_password = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_multi_instance_mongos = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_multi_instance_shard = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
}`

const testMongoDBShardingInstance_multi_instance_backup = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test"
  account_password = "inputYourCodeHere"
  backup_period    = ["Wednesday"]
  backup_time      = "11:00Z-12:00Z"
}`

const testMongoDBShardingInstance_multi_instance_together = `
provider "alibabacloudstack" {
	assume_role {}
}
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}
variable "name" {
  default = "tf-testAccMongoDBShardingInstance_multi_instance"
}

resource "alibabacloudstack_mongodb_sharding_instance" "default" {
  count          = 3
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "3.4"
  shard_list {
    node_class   = "dds.shard.mid"
    node_storage = 10
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
    }
  shard_list {
    node_class   = "dds.shard.standard"
    node_storage = 20
  }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
  mongo_list {
    node_class = "dds.mongos.mid"
    }
   mongo_list {
    node_class = "dds.mongos.mid"
  }
  name             = "tf-testAccMongoDBShardingInstance_test_together"
  account_password = "inputYourCodeHere"
  backup_period    = ["Tuesday", "Wednesday"]
  backup_time      = "10:00Z-11:00Z"
  security_ip_list = ["10.168.1.12", "10.168.1.13"]
}`
