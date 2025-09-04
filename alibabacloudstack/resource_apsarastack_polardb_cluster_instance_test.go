package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbClusterInstance_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_polardb_cluster_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbClusterInstance")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardb-shared-%d", rand)
	// name := "tfacc-polardb-shared-18982"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceoPolardbClusterInstanceDependence)

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
					"zone_id":                "cn-wulan-env149-amtest149001-a",
					"db_version":             "14",
					"storage_space":          "20",
					"vpc_id":                 "vpc-c3d79nt47rh1mghkysi00",
					"vswitch_id":             "vsw-c3dijzuatzymqi5ah3pcz",
					"db_node_class":          "polar.pg.x4.medium.c",
					"db_type":                "PostgreSQL",
					"db_node_num":            "2",
					"storage_type":           "ESSDPL1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_description": name,
						"db_version":             "14",
						"storage_space":          "20",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"db_node_class":          "polar.pg.x4.medium.c",
						"db_type":                "PostgreSQL",
						"db_node_num":            "2",
						"storage_type":           "ESSDPL1",
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
						"security_ips_groups.%":     "2",
						"security_ips_groups.test1": CHECKSET,
						"security_ips_groups.test2": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"name":  "auto_explain.log_analyze",
							"value": "on",
						},
						{
							"name":  "auto_explain.sample_rate",
							"value": "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"parameters.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,     // 资源地址
						"parameters.*", // TypeSet 属性路径（通配符 `*` 表示集合中的任意元素）
						map[string]string{
							"name":  "auto_explain.log_analyze",
							"value": "on",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,     // 资源地址
						"parameters.*", // TypeSet 属性路径（通配符 `*` 表示集合中的任意元素）
						map[string]string{
							"name":  "auto_explain.sample_rate",
							"value": "0",
						},
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_description": "${var.name}_update",
					"deletion_lock":          1,
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
					"ssl_enabled": "true",
					"tde_enabled": "true",
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

func resourceoPolardbClusterInstanceDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}	
 `, name)
}
