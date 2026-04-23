package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackOssBucketsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_oss_buckets.default"
	name := fmt.Sprintf("tf-testacc-bucket-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOssBucketsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oss_bucket.demo.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_oss_bucket.demo.bucket}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_oss_bucket.demo.bucket}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-bucket-id-12345"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_oss_bucket.demo.bucket}"},
			"name_regex": "${alibabacloudstack_oss_bucket.demo.bucket}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"fake-bucket-id-12345"},
			"name_regex": "${alibabacloudstack_oss_bucket.demo.bucket}_fake",
		}),
	}

	var existOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "1",
			"names.#":        "1",
			"buckets.#":      "1",
			"buckets.0.name": name,
			// Other bucket attributes are computed but we don't know exact values
			// so we only validate the fields that are guaranteed to be present and known
		}
	}

	var fakeOssBucketsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"buckets.#": "0",
		}
	}

	var ossBucketsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOssBucketsMapFunc,
		fakeMapFunc:  fakeOssBucketsMapFunc,
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckOss(t)
		},
	}
	ossBucketsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceOssBucketsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_oss_bucket" "demo" {
  bucket = var.name
  acl    = "public-read"
}

`, name)
}
