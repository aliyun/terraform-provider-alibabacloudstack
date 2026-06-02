---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_keypair"
description: |-
  Provides a BCMP Keypair resource.
---

# alibabacloudstack_bcmp_keypair

Provides a BCMP Keypair resource.

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_bcmp_keypair" "basic" {
  key_pair_name = "terraform-test-key-pair"
}

// Import an existing public key to build a key pair
resource "alibabacloudstack_bcmp_keypair" "publickey" {
  key_pair_name = "my_public_key"
  public_key    = "-----BEGIN CERTIFICATE-----\nMIIDRjCCAq*******<Your Server Certificate String>*****bJJyOm5LqoiA=\n-----END CERTIFICATE-----"
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Optional, ForceNew) Name of the key pair. The value contains 2 to 128 English or Chinese characters. The name must start with a letter or a Chinese character. It cannot start with `http://` or `https://`. The value can contain digits, colons (`:`), underscores (`_`), or dashes (`-`). If not specified, Terraform will generate a unique name. The name is unique within the region.
* `public_key` - (Optional) An existing public key that you want to import and manage using AlibabaCloudStack key pairs. Leading and trailing whitespaces are trimmed. Modifying this parameter does not trigger resource recreation, but the change will not be persisted as the resource does not support updating the public key.
* `key_file` - (Optional, ForceNew) The file name where the private key of the newly created key pair will be saved. It is strongly recommended to specify this parameter when creating a key pair, as you will not be able to retrieve the private key afterward if it is not saved. Modifying this parameter forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the key pair, which is the same as `key_pair_name`.
* `key_pair_name` - The name of the key pair.
* `finger_print` - The fingerprint of the key pair.
* `region` - The region of the key pair.
* `resource_group_name` - The resource group name of the key pair.
* `create_time` - The creation time of the key pair.
* `update_time` - The update time of the key pair.

> **NOTE:** The `resource_group` attribute is defined in the schema but is not currently populated by the API read operation.

## Import

BCMP Keypair can be imported using the key pair name, e.g.

```
$ terraform import alibabacloudstack_bcmp_keypair.example my-key-pair-name
```
