package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbSharedInstance_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_polardb_shared_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBSharedInstance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardb-shared-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceoPolardbSharedInstanceDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": "${var.name}",
					"zone_id":                "cn-ece-entity46-amtest11001-a",
					"db_version":             "5.7",
					"storage_space":          "50",
					"vpc_id":                 "vpc-wz9i8dcp5yiq2me6f6ndc",
					"vswitch_id":             "vsw-wz9bb8mem4qmwqnmqjcs2",
					"db_node_class":          "polar.mysql.x1.small.c",
					// "sub_category":           "normal_exclusive",
					"db_type":      "MySQL",
					"db_node_num":  "2",
					"storage_type": "ESSDPL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
						"db_version":             "5.7",
						"storage_space":          "50",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"db_node_class":          "polar.mysql.x1.small.c",
						// "sub_category":           "normal_exclusive",
						"db_type":      "MySQL",
						"db_node_num":  "2",
						"storage_type": "ESSDPL1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": "${var.name}_update",
					"deletion_lock":          1,
					"public_connection_port": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": fmt.Sprintf("%s_update", name),
						"deletion_lock":          "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips_groups": map[string]interface{}{
						"test1": "10.0.0.10,192.168.1.10/24",
						"test2": "10.0.0.1,192.168.1.1/24",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips_groups.#":     "2",
						"security_ips_groups.test1": "10.0.0.10,192.168.1.10/24",
						"security_ips_groups.test2": "10.0.0.1,192.168.1.1/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "true",
					"enable_tde":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "true",
						"tde_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_num": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_num": "3",
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

func resourceoPolardbSharedInstanceDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}	
 `, name)
}
