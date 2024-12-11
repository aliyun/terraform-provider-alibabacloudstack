package alibabacloudstack

import (
	"fmt"
	
	"testing"
)

func TestAccAlibabacloudStackImagesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_images.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceImagesConfigDependence)

	ownerConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"owners": "system",
		}),
	}
	ownerNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*-fake",
			"owners":     "system",
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

	imagesCheckInfo.dataSourceTestCheck(t, rand, ownerConf, ownerNameRegexConf, ownerRecentConf, allConf)
}

func TestAccAlibabacloudStackImagesDataSource_win(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_images.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceImagesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^win.*fake",
			"owners":     "system",
		}),
	}

	var existImagesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"images.#":                        CHECKSET,
			"images.0.architecture":           CHECKSET,
			"images.0.disk_device_mappings.#": "0",
			"images.0.creation_time":          CHECKSET,
			"images.0.image_id":               CHECKSET,
			"images.0.image_owner_alias":      CHECKSET,
			"images.0.os_type":                "windows",
			"images.0.name":                   CHECKSET,
			"images.0.os_name":                REGEXMATCH + "^Windows Server.*版.*",
			"images.0.os_name_en":             CHECKSET,
			"images.0.progress":               "100%",
			"images.0.state":                  "Available",
			"images.0.status":                 "Available",
			"images.0.usage":                  "instance",
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

	imagesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func TestAccAlibabacloudStackImagesDataSource_linux_english(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_images.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceImagesConfigDependence)

	slesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^sles.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^sles.*fake",
			"owners":     "system",
		}),
	}

	openSuseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^opensuse.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^opensuse.*fake",
			"owners":     "system",
		}),
	}

	coreOsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^coreos.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^coreos.*fake",
			"owners":     "system",
		}),
	}

	var existImagesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"images.#":                        CHECKSET,
			"images.0.architecture":           CHECKSET,
			"images.0.disk_device_mappings.#": "0",
			"images.0.creation_time":          CHECKSET,
			"images.0.image_id":               CHECKSET,
			"images.0.image_owner_alias":      CHECKSET,
			"images.0.os_type":                "linux",
			"images.0.name":                   CHECKSET,
			"images.0.os_name":                REGEXMATCH + "^.*Bit.*",
			"images.0.progress":               "100%",
			"images.0.state":                  "Available",
			"images.0.status":                 "Available",
			"images.0.usage":                  "instance",
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

	imagesCheckInfo.dataSourceTestCheck(t, rand, slesConf, openSuseConf, coreOsConf)
}

func TestAccAlibabacloudStackImagesDataSource_linux_chinese(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_images.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-%d", rand),
		dataSourceImagesConfigDependence)

	ubuntuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^ubuntu.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^ubuntu.*fake",
			"owners":     "system",
		}),
	}

	freebsdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^freebsd.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^freebsd.*fake",
			"owners":     "system",
		}),
	}

	centOsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^centos.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^centos.*fake",
			"owners":     "system",
		}),
	}

	debianConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^debian.*",
			"owners":     "system",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^debian.*fake",
			"owners":     "system",
		}),
	}

	var existImagesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"images.#":                        CHECKSET,
			"images.0.architecture":           CHECKSET,
			"images.0.disk_device_mappings.#": "0",
			"images.0.creation_time":          CHECKSET,
			"images.0.image_id":               CHECKSET,
			"images.0.image_owner_alias":      CHECKSET,
			"images.0.os_type":                "linux",
			"images.0.name":                   CHECKSET,
			"images.0.os_name":                REGEXMATCH + "^.*位.*",
			"images.0.progress":               "100%",
			"images.0.state":                  "Available",
			"images.0.status":                 "Available",
			"images.0.usage":                  "instance",
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

	imagesCheckInfo.dataSourceTestCheck(t, rand, ubuntuConf, freebsdConf, centOsConf, debianConf)
}

func dataSourceImagesConfigDependence(name string) string {
	return ""
}
