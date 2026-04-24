package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackKmsCiphertext_basic(t *testing.T) {
	resourceId := "alibabacloudstack_kms_ciphertext.default"
	ra := resourceAttrInit(resourceId, KmsCiphertextMap)
	// rc := resourceCheckInit(resourceId, nil, func() interface{} {
	// 	return &KmsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	// })
	// rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-kms-ciphertext%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KmsCiphertextBasicDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		// This resource only has local attributes and does not support deletion.
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_id":             "${alibabacloudstack_kms_key.default.id}",
					"plaintext":          "plaintext-basic",
					"encryption_context": map[string]string{"name": "value"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ciphertext_blob":         CHECKSET,
						"encryption_context.%":    "1",
						"encryption_context.name": "value",
					}),
				),
			},
		},
	})
}

var KmsCiphertextMap = map[string]string{
	"ciphertext_blob": CHECKSET,
}

func KmsCiphertextBasicDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_kms_key" "default" {
  	description = "%s"
	is_enabled  = true
}
`, name)
}
