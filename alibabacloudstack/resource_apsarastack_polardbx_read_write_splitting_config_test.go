package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

func TestAccAlibabacloudStackPolardbxReadWriteSplittingConfig_basic0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_polardbx_read_write_splitting_config.default"
	ra := resourceAttrInit(resourceId, map[string]string{})

	serviceFunc := func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolardbxReadWriteSplittingConfig")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tfacc-polardbx-rws-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxReadWriteSplittingConfigDependence)

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
					"db_instance_id": "${local.polardbx_instance.id}",
					"attend_htap_list": []string{
						"${local.polardbx_instance.id}",
						"${alibabacloudstack_polardbx_readonly_instance.default.id}",
					},
					"auto_attend_htap":               "true",
					"delay_execution_strategy":       "1",
					"enable_consistent_replica_read": "true",
					"enable_htap":                    "true",
					"master_read_weight":             "30",
					"storage_delay_threshold":        "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attend_htap_list.#":             "2",
						"auto_attend_htap":               "true",
						"delay_execution_strategy":       "1",
						"enable_consistent_replica_read": "true",
						"enable_htap":                    "true",
						"master_read_weight":             "30",
						"storage_delay_threshold":        "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_attend_htap":               "false",
					"delay_execution_strategy":       "0",
					"enable_consistent_replica_read": "false",
					"master_read_weight":             "40",
					"storage_delay_threshold":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_attend_htap":               "false",
						"delay_execution_strategy":       "0",
						"enable_consistent_replica_read": "false",
						"master_read_weight":             "40",
						"storage_delay_threshold":        "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_htap": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_htap": "true",
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

func resourcePolardbxReadWriteSplittingConfigDependence(name string) string {

	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s

resource "alibabacloudstack_polardbx_readonly_instance" "default" {
	primary_db_instance_id = "${local.polardbx_instance.id}"
	zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage        = 50
	vswitch_id     = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id}"
	cn_node_count  = "2"
	dn_node_class  = "${data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id}"
	dn_node_count  = "2"
}

 `, name, VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase())
}
