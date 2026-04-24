package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOssSingleTunnelsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_oss_single_tunnels.default"
	name := fmt.Sprintf("tf-tunnel-data%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		name,
		dataSourceOssSingleTunnelsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oss_single_tunnel.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oss_single_tunnel.default.id}-fake"},
		}),
	}

	clusterConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oss_single_tunnel.default.id}"},
			"oss_cluster": "${alibabacloudstack_oss_single_tunnel.default.cluster}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_oss_single_tunnel.default.id}"},
			"oss_cluster": "${alibabacloudstack_oss_single_tunnel.default.cluster}-fake",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oss_single_tunnel.default.label}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oss_single_tunnel.default.label}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_oss_single_tunnel.default.id}"},
			"name_regex": "${alibabacloudstack_oss_single_tunnel.default.label}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_oss_single_tunnel.default.id}-fake"},
			"name_regex": "${alibabacloudstack_oss_single_tunnel.default.label}-fake",
		}),
	}

	var existOtstunnelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"tunnels.#":         "1",
			"tunnels.0.id":      CHECKSET,
			"tunnels.0.label":   name,
			"tunnels.0.cluster": CHECKSET,
			"tunnels.0.vip":     CHECKSET,
			"tunnels.0.vpc_id":  CHECKSET,
			"tunnels.0.shared":  CHECKSET,
		}
	}

	var fakeOtstunnelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"tunnels.#": "0",
		}
	}

	var otstunnelsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtstunnelsMapFunc,
		fakeMapFunc:  fakeOtstunnelsMapFunc,
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOssEndpointList(t)
		},
	}

	otstunnelsCheckInfo.dataSourceTestCheck(t, rand, idsConf, clusterConf, nameRegexConf, allConf)
}

func dataSourceOssSingleTunnelsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_oss_clusters" "default" {
}

%s

resource "alibabacloudstack_oss_single_tunnel" "default" {
  shared = "0"
  label = "${var.name}"
  cluster = "${data.alibabacloudstack_oss_clusters.default.clusters.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
	`, name, VSwitchCommonTestCase)
}
