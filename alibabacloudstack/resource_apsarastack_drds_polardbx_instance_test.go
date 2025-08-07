package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackDrdsPolardbxInstance_basic0(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_drds_polardbx_instance.default"
	ra := resourceAttrInit(resourceId, drdsPolardbxbasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_polardb_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsPolardbxDependence)

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
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"cn_node_count": "2",
						"dn_node_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cn_node_count": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cn_node_count": "3",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"primary_db_instance_id", "topology_type"},
			},
		},
	})
}

func TestAccAlibabacloudStackDrdsPolardbxInstance_onlyreadInstance(t *testing.T) {
	var v *PolardbxDescribedbinstanceattributeResponse

	resourceId := "alibabacloudstack_drds_polardbx_instance.onlyread_instance"
	ra := resourceAttrInit(resourceId, drdsPolardbxbasicMap)

	serviceFunc := func() interface{} {
		return &DrdsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DoPolardbxDescribedbinstanceattributeRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_drds_polardb_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDrdsPolardbxOnlyReadInstanceDependence)

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
					"primary_db_instance_id": "${alibabacloudstack_drds_polardbx_instance.default.id}",
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

func resourceDrdsPolardbxDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
		default = "%s"
	}

	%s
`, name, VSwitchCommonTestCase)
}

var drdsPolardbxbasicMap = map[string]string{}

func resourceDrdsPolardbxOnlyReadInstanceDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}
%s
resource "alibabacloudstack_drds_polardbx_instance" "default" {
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
