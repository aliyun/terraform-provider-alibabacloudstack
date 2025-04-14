---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_ciphertext"
sidebar_current: "docs-alibabacloudstack-datasource-kms-ciphertext"
description: |-
    Encrypt data with KMS.
---

# alibabacloudstack_kms_ciphertext

Encrypt a given plaintext with KMS. 

~> **NOTE**: Using this data provider will allow you to conceal secret data within your resource definitions but does not take care of protecting that data in all Terraform logging and state output. Please take care to secure your secret data beyond just the Terraform configuration.

## Example Usage

```
resource "alibabacloudstack_kms_key" "key" {
  description             = "example key"
  is_enabled              = true
}

data "alibabacloudstack_kms_ciphertext" "encrypted" {
  key_id    = alibabacloudstack_kms_key.key.id
  plaintext = "example"
}

output "alibabacloudstack_kms_ciphertext" {
  value = "${data.alibabacloudstack_kms_ciphertext.encrypted}"
}
```

## Argument Reference

The following arguments are supported:

* `plaintext` - The plaintext to be encrypted which must be encoded in Base64. 
* `key_id` - The globally unique ID of the CMK.
* `encryption_context` - 
  (Optional) The Encryption context. If you specify this parameter here, it is also required when you call the Decrypt API operation. 
* `sensitive` - (Required) Indicates whether the attribute is sensitive.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ciphertext_blob` - The ciphertext of the data key encrypted with the primary CMK version.  
* `sensitive` - Indicates whether the attribute is sensitive. This attribute is computed automatically.