---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypair"
sidebar_current: "docs-Alibabacloudstack-ecs-keypair"
description: |- 
  Provides a ecs Keypair resource.
---

# alibabacloudstack_ecs_keypair
-> **NOTE:** Alias name has: `alibabacloudstack_key_pair`

Provides a ecs Keypair resource.

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_ecs_keypair" "basic" {
  key_pair_name = "terraform-test-key-pair"
}

// Using name prefix to build key pair
resource "alibabacloudstack_ecs_keypair" "prefix" {
  key_name_prefix = "terraform-test-key-pair-prefix"
}

// Import an existing public key to build a key pair
resource "alibabacloudstack_ecs_keypair" "publickey" {
  key_pair_name = "my_public_key"
  public_key    = "ssh-rsa AB3Napapsod45678qwertyuudsfsg"
}
```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Required, ForceNew) Name of the key pair. The value contains 2 to 128 English or Chinese characters. The name must start with a letter or a Chinese character. It cannot start with `http://` or `https://`. The value can contain digits, colons (`:`), underscores (`_`), or dashes (`-`).
* `key_name_prefix` - (Optional, ForceNew) The prefix of the key pair name. It is conflicting with `key_pair_name`. If it is specified, Terraform will use it to generate a unique key pair name.
* `public_key` - (Optional, ForceNew) An existing public key that you want to import and manage using AlibabaCloudStack key pairs.
* `key_file` - (Optional, ForceNew) The file name where the private key of the newly created key pair will be saved. It is strongly recommended to specify this parameter when creating a key pair, as you will not be able to retrieve the private key afterward if it is not saved.
* `tags` - (Optional) A mapping of tags to assign to the resource.

> **NOTE:** If neither `key_pair_name` nor `key_name_prefix` is set, Terraform will generate a unique ID to replace them.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `key_pair_name` - The name of the key pair.
* `finger_print` - The fingerprint of the key pair. The message-digest algorithm 5 (MD5) is used based on the public key fingerprint format defined in RFC 4716. For more information, see [RFC 4716](https://tools.ietf.org/html/rfc4716).