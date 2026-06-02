---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_ciphertexts"
description: |-
    Encrypt data with KMS.
---

# alibabacloudstack_kms_ciphertexts

Encrypt a given plaintext with KMS. 

~> **NOTE**: Using this data provider will allow you to conceal secret data within your resource definitions but does not take care of protecting that data in all Terraform logging and state output. Please take care to secure your secret data beyond just the Terraform configuration.

## Example Usage

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

data "alibabacloudstack_kms_ciphertexts" "encrypted" {
  key_id    = alibabacloudstack_kms_key.key.id
  plaintext = "example"
}

output "alibabacloudstack_kms_ciphertexts" {
  value = "${data.alibabacloudstack_kms_ciphertexts.encrypted}"
}
```

## Argument Reference

The following arguments are supported:

* `plaintext` - (Required) The plaintext data to be encrypted.
* `key_id` - (Required) The globally unique ID of the CMK.
* `encryption_context` - (Optional) The Encryption context. If you specify this parameter here, it is also required when you call the Decrypt API operation.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ciphertext_blob` - The ciphertext of the data key encrypted with the primary CMK version. 