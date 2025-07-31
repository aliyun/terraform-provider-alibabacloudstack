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
	resource.AddTestSweepers("alibabacloudstack_mongodb_sharding_instance", &resource.Sweeper{
		Name: "alibabacloudstack_mongodb_sharding_instance",
		F:    testSweepMongoDBShardingInstances,
	})
}
func testSweepMongoDBShardingInstances(region string) error {
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
	request.DBInstanceType = "sharding"
	for {
		raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
			return ddsClient.DescribeDBInstances(request)
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "testSweepMongoDBShardingInstances", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
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
				if errmsgs.NotFoundError(err) {
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
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}
		return err
	}
	return nil
}

func TestAccAlibabacloudStackMongoDBShardingInstance_basicv4(t *testing.T) {
	var v dds.DBInstance
	resourceId := "alibabacloudstack_mongodb_sharding_instance.default"
	serverFunc := func() interface{} {
		return &MongoDBService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tfaccount%d", rand)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serverFunc, "DescribeMongoDBInstance")
	ra := resourceAttrInit(resourceId, nil)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testMongoDBShardingInstance_base(true))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		// CheckDestroy:  testAccCheckMongoDBShardingInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{

					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version": "4.0",
					"shard_list": []map[string]interface{}{
						{
							"description":  "shard1",
							"node_class":   "dds.shard.mid",
							"node_storage": 10,
						},
						{
							"description":  "shard2",
							"node_class":   "dds.shard.mid",
							"node_storage": 10,
						},
					},
					"mongo_list": []map[string]interface{}{
						{
							"description": "mongo1",
							"node_class":  "dds.mongos.mid",
						},
						{
							"description": "mongo2",
							"node_class":  "dds.mongos.mid",
						},
					},
					"configserver_list": []map[string]interface{}{
						{
							"description":  "cs1",
							"node_class":   "dds.cs.mid",
							"node_storage": 20,
						},
					},
					"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
					//					"security_group_id": "${alibabacloudstack_ecs_securitygroup.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":        CHECKSET,
						"engine_version": "4.0",
						"shard_list.#":   "2",
						"mongo_list.#":   "2",
						"name":           "",
						"storage_engine": "WiredTiger",
						//						"instance_charge_type": "PostPaid",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configserver_list.*.account_name", "configserver_list.*.account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mongo_list": []map[string]interface{}{
						{
							"description": "mongo1",
							"node_class":  "dds.mongos.mid",
						},
						{
							"description": "mongo3",
							"node_class":  "dds.mongos.mid",
						},
					},
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mongo_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_list": []map[string]interface{}{
						{
							"description":  "shard1",
							"node_class":   "dds.shard.mid",
							"node_storage": "10",
						},
						{
							"description":  "shard3",
							"node_class":   "dds.shard.mid",
							"node_storage": "20",
						},
						{
							"description":  "shard4",
							"node_class":   "dds.shard.mid",
							"node_storage": "30",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configserver_list": []map[string]interface{}{
						{
							"description":  "cs1",
							"node_class":   "dds.cs.mid",
							"node_storage": 30,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configserver_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"audit_status": "Enable",
					"audit_filter": []map[string]interface{}{{
						"role_type": "db",
						"filters":   []string{"update", "delete", "command"},
					}, {
						"role_type": "mongos",
						"filters":   []string{"admin", "insert"},
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"audit_status":   "Enable",
						"audit_filter.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"audit_filter.*",
						map[string]string{
							"role_type": "db",
							"filters.#": "3",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"audit_filter.*",
						map[string]string{
							"role_type": "mongos",
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
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"admin",
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"command",
					),
					resource.TestCheckTypeSetElemAttr(
						resourceId,
						"audit_filter.*.filters.*",
						"insert",
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "${random_password.password.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
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
					"name":             "${var.name}_update",
					"backup_period":    []string{"Tuesday", "Wednesday"},
					"backup_time":      "10:00Z-11:00Z",
					"security_ip_list": []string{"10.168.1.12", "10.168.1.13"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               fmt.Sprintf("%s_update", name),
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
				Config: testAccConfig(map[string]interface{}{
					"tde_status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
		},
	})
}

func testMongoDBShardingInstance_base(enableVpc bool) func(string) string {
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
	
%s

%s

`, name, RandomPasswordTestCase(12), vpcString)
	}
}
