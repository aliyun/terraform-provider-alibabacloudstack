package alibabacloudstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackKmsCiphertext_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_basic(acctest.RandomWithPrefix("tf-testacc-basic")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"alibabacloudstack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackKmsCiphertext_validate(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_validate(acctest.RandomWithPrefix("tf-testacc-validate")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alibabacloudstack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccAlibabacloudStackKmsCiphertext_validate_withContext(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlibabacloudStackKmsCiphertextConfig_validate_withContext(acctest.RandomWithPrefix("tf-testacc-validate-withcontext")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alibabacloudstack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

var testAccAlibabacloudStackKmsCiphertextConfig_basic = func(keyId string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_kms_key" "default" {
  	description = "%s"
	is_enabled  = true
}

resource "alibabacloudstack_kms_ciphertext" "default" {
	key_id = "${alibabacloudstack_kms_key.default.id}"
	plaintext = "plaintext"
}
`, keyId)
}

var testAccAlibabacloudStackKmsCiphertextConfig_validate = func(keyId string) string {
	return fmt.Sprintf(`
	resource "alibabacloudstack_kms_key" "default" {
        description = "%s"
	}
	
	resource "alibabacloudstack_kms_ciphertext" "default" {
		key_id = "${alibabacloudstack_kms_key.default.id}"
		plaintext = "plaintext"
	}
	`, keyId)
}

var testAccAlibabacloudStackKmsCiphertextConfig_validate_withContext = func(keyId string) string {
	return fmt.Sprintf(`
	resource "alibabacloudstack_kms_key" "default" {
        description = "%s"
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
