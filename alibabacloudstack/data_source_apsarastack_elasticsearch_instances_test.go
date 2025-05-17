package alibabacloudstack

import (
	"fmt"
	"testing"

	
)

func TestAccAlibabacloudStackElasticsearchDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_elasticsearch_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccES%d", rand),
		dataSourceElasticsearchConfigDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}-F",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_elasticsearch_instance.default.id}-F"},
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}",
			"ids":               []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
			"version":           "5.5.3_with_X-Pack",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_elasticsearch_instance.default.description}-F",
			"ids":               []string{"${alibabacloudstack_elasticsearch_instance.default.id}"},
			"version":           "6.7.0_with_X-Pack",
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}

	var elasticsearchCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existElasticsearchMapFunc,
		fakeMapFunc:  fakeElasticsearchMapFunc,
	}
	elasticsearchCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, idsConf, tagsConf, allConf)
}

var existElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                            "1",
		"ids.0":                            CHECKSET,
		"descriptions.#":                   "1",
		"descriptions.0":                   fmt.Sprintf("tf-testAccES%d", rand),
		"instances.#":                      "1",
		"instances.0.id":                   CHECKSET,
		"instances.0.instance_charge_type": string(PostPaid),
		"instances.0.description":          fmt.Sprintf("tf-testAccES%d", rand),
		"instances.0.data_node_amount":     "2",
		"instances.0.data_node_spec":       "elasticsearch.sn2ne.large",
		"instances.0.status":               "active",
		"instances.0.version":              CHECKSET,
		"instances.0.tags.%":               CHECKSET,
		"instances.0.created_at":           CHECKSET,
		"instances.0.updated_at":           CHECKSET,
		"instances.0.vswitch_id":           CHECKSET,
	}
}

var fakeElasticsearchMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"instances.#":    "0",
		"ids.#":          "0",
		"descriptions.#": "0",
	}
}

func dataSourceElasticsearchConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alibabacloudstack_elasticsearch_zones" "default" {}
data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alibabacloudstack_vswitches" "default" {
  vpc_id = data.alibabacloudstack_vpcs.default.ids.0
  zone_id = data.alibabacloudstack_elasticsearch_zones.default.ids.0
}
resource "alibabacloudstack_vswitch" "vswitch" {
  count             = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alibabacloudstack_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alibabacloudstack_vpcs.default.vpcs[0].cidr_block, 8, 8)
  availability_zone = data.alibabacloudstack_elasticsearch_zones.default.ids.0
 name              = var.name
}

locals {
  vswitch_id = length(data.alibabacloudstack_vswitches.default.ids) > 0 ? data.alibabacloudstack_vswitches.default.ids[0] : concat(alibabacloudstack_vswitch.vswitch.*.id, [""])[0]
}

resource "alibabacloudstack_elasticsearch_instance" "default" {
  description          = var.name
  password             = "%s"
  vswitch_id           = local.vswitch_id
  data_node_amount     = "2"
  data_node_spec       = "elasticsearch.sn2ne.large"
  data_node_disk_size  = "20"
  data_node_disk_type  = "cloud_ssd"
  instance_charge_type = "PostPaid"
  version              = "5.5.3_with_X-Pack"
  tags                 = {
	  "Created": "TF",
	  "For":     "acceptance test",
  }
}
`, name, GeneratePassword())
}
