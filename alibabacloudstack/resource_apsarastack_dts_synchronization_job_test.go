package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDTSSynchronizationJob_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_dts_synchronization_job.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackDTSSynchronizationJobMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeDtsSynchronizationJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtssynchronizationjob%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackDTSSynchronizationJobBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_instance_id":                    "${alibabacloudstack_dts_synchronization_instance.default.id}",
					"dts_job_name":                       "tf-testAccCase",
					"source_endpoint_instance_type":      "RDS",
					"source_endpoint_instance_id":        "${alibabacloudstack_db_instance.rsinstance.id}",
					"source_endpoint_engine_name":        "MySQL",
					"source_endpoint_database_name":      "tfaccountpri_0",
					"source_endpoint_user_name":          "tftestdts",
					"source_endpoint_password":           "inputYourCodeHere",
					"destination_endpoint_instance_type": "RDS",
					"destination_endpoint_instance_id":   "${alibabacloudstack_db_instance.dsinstance.id}",
					"destination_endpoint_engine_name":   "MySQL",
					"destination_endpoint_database_name": "tfaccountpri_0",
					"destination_endpoint_user_name":     "tftestdts",
					"destination_endpoint_password":      "inputYourCodeHere",
					"db_list":                            "{\\\"tfaccountpri_0\\\":{\\\"name\\\":\\\"tfaccountpri_0\\\",\\\"all\\\":true,\\\"state\\\":\\\"normal\\\"}}",
					"structure_initialization":           "true",
					"data_initialization":                "true",
					"data_synchronization":               "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_name":                       "tf-testAccCase",
						"source_endpoint_instance_type":      "RDS",
						"source_endpoint_engine_name":        "MySQL",
						"source_endpoint_database_name":      "tfaccountpri_0",
						"source_endpoint_user_name":          "tftestdts",
						"source_endpoint_password":           "inputYourCodeHere",
						"destination_endpoint_instance_type": "RDS",
						"destination_endpoint_engine_name":   "MySQL",
						"destination_endpoint_database_name": "tfaccountpri_0",
						"destination_endpoint_user_name":     "tftestdts",
						"destination_endpoint_password":      "inputYourCodeHere",
						"db_list":                            "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}",
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
					"destination_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Suspending",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Suspending",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_endpoint_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_endpoint_password": "inputYourCodeHere",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Synchronizing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Synchronizing",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"delay_notice", "error_phone", "delay_rule_time", "error_notice", "delay_phone", "reserve", "destination_endpoint_password", "source_endpoint_password"},
			},
		},
	})
}

var AlibabacloudStackDTSSynchronizationJobMap0 = map[string]string{
	"error_phone":                      NOSET,
	"error_notice":                     NOSET,
	"delay_rule_time":                  NOSET,
	"delay_phone":                      NOSET,
	"source_endpoint_engine_name":      CHECKSET,
	"reserve":                          NOSET,
	"delay_notice":                     NOSET,
	"destination_endpoint_engine_name": CHECKSET,
	"status":                           CHECKSET,
}

func AlibabacloudStackDTSSynchronizationJobBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "creation" {
  default = "Rds"
}
data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}
resource "alibabacloudstack_vpc" "default" {
  vpc_name       = var.name
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone  = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name      = var.name
}
resource "alibabacloudstack_db_instance" "dsinstance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  vswitch_id       = alibabacloudstack_vswitch.default.id
  instance_name    = var.name
  storage_type         = "local_ssd"
}
resource "alibabacloudstack_db_instance" "rsinstance" {
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
  instance_id = alibabacloudstack_db_instance.dsinstance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
  character_set =  "UTF8"
}

resource "alibabacloudstack_db_account" "account" {
  instance_id      = alibabacloudstack_db_instance.dsinstance.id
  name        = "tftestdts"
  password    = "inputYourCodeHere"
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege" {
  instance_id  = alibabacloudstack_db_instance.dsinstance.id
  account_name = alibabacloudstack_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alibabacloudstack_db_database.db.*.name
}

resource "alibabacloudstack_db_database" "db_r" {
  count       = 2
  instance_id = alibabacloudstack_db_instance.rsinstance.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
character_set =  "UTF8"
}

resource "alibabacloudstack_db_account" "account_r" {
  instance_id      =alibabacloudstack_db_instance.rsinstance.id
  name        = "tftestdts"
  password    = "inputYourCodeHere"
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege_r" {
  instance_id  = alibabacloudstack_db_instance.rsinstance.id
  account_name = alibabacloudstack_db_account.account_r.name
  privilege    = "ReadWrite"
  db_names     = alibabacloudstack_db_database.db_r.*.name
}

resource "alibabacloudstack_dts_synchronization_instance" "default" {
  payment_type                        = "PostPaid"
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = "cn-qingdao-env17-d01"
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = "cn-qingdao-env17-d01"
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}

`, name)
}
