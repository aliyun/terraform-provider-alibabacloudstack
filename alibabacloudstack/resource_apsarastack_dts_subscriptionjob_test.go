package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alibabacloudstack_dts_subscription_job",
		&resource.Sweeper{
			Name: "alibabacloudstack_dts_subscription_job",
			F:    testSweepDtsSubscriptionJob,
		})
}

func testSweepDtsSubscriptionJob(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return errmsgs.WrapErrorf(err, "Error getting AlibabacloudStack client.")
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDtsJobs"
	request := make(map[string]interface{})
	request["JobType"] = "SUBSCRIBE"
	request["Region"] = region
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return errmsgs.WrapError(err)
	}

	for {
		runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if errmsgs.NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch Dts SubscriptionJobs: %s", errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_dts_subscription_jobs", action, errmsgs.AlibabacloudStackSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			log.Printf("[ERROR] Failed to parse Dts SubscriptionJobs: %s", errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.DtsJobList", response))
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DtsJobName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Dts SubscriptionJob: %s", item["DtsJobName"].(string))
				continue
			}

			action := "DeleteDtsJob"
			request := map[string]interface{}{
				"DtsJobId": item["DtsJobId"],
			}
			request["DtsInstanceId"] = item["DtsInstanceID"]
			request["RegionId"] = client.RegionId
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dts SubscriptionJob (%s): %s", item["DtsJobName"].(string), err)
			}
			log.Printf("[INFO] Delete Dts SubscriptionJob success: %s ", item["DtsJobName"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

func TestAccAlibabacloudStackDTSSubscriptionJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDTSSubscriptionJobBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase",
					"payment_type":                       "PayAsYouGo",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "cn-hangzhou",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alibabacloudstack_db_instance.instance.id}",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestprivilege",
					"source_endpoint_password":           "inputYourCodeHere",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${alibabacloudstack_vpc.default1.id}",
					"subscription_instance_vswitch_id":   "${alibabacloudstack_vswitch.default1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"payment_type":                       "PayAsYouGo",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             "cn-hangzhou",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestprivilege",
						"source_endpoint_password":           "inputYourCodeHere",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name": "tf-testAccCase1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccCase1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_list": "{\\\"dtstestdata\\\": {   \\\"name\\\": \\\"tfaccountpri_0\\\",   \\\"all\\\": true }}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_list": "{\"dtstestdata\": {   \"name\": \"tfaccountpri_0\",   \"all\": true }}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Abnormal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Abnormal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${alibabacloudstack_vpc.default2.id}",
					"subscription_instance_vswitch_id":   "${alibabacloudstack_vswitche.default2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase2",
					"source_endpoint_instance_id":        "${alibabacloudstack_db_instance.instance.id}",
					"subscription_instance_network_type": "vpc",
					"subscription_instance_vpc_id":       "${data.alibabacloudstack_vpcs.default1.ids[0]}",
					"subscription_instance_vswitch_id":   "${data.alibabacloudstack_vswitches.default_1.ids[0]}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase2",
						"subscription_instance_network_type": "vpc",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

func TestAccAlibabacloudStackDTSSubscriptionJob_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDTSSubscriptionJobBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase",
					"payment_type":                       "PayAsYouGo",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "cn-hangzhou",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alibabacloudstack_db_instance.instance.id}",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestprivilege",
					"source_endpoint_password":           "inputYourCodeHere",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "classic",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"payment_type":                       "PayAsYouGo",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             "cn-hangzhou",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestprivilege",
						"source_endpoint_password":           "inputYourCodeHere",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "classic",
						"tags.%":                             "2",
						"tags.Created":                       "TF",
						"tags.For":                           "acceptance test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "subscribeJob",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "subscribeJob",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name": "tf-testAccCase1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name": "tf-testAccCase1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Normal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Normal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_list": "{\\\"dtstestdata\\\": {   \\\"name\\\": \\\"tfaccountpri_0\\\",   \\\"all\\\": true }}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_list": "{\"dtstestdata\": {   \"name\": \"tfaccountpri_0\",   \"all\": true }}",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

func TestAccAlibabacloudStackDTSSubscriptionJob_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dts_subscription_job.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDTSSubscriptionJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDtsSubscriptionJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssubscriptionjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDTSSubscriptionJobBasicDependence0)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime2(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_name":                       "tf-testAccCase",
					"payment_type":                       "Subscription",
					"payment_duration_unit":              "Month",
					"payment_duration":                   "1",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_region":             "cn-hangzhou",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alibabacloudstack_db_instance.instance.id}",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestprivilege",
					"source_endpoint_password":           "inputYourCodeHere",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"subscription_instance_network_type": "classic",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"payment_type":                       "Subscription",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_region":             "cn-hangzhou",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestprivilege",
						"source_endpoint_password":           "inputYourCodeHere",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
						"subscription_instance_network_type": "classic",
						"tags.%":                             "2",
						"tags.Created":                       "TF",
						"tags.For":                           "acceptance test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"source_endpoint_password", "error_phone", "payment_duration_unit", "destination_region", "database_count", "delay_notice", "reserve", "synchronization_direction", "instance_class", "destination_endpoint_engine_name", "payment_duration", "delay_rule_time", "delay_phone", "compute_unit", "error_notice", "sync_architecture"},
			},
		},
	})
}

var AlibabacloudStackDTSSubscriptionJobMap0 = map[string]string{
	"auto_pay":                         NOSET,
	"delay_rule_time":                  NOSET,
	"compute_unit":                     NOSET,
	"delay_phone":                      NOSET,
	"subscription_data_type_dml":       CHECKSET,
	"error_notice":                     NOSET,
	"sync_architecture":                NOSET,
	"quantity":                         NOSET,
	"error_phone":                      NOSET,
	"period":                           NOSET,
	"destination_region":               NOSET,
	"delay_notice":                     NOSET,
	"reserve":                          NOSET,
	"synchronization_direction":        NOSET,
	"auto_start":                       NOSET,
	"database_count":                   NOSET,
	"instance_class":                   NOSET,
	"subscription_data_type_ddl":       CHECKSET,
	"destination_endpoint_engine_name": NOSET,
	"used_time":                        NOSET,
	"status":                           CHECKSET,
}

func AlibabacloudStackDTSSubscriptionJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "tf-testaccdts%s"
}

variable "creation" {
  default = "Rds"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = var.creation
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name       = var.name
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone           = data.alibabacloudstack_zones.default.zones[0].id
  name      = var.name
}

data "alibabacloudstack_db_zones" "default"{
	
}



resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  vswitch_id       = alibabacloudstack_vswitch.default.id
  instance_name    = var.name
  storage_type         = "local_ssd"
}

resource "alibabacloudstack_db_database" "db" {
  count       = 2
  instance_id = alibabacloudstack_db_instance.instance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
  character_set = "UTF8"
}

resource "alibabacloudstack_db_account" "account" {
  instance_id      = alibabacloudstack_db_instance.instance.id
  name        = "tftestprivilege"
  password    = "%s"
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege" {
  instance_id  = alibabacloudstack_db_instance.instance.id
  account_name = alibabacloudstack_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alibabacloudstack_db_database.db.*.name
}
resource "alibabacloudstack_vpc" "default1" {
  vpc_name       = var.name
  cidr_block     = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default1" {
  vpc_id            = alibabacloudstack_vpc.default1.id
  cidr_block        = "10.1.0.0/16"
  availability_zone           = data.alibabacloudstack_zones.default.zones[0].id
  name      = var.name
}
resource "alibabacloudstack_vpc" "default2" {
  vpc_name       = var.name
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default2" {
  vpc_id            = alibabacloudstack_vpc.default2.id
  cidr_block        = "172.16.0.0/24"
  availability_zone           = data.alibabacloudstack_zones.default.zones[0].id
  name      = var.name
}


`, name, GeneratePassword())
}
