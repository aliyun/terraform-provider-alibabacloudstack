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
	ra := resourceAttrInit(resourceId, PolardbxbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_polardb_%d", rand)
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
					"description":    "testtf1111",
					"zone_id":        "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version": "5.7",
					"storage":        "50",
					"network_type":   "vpc",
					"vpc_id":         "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":     "${alibabacloudstack_vpc_vswitch.default.id}",
					"cn_node_class":  "polarx.x4.medium.2e",
					"cn_node_count":  "2",
					"dn_node_class":  "mysql.n4.medium.25",
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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                name,
						"cn_node_count":              "2",
						"dn_node_count":              "2",
						"compute_parameters.#":       "2",
						"compute_parameters.0.name":  CHECKSET,
						"storage_parameters.#":       "2",
						"storage_parameters.0.value": CHECKSET,
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
					"cn_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// "primary_db_instance_id", "topology_type" 不支持回读， "compute_parameters", "storage_parameters" 只支持回读本地更新的参数。
				ImportStateVerifyIgnore: []string{"primary_db_instance_id", "topology_type", "compute_parameters", "storage_parameters"},
			},
		},
	})
}

func TestAccAlibabacloudStackPolardbxInstance_onlyreadInstance(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_polardbx_instance.onlyread_instance"
	ra := resourceAttrInit(resourceId, PolardbxbasicMap)

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_polardb_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxOnlyReadInstanceDependence)

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
					"is_read_db_instance":    "true",
					"primary_db_instance_id": "${alibabacloudstack_polardbx_instance.default.id}",
					"description":            "${var.name}",
					"zone_id":                "${data.alibabacloudstack_zones.default.zones.0.id}",
					"engine_version":         "5.7",
					"storage":                "50",
					"network_type":           "vpc",
					"vpc_id":                 "${alibabacloudstack_vpc_vpc.default.id}",
					"vswitch_id":             "${alibabacloudstack_vpc_vswitch.default.id}",
					"cn_node_class":          "polarx.x4.medium.2e",
					"cn_node_count":          "2",
					"dn_node_class":          "mysql.n4.medium.25",
					"dn_node_count":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_read_db_instance": "true",
						"description":         name,
						"cn_node_count":       "2",
						"dn_node_count":       "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"primary_db_instance_id", "topology_type"},
			},
		},
	})
}

func resourcePolardbxDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}

`, name)
}

var PolardbxbasicMap = map[string]string{}

func resourcePolardbxOnlyReadInstanceDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}
%s
resource "alibabacloudstack_polardbx_instance" "default" {
 description = "testtf1111"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

 `, name, VSwitchCommonTestCase)
}
