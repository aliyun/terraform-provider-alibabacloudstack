package apsarastack

import (
	"fmt"
	"testing"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackKVStoreAccountUpdateV4(t *testing.T) {
	var v *r_kvstore.Account
	rand := acctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccKVstoreAccount-%d", rand)
	var basicMap = map[string]string{
		"instance_id":      CHECKSET,
		"account_name":     "tftestnormal",
		"account_password": "inputYourCodeHere",
		"account_type":     "Normal",
	}
	resourceId := "apsarastack_kvstore_account.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeKVstoreAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKVstoreAccountConfigDependenceV4)
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
					"instance_id":      "${apsarastack_kvstore_instance.instance.id}",
					"account_name":     "tftestnormal",
					"account_password": "inputYourCodeHere",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_privilege": "RoleRepl",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_privilege": "RoleRepl",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":       "tf test",
					"account_password":  "inputYourCodeHere",
					"account_privilege": "RoleReadOnly",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":       "tf test",
						"account_password":  "inputYourCodeHere",
						"account_privilege": "RoleReadOnly",
					}),
				),
			},
		},
	})
}

func resourceKVstoreAccountConfigDependenceV4(name string) string {
	return fmt.Sprintf(`
provider "apsarastack" {
	assume_role {}
}
	data "apsarastack_zones" "default" {
	}
	variable "name" {
		default = "%v"
	}
	resource "apsarastack_kvstore_instance" "instance" {
		availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
		instance_class = "redis.master.small.default"
		instance_name  = "${var.name}"
		instance_charge_type = "PostPaid"
		engine_version = "4.0"
	}
	`, name)
}

//func resourceKVstoreAccountConfigDependenceV5(name string) string {
//	return fmt.Sprintf(`
//	data "apsarastack_zones" "default" {
//		available_resource_creation = "KVStore"
//	}
//	variable "name" {
//		default = "%v"
//	}
//	resource "apsarastack_kvstore_instance" "instance" {
//		availability_zone = "${lookup(data.apsarastack_zones.default.zones[(length(data.apsarastack_zones.default.zones)-1)%%length(data.apsarastack_zones.default.zones)], "id")}"
//		instance_class = "redis.master.small.default"
//		instance_name  = "${var.name}"
//		instance_charge_type = "PostPaid"
//		engine_version = "5.0"
//	}
//	`, name)
//}
