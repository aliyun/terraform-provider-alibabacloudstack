package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxInstance_basic0(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, map[string]string{"private_connection_string": CHECKSET})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardbx-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxDependence)

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
					"description":    "testtf",
					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version": "5.7",
					"storage":        "50",
					"vswitch_id":     "${alibabacloudstack_vpc_vswitch.default.id}",
					"gms_node_class": "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.1.id}",
					"cn_node_class":  "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}",
					"cn_node_count":  "2",
					"dn_node_class":  "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}",
					"dn_node_count":  "2",
					"compute_parameters": []map[string]interface{}{
						{
							"name":  "CONN_POOL_BLOCK_TIMEOUT",
							"value": "40000",
						},
						{
							"name":  "CONN_POOL_IDLE_TIMEOUT",
							"value": "35",
						},
					},
					"storage_parameters": []map[string]interface{}{
						{
							"name":  "innodb_stats_sample_pages",
							"value": "10",
						},
						{
							"name":  "ft_query_expansion_limit",
							"value": "35",
						},
					},
					"security_groups": []map[string]interface{}{
						{
							"group_name": "test123",
							"ips":        "10.0.0.1,10.0.0.2",
						},
					},
					"private_connection_string_prefix": name + "-pr",
					"private_connection_port":          4001,
					"enable_public_connection":         true,
					"public_connection_string_prefix":  name + "-pu",
					"public_connection_port":           4002,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						// "description":                  name,
						"cn_node_count":                "2",
						"dn_node_count":                "2",
						"compute_parameters.#":         "2",
						"compute_parameters.0.name":    CHECKSET,
						"storage_parameters.#":         "2",
						"storage_parameters.0.value":   CHECKSET,
						"security_groups.#":            "1",
						"security_groups.0.group_name": "test123",
						"security_groups.0.ips":        "10.0.0.1,10.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_public_connection": false,
					"public_connection_string_prefix": REMOVEKEY,
					"public_connection_port":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"public_connection_string_prefix": REMOVEKEY,
						"public_connection_port":          REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []map[string]interface{}{
						{
							"group_name": "test123",
							"ips":        "10.0.0.1,192.168.1.1/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#":            "1",
						"security_groups.0.group_name": "test123",
						"security_groups.0.ips":        "10.0.0.1,192.168.1.1/24",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "true",
					"enable_tde": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "true",
						"enable_tde": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_ssl": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_ssl": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_count": "3",
					"dn_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "3",
						"dn_node_count": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dn_node_count": "2",
					"cn_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dn_node_count": "2",
						"cn_node_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_class": "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.1.id}",
					"dn_node_class": "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// "topology_type" does not support readback, "compute_parameters", "storage_parameters" only support reading back locally updated parameters.
				ImportStateVerifyIgnore: []string{"topology_type", "compute_parameters", "storage_parameters"},
			},
		},
	})
}

func resourcePolardbxDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

data "alibabacloudstack_polardbx_instance_types" "cn" {
	sorted_by = "CPU"
	spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
	sorted_by = "CPU"
	spec_type = "DN"
}
	
%s
 `, name, VSwitchCommonTestCase)
}
