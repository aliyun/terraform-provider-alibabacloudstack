---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairs"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-keypairs"
description: |- 
  Provides a list of ecs keypairs owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_keypairs
-> **NOTE:** Alias name has: `alibabacloudstack_key_pairs`

This data source provides a list of ECS key pairs in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
# Declare the resource
resource "alibabacloudstack_key_pair" "default" {
  key_name = "exampleKeyPair"
}

# Retrieve all key pairs matching the name_regex
data "alibabacloudstack_ecs_keypairs" "default" {
  name_regex = "${alibabacloudstack_key_pair.default.key_name}"
}

output "key_pairs" {
  value = data.alibabacloudstack_ecs_keypairs.default.key_pairs
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter the resulting key pairs by their names.
* `ids` - (Optional) A list of key pair IDs. If provided, only the key pairs with these IDs will be returned.
* `finger_print` - (Optional) The fingerprint of the key pair. The message-digest algorithm 5 (MD5) is used based on the public key fingerprint format defined in RFC 4716. For more information, see [RFC 4716](https://tools.ietf.org/html/rfc4716).
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of key pair names.
* `key_pairs` - A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_name` - Name of the key pair.
  * `finger_print` - The fingerprint of the key pair. The message-digest algorithm 5 (MD5) is used based on the public key fingerprint format defined in RFC 4716. For more information, see [RFC 4716](https://tools.ietf.org/html/rfc4716).
  * `instances` - A list of ECS instances that have been bound to this key pair. Each instance includes the following attributes:
    * `availability_zone` - The ID of the availability zone where the ECS instance is located.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vswitch_id` - The ID of the VSwitch attached to the ECS instance.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
  * `tags` - (Optional) A mapping of tags assigned to the key pair.