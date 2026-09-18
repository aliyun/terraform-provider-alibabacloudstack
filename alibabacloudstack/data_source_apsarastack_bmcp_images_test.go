package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackBmcpImagesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_bmcp_images.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccBmcpImages%d", rand),
		dataSourceBmcpImagesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"images.#":        CHECKSET,
				"images.0.id":     CHECKSET,
				"images.0.name":   CHECKSET,
				"images.0.type":   CHECKSET,
				"images.0.status": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"images.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":       "public",
			"name_regex": "ubuntu",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":       "public",
			"name_regex": "fake-name",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type": "public",
			"ids":  []string{"i-d"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type": "public",
			"ids":  []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"type":       "public",
			"name_regex": "ubuntu",
			"ids":        []string{"i-d"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"type":       "public",
			"name_regex": "fake-name",
			"ids":        []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceBmcpImagesConfigDependence(name string) string {
	return fmt.Sprintf(``)
}
