package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackKmsCiphertextDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_kms_ciphertext.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceKmsCiphertextDependence)

	plaintextConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alibabacloudstack_kms_key.default.id}",
			"plaintext": "plaintext",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alibabacloudstack_kms_key.default.id}",
			"plaintext": "plaintext",
			"encryption_context": map[string]string{
				"key": "value",
			},
		}),
	}

	var existKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": CHECKSET,
		}
	}

	var fakeKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": NOSET,
		}
	}

	var kmsCipherCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsCiphertextMapFunc,
		fakeMapFunc:  fakeKmsCiphertextMapFunc,
	}

	kmsCipherCheckInfo.dataSourceTestCheck(t, 0, plaintextConf, allConf)
}

func dataSourceKmsCiphertextDependence(name string) string {
	return `
	resource "alibabacloudstack_kms_key" "default" {
    	is_enabled = true
	}
	`
}
