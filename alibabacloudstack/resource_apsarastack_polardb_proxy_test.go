package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbProxy_basic(t *testing.T) {
	var v *PolardbDescribedbproxyResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbproxy-%d", rand)
	resourceId := "alibabacloudstack_polardb_proxy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbDescribedbproxyRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbProxyConfigDependence)
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
					"db_instance_id":        "${alibabacloudstack_polardb_dbinstance.default.id}",
					"db_proxy_instance_num": "1",
					"instance_network_type": "Classic",
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
					"db_proxy_instance_num": "2",
					"effective_time":        "Immediate",
					// "effective_specific_time": "${var.db_proxy_effective_time}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"persistent_connection_status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"persistent_connection_status": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_connect_string_port": "8000",
					"db_proxy_connect_string":      "test1234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_connect_string": "test1234",
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

func resourcePolardbProxyConfigDependence(name string) string {
	now := time.Now().UTC()
	oneminago := now.Add(+1 * time.Minute)
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}
		
	variable "db_proxy_effective_time" {
		default = "%v"
	}

	resource "alibabacloudstack_polardb_dbinstance" "default" {
	instance_storage = "5"
	instance_name = "${var.name}"
	storage_type = "local_ssd"
	engine = "MySQL"
	engine_version = "5.7"
	instance_type = "rds.mysql.t1.small"
	}
	`, name, oneminago.Format("2006-01-02T15:04:05Z"))
}
