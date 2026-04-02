---
subcategory: "BCMP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_keypair"
sidebar_current: "docs-Alibabacloudstack-bcmp-keypair"
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
  public_key    = "xx"
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Optional, ForceNew) Name of the key pair. The value contains 2 to 128 English or Chinese characters. The name must start with a letter or a Chinese character. It cannot start with `http://` or `https://`. The value can contain digits, colons (`:`), underscores (`_`), or dashes (`-`).
* `public_key` - (Optional, ForceNew) An existing public key that you want to import and manage using AlibabaCloudStack key pairs.
* `key_file` - (Optional, ForceNew) The file name where the private key of the newly created key pair will be saved. It is strongly recommended to specify this parameter when creating a key pair, as you will not be able to retrieve the private key afterward if it is not saved.

> **NOTE:** If `key_pair_name` is not set, Terraform will generate a unique ID to replace it.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_pair_name` - The name of the key pair.
* `finger_print` - The fingerprint of the key pair.
* `region` - The region of the key pair.
* `resource_group` - The resource group ID of the key pair.
* `resource_group_name` - The resource group name of the key pair.
* `create_time` - The creation time of the key pair.
* `update_time` - The update time of the key pair.
