---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ecs_key_pair"
sidebar_current: "docs-apsarastack-resource-ecs-key-pair"
description: |-
  Provides a Apsarastack ECS Key Pair resource.
---

# apsarastack\_ecs\_key\_pair

Provides a ECS Key Pair resource.

For information about ECS Key Pair and how to use it, see [What is Key Pair](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/ecs/enterprise-developer-guide/CreateKeyPair-1.html?spm=a2c4g.14484438.10001.356).

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_key_pair" "default" {
	key_name ="tf-testAccKeyPairConfig4427256049161700561"
	public_key = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

```

## Argument Reference

The following arguments are supported:

* `key_file` - (Optional, ForceNew) The key file.
* `key_name` - (Optional, ForceNew) The key pair's name. It is the only in one Alicloud account.
* `key_name_prefix` - (Optional, ForceNew) The key pair name's prefix. It is conflict with `key_pair_name`. If it is specified, terraform will using it to build the only key name.
* `public_key` - (Optional, ForceNew) You can import an existing public key and using Alicloud key pair to manage it. If this parameter is specified, `resource_group_id` is the key pair belongs.
* `resource_group_id` - (Optional) The Id of resource group which the key pair belongs.
* `tags` - (Optional) A mapping of tags to assign to the resource.

-> **NOTE:** If `key_pair_name` and `key_name_prefix` are not set, terraform will produce a specified ID to replace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Key Pair. Value as `key_pair_name`.
* `finger_print` The finger print of the key pair.

## Import

ECS Key Pair can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_key_pair.example <key_name>
```
