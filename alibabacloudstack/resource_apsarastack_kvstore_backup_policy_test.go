package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackKVStoreRedisBackupPolicy_classic(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alibabacloudstack_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_classic(rand, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(rand, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateTime(rand, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "12:00Z-13:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_classicUpdateAll(rand, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "13:00Z-14:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

/*func TestAccAlibabacloudStackKVStoreMemcacheBackupPolicy_classic(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alibabacloudstack_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			//{
			//	Config: testAccKVStoreBackupPolicy_classicUpdatePeriod(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"backup_period.#": "3",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccKVStoreBackupPolicy_classicUpdateTime(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"backup_time": "12:00Z-13:00Z",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccKVStoreBackupPolicy_classicUpdateAll(string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore4Dot0)),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"backup_time":     "13:00Z-14:00Z",
			//			"backup_period.#": "1",
			//		}),
			//	),
			//},
		},
	})

}*/

func TestAccAlibabacloudStackKVStoreRedisBackupPolicy_vpc(t *testing.T) {
	var policy *r_kvstore.DescribeBackupPolicyResponse

	resourceId := "alibabacloudstack_kvstore_backup_policy.default"
	ra := resourceAttrInit(resourceId, kvStoreMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 99999)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_vpc(rand, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(rand, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateTime(rand, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time": "11:00Z-12:00Z",
					}),
				),
			},
			{
				Config: testAccKVStoreBackupPolicy_vpcUpdateAll(rand, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_time":     "12:00Z-13:00Z",
						"backup_period.#": "1",
					}),
				),
			},
		},
	})

}

//	func TestAccAlibabacloudStackKVStoreMemcacheBackupPolicy_vpc(t *testing.T) {
//		var policy *r_kvstore.DescribeBackupPolicyResponse
//		resourceId := "alibabacloudstack_kvstore_backup_policy.default"
//		ra := resourceAttrInit(resourceId, kvStoreMap)
//		serviceFunc := func() interface{} {
//			return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
//		}
//		rc := resourceCheckInitWithDescribeMethod(resourceId, &policy, serviceFunc, "DescribeKVstoreBackupPolicy")
//		rac := resourceAttrCheckInit(rc, ra)
//		testAccCheck := rac.resourceAttrMapUpdateSet()
//		ResourceTest(t, resource.TestCase{
//			PreCheck: func() {
//				testAccPreCheck(t)
//			},
//			// module name
//			IDRefreshName: resourceId,
//			Providers:     testAccProviders,
//			CheckDestroy:  testAccCheckKVStoreBackupPolicyDestroy,
//			Steps: []resource.TestStep{
//				{
//				     Config: testAccKVStoreBackupPolicy_vpc(VSwitchCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
//				     Check: resource.ComposeTestCheckFunc(
//				        testAccCheck(nil),
//				     ),
//				  },
//				  {
//				     ResourceName:      resourceId,
//				     ImportState:       true,
//				     ImportStateVerify: true,
//				  },
//				  {
//				     Config: testAccKVStoreBackupPolicy_vpcUpdatePeriod(VSwitchCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
//				     Check: resource.ComposeTestCheckFunc(
//				        testAccCheck(map[string]string{
//				           "backup_period.#": "3",
//				        }),
//				     ),
//				  },
//				  {
//				     Config: testAccKVStoreBackupPolicy_vpcUpdateTime(VSwitchCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
//				     Check: resource.ComposeTestCheckFunc(
//				        testAccCheck(map[string]string{
//				           "backup_time": "11:00Z-12:00Z",
//				        }),
//				     ),
//				  },
//				{
//					Config: testAccKVStoreBackupPolicy_vpcUpdateAll(VSwitchCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
//					Check: resource.ComposeTestCheckFunc(
//						testAccCheck(map[string]string{
//							"backup_time":     "12:00Z-13:00Z",
//							"backup_period.#": "1",
//						}),
//					),
//				},
//			},
//		})
//	}
func testAccCheckKVStoreBackupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_kvstore_instance" {
			continue
		}

		if _, err := kvstoreService.DescribeKVstoreBackupPolicy(rs.Primary.ID); err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error Describe DB backup policy: %#v", err)
		}
		return fmt.Errorf("KVStore Instance %s Policy sitll exists.", rs.Primary.ID)
	}

	return nil
}

var kvStoreMap = map[string]string{
	"instance_id":     CHECKSET,
	"backup_time":     "10:00Z-11:00Z",
	"backup_period.#": "2",
}

func testAccKVStoreBackupPolicy_classic(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`

	data "alibabacloudstack_zones" "default" {
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic%d"
}

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "%s"
}

%s 

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "%s"
	node_type = "double"
	architecture_type = "standard"
	password       = "1qaz@WSX"
}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, rand, instanceType, KVRInstanceClassCommonTestCase, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdatePeriod(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`

	data "alibabacloudstack_zones" "default" {
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic%d"
}

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "%s"
}

%s 

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "%s"
	node_type = "double"
	architecture_type = "standard"
	password       = "1qaz@WSX"
}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, rand, instanceType, KVRInstanceClassCommonTestCase, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateTime(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`

	data "alibabacloudstack_zones" "default" {
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic%d"
}

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "%s"
}

%s 

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "%s"
	node_type = "double"
	architecture_type = "standard"
	password       = "1qaz@WSX"
}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, rand, instanceType, KVRInstanceClassCommonTestCase, engineVersion)
}

func testAccKVStoreBackupPolicy_classicUpdateAll(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`

	data "alibabacloudstack_zones" "default" {
	}
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_classic%d"
}

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "%s"
}

%s 

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "%s"
	node_type = "double"
	architecture_type = "standard"
	password       = "1qaz@WSX"
}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Sunday"]
		backup_time = "13:00Z-14:00Z"
	}
	`, rand, instanceType, KVRInstanceClassCommonTestCase, engineVersion)
}

func testAccKVStoreBackupPolicy_vpc(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc%d"
	}
	variable "kv_edition" {
    default = "enterprise"
	}
	
	variable "kv_engine" {
    default = "%s"
	}

	%s 

	resource "alibabacloudstack_kvstore_instance" "default" {
		zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
		instance_class = local.default_kv_instance_classes
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
		security_ips = ["10.0.0.1"]
		instance_type = var.kv_engine
		engine_version = "%s"
	}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, rand, instanceType, VSwitchCommonTestCase+KVRInstanceClassCommonTestCase, engineVersion)
}

func testAccKVStoreBackupPolicy_vpcUpdatePeriod(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc%d"
	}
	variable "kv_edition" {
    default = "enterprise"
	}
	
	variable "kv_engine" {
    default = "%s"
	}

	%s 

	resource "alibabacloudstack_kvstore_instance" "default" {
		zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
		instance_class = local.default_kv_instance_classes
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
		security_ips = ["10.0.0.1"]
		instance_type = var.kv_engine
		engine_version = "%s"
	}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "10:00Z-11:00Z"
	}
	`, rand, instanceType, VSwitchCommonTestCase+KVRInstanceClassCommonTestCase, engineVersion)
}
func testAccKVStoreBackupPolicy_vpcUpdateTime(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc%d"
	}
	variable "kv_edition" {
    default = "enterprise"
	}
	
	variable "kv_engine" {
    default = "%s"
	}

	%s 

	resource "alibabacloudstack_kvstore_instance" "default" {
		zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
		instance_class = local.default_kv_instance_classes
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
		security_ips = ["10.0.0.1"]
		instance_type = var.kv_engine
		engine_version = "%s"
	}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday", "Wednesday", "Sunday"]
		backup_time = "11:00Z-12:00Z"
	}
	`, rand, instanceType, VSwitchCommonTestCase+KVRInstanceClassCommonTestCase, engineVersion)
}
func testAccKVStoreBackupPolicy_vpcUpdateAll(rand int, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccKVStoreBackupPolicy_vpc%d"
	}
	variable "kv_edition" {
    default = "enterprise"
	}
	
	variable "kv_engine" {
    default = "%s"
	}

	%s 

	resource "alibabacloudstack_kvstore_instance" "default" {
		zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
		instance_class = local.default_kv_instance_classes
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
		security_ips = ["10.0.0.1"]
		instance_type = var.kv_engine
		engine_version = "%s"
	}
	resource "alibabacloudstack_kvstore_backup_policy" "default" {
		instance_id = "${alibabacloudstack_kvstore_instance.default.id}"
		backup_period = ["Tuesday"]
		backup_time = "12:00Z-13:00Z"
	}
	`, rand, instanceType, VSwitchCommonTestCase+KVRInstanceClassCommonTestCase, engineVersion)
}
