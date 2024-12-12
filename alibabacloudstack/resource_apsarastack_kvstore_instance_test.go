package alibabacloudstack

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var redisInstanceConnectionDomainRegexp = "^r-[a-z0-9]+.redis[.a-z-0-9]*.rds.aliyuncs.com"
var redisInstanceClassForTest = "redis.master.small.default"
var redisInstanceClassForTestUpdateClass = "redis.master.mid.default"
var memcacheInstanceConnectionDomainRegexp = "^m-[a-z0-9]+.memcache[.a-z-0-9]*.rds.aliyuncs.com"
var memcacheInstanceClassForTest = "memcache.master.small.default"
var memcacheInstanceClassForTestUpdateClass = "memcache.master.mid.default"

func init() {
	resource.AddTestSweepers("alibabacloudstack_kvstore_instance", &resource.Sweeper{
		Name: "alibabacloudstack_kvstore_instance",
		F:    testSweepKVStoreInstances,
	})
}

func testSweepKVStoreInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AlibabacloudStack client: %s", err)
	}
	client := rawClient.(*connectivity.AlibabacloudStackClient)

	prefixes := []string{
		"tf-testAcc",
		"testAcc",
	}

	var insts []r_kvstore.KVStoreInstanceInDescribeInstances
	req := r_kvstore.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		req.InstanceType = instanceType
		for {
			raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.DescribeInstances(req)
			})
			if err != nil {
				return fmt.Errorf("Error retrieving KVStore Instances: %s", err)
			}
			resp, _ := raw.(*r_kvstore.DescribeInstancesResponse)
			if resp == nil || len(resp.Instances.KVStoreInstance) < 1 {
				break
			}
			insts = append(insts, resp.Instances.KVStoreInstance...)

			if len(resp.Instances.KVStoreInstance) < PageSizeLarge {
				break
			}

			page, err := getNextpageNumber(req.PageNumber)
			if err != nil {
				return err
			}
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping KVStore Instance: %s (%s)", name, id)
			continue
		}

		sweeped = true
		log.Printf("[INFO] Deleting KVStore Instance: %s (%s)", name, id)
		req := r_kvstore.CreateDeleteInstanceRequest()
		req.InstanceId = id
		_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete KVStore Instance (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these KVStore instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAlibabacloudStackKVStoreRedisInstance_classictest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alibabacloudstack_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, KVStoreInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "cpu_type"},
			},
		},
	})
}

func TestAccAlibabacloudStackKVStoreRedisInstance_vpctest(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alibabacloudstack_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, KVStoreInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAlibabacloudStackKVStoreRedisInstance_vpcmulti(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alibabacloudstack_kvstore_instance.default.2"
	ra := resourceAttrInit(resourceId, KVStoreInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpcmulti(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackKVStoreRedisInstance_classicmulti(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alibabacloudstack_kvstore_instance.default.2"
	ra := resourceAttrInit(resourceId, KVStoreInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_classicmulti(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore2Dot8)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckKVStoreInstanceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	kvstoreService := KvstoreService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_kvstore_instance" {
			continue
		}

		_, err := kvstoreService.DescribeKVstoreInstance(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func testAccKVStoreInstance_classic(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`

	
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstancesDataSource4"
}

%s

resource "alibabacloudstack_kvstore_instance" "default" {
	instance_name  = var.name
	vswitch_id     = alibabacloudstack_vswitch.default.id
	private_ip     = "172.16.0.10"
	security_ips   = ["10.0.0.1"]
	instance_type  = "%s"
	instance_class = "%s"
	engine_version = "%s"
    cpu_type = "intel"
}

	`, VSwichCommonTestCase, instanceType, instanceClass, engineVersion)
}

var KVStoreInstanceCheckMap = map[string]string{
	"instance_name":  CHECKSET,
	"instance_class": CHECKSET,
}

func testAccKVStoreInstance_classicUpdateParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters {
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			}
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicAddParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters {
				name = "maxmemory-policy"
				value = "volatile-ttl"
			  }
		parameters {
				  name = "slowlog-max-len"
				  value = "1111"
			  }
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicDeleteParameter(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		parameters {
				name = "slowlog-max-len"
				value = "1111"
			}
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_classicUpdateSecuirtyIps(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateClass(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateAttr(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		password = "inputYourCodeHere"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateTags(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		password = "inputYourCodeHere"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		tags = {
			Created = "TF"
			For		= "acceptance test"
		}
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateMaintainStartTime(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	data "alibabacloudstack_zones" "default" {
		available_resource_creation = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		availability_zone = "${lookup(data.alibabacloudstack_zones.default.zones[(length(data.alibabacloudstack_zones.default.zones)-1)%%length(data.alibabacloudstack_zones.default.zones)], "id")}"
		password = "inputYourCodeHere"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
		maintain_start_time = "02:00Z"
		maintain_end_time = "03:00Z"
		tags = {
			Created = "TF"
			For		= "acceptance test"
		}
	}
	`, instanceType, instanceClass, engineVersion)
}
func testAccKVStoreInstance_classicUpdateAll(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`
	
	variable "name" {
		default = "tf-testAccKVStoreInstance_classicUpdateAll"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		password = "inputYourCodeHere"
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.2","10.0.0.3"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
	}
	`, instanceType, instanceClass, engineVersion)
}

func testAccKVStoreInstance_vpc(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s

	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		engine_version = "%s"
cpu_type = "intel"
       
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreInstance_vpcUpdateSecurityIps(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpcc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateSecurityGroupIds(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	data "alibabacloudstack_security_groups" "default" {
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
		security_group_id    = "${data.alibabacloudstack_security_groups.default.groups.0.id}"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateVpcAuthMode(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		vpc_auth_mode = "Close"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters {
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			}
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcAddParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters {
			  name = "maxmemory-policy"
			  value = "volatile-ttl"
			}
        parameters {
				name = "slowlog-max-len"
				value = "1111"
			}
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcDeleteParameter(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		parameters {
				name = "slowlog-max-len"
				value = "1111"
			}
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcUpdateClass(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.3", "10.0.0.2"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}
func testAccKVStoreInstance_vpcUpdateAll(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstancevpcUpdateAlll"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		instance_class = "%s"
		instance_name  = "${var.name}"
		password       = "inputYourCodeHere"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		private_ip     = "172.16.0.10"
		security_ips = ["10.0.0.1", "10.0.0.4"]
		instance_type = "%s"
		engine_version = "%s"
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_vpcmulti(common, instanceClass, instanceType, engineVersion string) string {
	return fmt.Sprintf(`
	%s

	variable "creation" {
		default = "KVStore"
	}
	variable "name" {
		default = "tf-testAccKVStoreInstance_vpc"
	}
	resource "alibabacloudstack_kvstore_instance" "default" {
		count		   = 3
		instance_class = "%s"
		instance_name  = "${var.name}"
		password       = "inputYourCodeHere"
		vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
		security_ips   = ["10.0.0.1"]
		instance_type  = "%s"
		engine_version = "%s"
cpu_type = "intel"
       
	}
	`, common, instanceClass, instanceType, engineVersion)
}

func testAccKVStoreInstance_classicmulti(instanceType, instanceClass, engineVersion string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "tf-testAccKVStoreInstance_classic"
	}

	resource "alibabacloudstack_kvstore_instance" "default" {
		count = 3
		instance_name  = "${var.name}"
		security_ips = ["10.0.0.1"]
		instance_type = "%s"
		instance_class = "%s"
		engine_version = "%s"
cpu_type = "intel"
       
	}
	`, instanceType, instanceClass, engineVersion)
}
