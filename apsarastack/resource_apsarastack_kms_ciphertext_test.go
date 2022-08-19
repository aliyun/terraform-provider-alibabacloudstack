package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackKmsCiphertext_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccApsaraStackKmsCiphertextConfig_basic(acctest.RandomWithPrefix("tf-testacc-basic")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"apsarastack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccApsaraStackKmsCiphertext_validate(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccApsaraStackKmsCiphertextConfig_validate(acctest.RandomWithPrefix("tf-testacc-validate")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("apsarastack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccApsaraStackKmsCiphertext_validate_withContext(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccApsaraStackKmsCiphertextConfig_validate_withContext(acctest.RandomWithPrefix("tf-testacc-validate-withcontext")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("apsarastack_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

var testAccApsaraStackKmsCiphertextConfig_basic = func(keyId string) string {
	return fmt.Sprintf(`
resource "apsarastack_kms_key" "default" {
  	description = "%s"
	is_enabled  = true
}

resource "apsarastack_kms_ciphertext" "default" {
	key_id = "${apsarastack_kms_key.default.id}"
	plaintext = "plaintext"
}
`, keyId)
}

var testAccApsaraStackKmsCiphertextConfig_validate = func(keyId string) string {
	return fmt.Sprintf(`
	resource "apsarastack_kms_key" "default" {
        description = "%s"
	}
	
	resource "apsarastack_kms_ciphertext" "default" {
		key_id = "${apsarastack_kms_key.default.id}"
		plaintext = "plaintext"
	}
	`, keyId)
}

var testAccApsaraStackKmsCiphertextConfig_validate_withContext = func(keyId string) string {
	return fmt.Sprintf(`
	resource "apsarastack_kms_key" "default" {
        description = "%s"
	}
	
	resource "apsarastack_kms_ciphertext" "default" {
		key_id = "${apsarastack_kms_key.default.id}"
		plaintext = "plaintext"
        encryption_context = {
    		name = "value"
  		}
	}
	
	
	`, keyId)
}
