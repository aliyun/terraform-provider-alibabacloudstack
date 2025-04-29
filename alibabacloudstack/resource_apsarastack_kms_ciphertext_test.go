package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"testing"
)

func TestAccAlibabacloudStackKmsCiphertext_basic(t *testing.T) {
	resourceId := "alibabacloudstack_kms_ciphertext.default"
	ResourceTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_basic(getAccTestRandInt(1000000, 9999999)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "ciphertext_blob"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlibabacloudStackKmsCiphertext_validate(t *testing.T) {

	ResourceTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_validate(getAccTestRandInt(1000000, 9999999)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alibabacloudstack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackKmsCiphertext_validate_withContext(t *testing.T) {
	ResourceTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_validate_withContext(getAccTestRandInt(1000000, 9999999)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alibabacloudstack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

var testAccAlibabacloudStackKmsCiphertextConfig_basic = func(keyId int) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_kms_key" "default" {
  	description = "tf-testacc-basic%d"
	is_enabled  = true
}

resource "alibabacloudstack_kms_ciphertext" "default" {
	key_id = "${alibabacloudstack_kms_key.default.id}"
	plaintext = "plaintext"
}
`, keyId)
}

var testAccAlibabacloudStackKmsCiphertextConfig_validate = func(keyId int) string {
	return fmt.Sprintf(`
	resource "alibabacloudstack_kms_key" "default" {
        description = "tf-testacc-validate%d"
	}
	
	resource "alibabacloudstack_kms_ciphertext" "default" {
		key_id = "${alibabacloudstack_kms_key.default.id}"
		plaintext = "plaintext"
	}
	`, keyId)
}

var testAccAlibabacloudStackKmsCiphertextConfig_validate_withContext = func(keyId int) string {
	return fmt.Sprintf(`
	resource "alibabacloudstack_kms_key" "default" {
        description = "tf-testacc-validate-withcontext%d"
	}
	
	resource "alibabacloudstack_kms_ciphertext" "default" {
		key_id = "${alibabacloudstack_kms_key.default.id}"
		plaintext = "plaintext"
        encryption_context = {
    		name = "value"
  		}
	}
	
	
	`, keyId)
}
