package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackApiGatewayV2CascadeLinksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alibabacloudstack_api_gateway_v2_cascade_links.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccApiGwV2%d", rand),
		dataSourceApiGatewayV2CascadeLinksConfigDependence)

	// Define the dataSourceAttr with existMapFunc and fakeMapFunc
	dsa := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                           "1",
				"links.#":                         "1",
				"links.0.link_name":               fmt.Sprintf("tf-testAccApiGwV2%d", rand),
				"links.0.source_instance_id":      CHECKSET,
				"links.0.cascade_instance_id":     CHECKSET,
				"links.0.source_instance_address": "10.17.94.180",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":   "0",
				"links.#": "0",
			}
		},
	}

	// Run the test check with different configurations
	dsa.dataSourceTestCheck(t, rand,
		// Test by name_regex
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_api_gateway_v2_cascade_link.default.link_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "fake-name-regex",
			}),
		},
		// Test by ids
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_api_gateway_v2_cascade_link.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		// Test by source_instance_name
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"source_instance_name": "${alibabacloudstack_api_gateway_v2_cascade_link.default.source_instance_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"source_instance_name": "fake-source-instance-name",
			}),
		},
		// Test by cascade_instance_name
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"cascade_instance_name": "${alibabacloudstack_api_gateway_v2_cascade_link.default.cascade_instance_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"cascade_instance_name": "fake-cascade-instance-name",
			}),
		},
	)
}

func dataSourceApiGatewayV2CascadeLinksConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_api_gateway_v2_instance" "source" {
  instance_name      = "${var.name}-source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_instance" "cascade" {
  instance_name      = "${var.name}-cascade"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = "${var.name}"
  cascade_instance_id = "${alibabacloudstack_api_gateway_v2_instance.cascade.id}"
}

resource "alibabacloudstack_api_gateway_v2_cascade_link" "default" {
  source_instance_id      = "${alibabacloudstack_api_gateway_v2_instance.source.id}"
  source_instance_address = "10.17.94.180"
  cascade_instance_id     = "${alibabacloudstack_api_gateway_v2_cascade_instance.default.id}"
  link_name               = "${var.name}"
}
`, name)
}
