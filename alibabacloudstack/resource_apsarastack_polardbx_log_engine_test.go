package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackPolardbxLogEngine_basic(t *testing.T) {
	var v []PolarDbXLogEngineInfo
	resourceId := "alibabacloudstack_polardbx_log_engine.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"cdc_node_class": CHECKSET,
		"cdc_node_count": CHECKSET,
	})
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &PolardbXService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1, 254)
	name := fmt.Sprintf("tf-testAccPolardbxLogengine_%v", rand)
	groupName := fmt.Sprintf("group%vt", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolardbxLogEngineBasicDependence)

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
					"instance_id":    "${local.polardbx_instance.id}",
					"cdc_node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
					"cdc_node_count": 2,
					"multi_stream": []map[string]interface{}{
						{
							"group_name": groupName + "1",
							"comment":    "test1",
							"hash_level": "RECORD",
							"node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
							"node_count": 3,
						},
						{
							"group_name": groupName + "2",
							"comment":    "test2",
							"hash_level": "RECORD",
							"node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.1.id}",
							"node_count": 2,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cdc_node_count": "2",
						"multi_stream.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"multi_stream.*",
						map[string]string{
							"group_name": groupName + "1",
							"hash_level": "RECORD",
							"node_count": "3",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"multi_stream.*",
						map[string]string{
							"group_name": groupName + "2",
							"hash_level": "RECORD",
							"node_count": "2",
						},
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cdc_node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.1.id}",
					"cdc_node_count": 3,
					"multi_stream": []map[string]interface{}{
						{
							"group_name": groupName + "1",
							"comment":    "test1",
							"hash_level": "DATABASE",
							"node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
							"node_count": 2,
						},
						{
							"group_name": groupName + "2",
							"comment":    "test2",
							"hash_level": "RECORD",
							"node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
							"node_count": 2,
						},
						{
							"group_name": groupName + "3",
							"comment":    "test3",
							"hash_level": "RECORD",
							"node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
							"node_count": 3,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cdc_node_count": "3",
						"multi_stream.#": "3",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"multi_stream.*",
						map[string]string{
							"group_name": groupName + "1",
							"hash_level": "DATABASE",
							"node_count": "2",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"multi_stream.*",
						map[string]string{
							"group_name": groupName + "2",
							"hash_level": "RECORD",
							"node_count": "2",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceId,
						"multi_stream.*",
						map[string]string{
							"group_name": groupName + "3",
							"hash_level": "RECORD",
							"node_count": "3",
						},
					),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cdc_node_class": "${data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id}",
					"cdc_node_count": 2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cdc_node_count": "2",
						"multi_stream.#": REMOVEKEY,
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

func resourcePolardbxLogEngineBasicDependence(name string) string {
	return fmt.Sprintf(`
%s
data "alibabacloudstack_polardbx_cdc_classes" "default" {
	instance_id= local.polardbx_instance.id
	sorted_by= "CPU"
}`, resourcePolardbxBackupBasicDependence(name))
}
