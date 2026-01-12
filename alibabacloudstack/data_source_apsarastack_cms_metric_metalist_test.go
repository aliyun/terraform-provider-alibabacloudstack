package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackCmsMetricMetalistDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cms_metric_metalist.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceCmsMetricMetalistConfigDependence)

	// Test with valid namespace that returns metrics
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace": "acs_slb_dashboard",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "fake-nonexistent-namespace-12345",
		}),
	}

	var existCmsMetricMetalistMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                CHECKSET, // Should contain at least one metric
			"resources.0.metric_name":    CHECKSET, // First metric should have a name
			"resources.0.namespace":      "acs_slb_dashboard", // Should match input namespace
			"resources.0.periods":        CHECKSET, // Should have periods
			"resources.0.description":    CHECKSET, // Should have description
			// Other fields may be empty strings but should be set
		}
	}

	var fakeCmsMetricMetalistMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": "0", // Fake namespace should return empty list
		}
	}

	var cmsMetricMetalistCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCmsMetricMetalistMapFunc,
		fakeMapFunc:  fakeCmsMetricMetalistMapFunc,
	}
	cmsMetricMetalistCheckInfo.dataSourceTestCheck(t, 0, namespaceConf)
}

func dataSourceCmsMetricMetalistConfigDependence(name string) string {
	return `
`
}
