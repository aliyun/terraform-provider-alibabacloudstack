package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackRedisParameterGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_kvstore_parameter_group.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccRedisParameterGroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KvstoreService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeParameterGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-kvparmgroup%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccRedisParameterGroupBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"character_type":       "logic",
					"engine":               "Redis",
					"parameter_group_name": "${var.name}",
					"engine_version":       "7.0",
					"parameter_group_desc": "${var.name}",
					"parameters": []map[string]interface{}{
						{
							"param_name": "resp_version",
							"value":      "3",
						},
						{
							"param_name": "rt_threshold_ms",
							"value":      "400",
						},
						{
							"param_name": "#no_loose_check-whitelist-always",
							"value":      "yes",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"character_type":       "logic",
						"engine":               "Redis",
						"parameter_group_name": name,
						"engine_version":       "7.0",
						"parameter_group_desc": name,
						"parameters.#":         "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"engine"},
			},
		},
	})
}

var AlibabacloudTestAccRedisParameterGroupCheckmap = map[string]string{
	"create_time": CHECKSET,
	"type":        CHECKSET,
	"is_dynamic":  CHECKSET,
}

func AlibabacloudTestAccRedisParameterGroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
