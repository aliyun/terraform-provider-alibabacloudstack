package alibabacloudstack

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackKVStoreInstance_basic(t *testing.T) {
	var instance *r_kvstore.DBInstanceAttribute
	resourceId := "alibabacloudstack_kvstore_instance.default"
	ra := resourceAttrInit(resourceId, KVStoreInstanceCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeKVstoreInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccKVStore%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccKVStoreInstanceTdeclassicDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		Providers:         testAccProviders,
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":  "${var.name}",
					"vswitch_id":     "${alibabacloudstack_vpc_vswitch.default.id}",
					"instance_type":  "Redis",
					"cpu_type":       "intel",
					"instance_class": "${data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.id}",
					"engine_version": string(KVStore5Dot0),
					"tde_status":     "true",
					"enable_ssl":     "true",
					"node_type":      "double",
					"security_ips":   []string{"10.168.1.11", "10.168.1.12"},
					"password":       "${random_password.password.0.result}",
					"encryption_key": "${alibabacloudstack_kms_key.key.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name":  name,
						"tde_status":     "true",
						"enable_ssl":     "true",
						"node_type":      "double",
						"security_ips.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "${random_password.password.1.result}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${var.name}_updated",
					"vpc_auth_mode": "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_updated",
						"vpc_auth_mode": "Close",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "${data.alibabacloudstack_kvstore_instance_classes.update.instance_classes.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": CHECKSET,
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
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"name":  "maxmemory-policy",
							"value": "volatile-ttl",
						},
						{
							"name":  "slowlog-max-len",
							"value": "1111",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// password sensitive field does not echo after setting
				ImportStateVerifyIgnore: []string{"password", "cpu_type", "encryption_key", "parameters"},
			},
		},
	})
}

var KVStoreInstanceCheckMap = map[string]string{
	"instance_name":  CHECKSET,
	"instance_class": CHECKSET,
}

func testAccKVStoreInstanceTdeclassicDependence(name string) string {
	return fmt.Sprintf(`
	
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_kms_key" "key" {
  description             = "Hello KMS"
  pending_window_in_days  = "7"
  key_state               = "Enabled"
}

%s

%s

data "alibabacloudstack_kvstore_instance_classes" "default" {
  edition_type   = "enterprise"
  engine_version = "5.0"
  engine 		 = "Redis"
  sorted_by      = "Memory"
  architecture   = "cluster"
  memory         = 2
}

data "alibabacloudstack_kvstore_instance_classes" "update" {
  edition_type   = "enterprise"
  engine_version = "5.0"
  engine 		 = "Redis"
  sorted_by      = "Memory"
  architecture   = "cluster"
  memory         = 4
}
`, name, VSwitchCommonTestCase, RandomPasswordTestCase(12, 2))
}
