package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackDBProxy_basic(t *testing.T) {
	var v *RdsDescribedbproxyResponse
	rand := getAccTestRandInt(10000, 999999)
	name := fmt.Sprintf("tf-testAccdbproxy-%d", rand)
	resourceId := "alibabacloudstack_db_proxy.default"
	ra := resourceAttrInit(resourceId, map[string]string{})
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
					"db_instance_id":          "${alibabacloudstack_db_instance.instance.id}",
					"config_db_proxy_service": "Startup",
					// "db_instance_id":          "rm-rw351e34bqk41yde5",
					"db_proxy_instance_num": "1",
					"instance_network_type": "VPC",
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
					"db_proxy_instance_num":   "3",
					"db_proxy_instance_type":  "common",
					"effective_specific_time": "${var.effective_specific_time}",
					"effective_time":          "SpecificTime",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "3",
						"db_proxy_instance_type": "common",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceDBProxyConfigDependence(name string) string {
	now := time.Now().UTC()
	effective_specific_time := now.Add(time.Minute * 30)
	return fmt.Sprintf(`
	variable "effective_specific_time" {
		default = "%s"
	}
	variable "name" {
		default = "%v"
	}

	resource "alibabacloudstack_db_instance" "instance" {
		engine               = "MySQL"
	    engine_version       = "5.7"
	    instance_type        = "rds.mysql.s2.large"
	    instance_storage     = "5"
	    instance_name 		 = "${var.name}"
	    storage_type         = "local_ssd"
	}
	
	`, effective_specific_time.Format("2006-01-02T15:04Z"), name)
}
