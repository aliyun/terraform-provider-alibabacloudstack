package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackHologramClustersDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_hologram_clusters.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceHologramClustersConfigDependence)

	// Basic configuration with zone_id only
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
		}),
	}

	// Test with ids filter
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
			"ids":     []string{`${data.alibabacloudstack_hologram_clusters.anyone.clusters.0.id}`},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
			"ids":     []string{"fake-nonexistent-cluster"},
		}),
	}

	// Test with name_regex filter
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":    `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
			"name_regex": `^${data.alibabacloudstack_hologram_clusters.anyone.clusters.0.cluster}$`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":    `${data.alibabacloudstack_zones.anyone.zones.0.id}`,
			"name_regex": "fake-nonexistent-cluster",
		}),
	}

	var existHologramClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           CHECKSET,
			"clusters.#":      CHECKSET,
			"clusters.0.id":   CHECKSET,
			"clusters.0.type": CHECKSET,
			"clusters.0.cpu":  CHECKSET,
		}
	}

	var fakeHologramClustersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"clusters.#": "0",
		}
	}

	var hologramClustersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existHologramClustersMapFunc,
		fakeMapFunc:  fakeHologramClustersMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
			testAccPreCheckApiAvailable(t, "POST", "Hologram", "2022-06-01", "ListClusters", "/api/v1/regions/listClusters")
		},
	}
	hologramClustersCheckInfo.dataSourceTestCheck(t, 0, basicConf, idsConf, nameRegexConf)
}

func dataSourceHologramClustersConfigDependence(name string) string {
	return `
data "alibabacloudstack_zones" "anyone" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "anyone" {
  zone_id    = data.alibabacloudstack_zones.anyone.zones.0.id
  name_regex = "Hologres.*"
}
`
}
