package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRdsDbProxy_basic(t *testing.T) {
	var v *RdsDescribedbproxyResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testaccdbproxy%d", rand)
	resourceId := "alibabacloudstack_db_proxy.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"db_proxy_connect_string":                 CHECKSET,
		"db_proxy_connect_string_port":            CHECKSET,
		"db_proxy_instance_current_minor_version": CHECKSET,
		"db_proxy_instance_latest_minor_version":  CHECKSET,
		"db_proxy_instance_status":                CHECKSET,
		"db_proxy_endpoint_aliases":               CHECKSET,
		"db_proxy_endpoint_name":                  CHECKSET,
		"db_proxy_endpoint_type":                  CHECKSET,
		"db_proxy_read_write_mode":                CHECKSET,
		"db_proxy_instance_type":                  CHECKSET,
		"db_proxy_service_status":                 CHECKSET,
		"connection_persist":                      CHECKSET,
		"causal_consist_read":                     CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoRdsDescribedbproxyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBProxyConfigDependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers: testAccProviders,
		// CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":        "${alibabacloudstack_db_instance.default.id}",
					"db_proxy_instance_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "1",
						"db_proxy_instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_connect_string_port": "8000",
					"db_proxy_connect_string":      "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_connect_string_port": "8000",
						"db_proxy_connect_string":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num": "2",
					//"effective_time":        "SpecificTime",
					//"effective_specific_time": "${local.thrity_seconds_after}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num": "2",
					}),
				),
			},
			//			{
			//				Config: testAccConfig(map[string]interface{}{
			//					"persistent_connection_status": "Enabled",
			//				}),
			//				Check: resource.ComposeTestCheckFunc(
			//					testAccCheck(map[string]string{
			//						"persistent_connection_status": "Enabled",
			//					}),
			//				),
			//			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_persist": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_persist": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"causal_consist_read": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"causal_consist_read": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceDBProxyConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
	
	%s
	%s

locals {
  thrity_seconds_after = formatdate("YYYY-MM-DD'T'hh:mm:ssZ", timeadd(timestamp(), "30s"))
}

	`, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
