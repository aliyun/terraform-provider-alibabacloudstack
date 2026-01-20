package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackEcsImagesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_images.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceImagesConfigDependence)
	ownerNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*-fake",
			"owners":     "self",
		}),
	}

	ownerRecentConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"most_recent": "true",
			"owners":      "system",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "^win.*",
			"most_recent": "true",
			"owners":      "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":  "^win.*-fake",
			"most_recent": "true",
			"owners":      "system",
		}),
	}

	var existImagesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"images.#":                        CHECKSET,
			"images.0.architecture":           CHECKSET,
			"images.0.disk_device_mappings.#": CHECKSET,
			"images.0.creation_time":          CHECKSET,
			"images.0.image_id":               CHECKSET,
			"images.0.image_owner_alias":      CHECKSET,
			"images.0.os_type":                CHECKSET,
			"images.0.name":                   CHECKSET,
			"images.0.os_name":                CHECKSET,
			"images.0.progress":               "100%",
			"images.0.state":                  "Available",
			"images.0.status":                 "Available",
			"images.0.usage":                  CHECKSET,
			"images.0.tags.%":                 "0",
		}
	}

	var fakeImagesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"images.#": "0",
		}
	}

	var imagesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existImagesMapFunc,
		fakeMapFunc:  fakeImagesMapFunc,
	}

	imagesCheckInfo.dataSourceTestCheck(t, rand, ownerNameRegexConf, ownerRecentConf, allConf)
}

func dataSourceImagesConfigDependence(name string) string {
	return ""
}
