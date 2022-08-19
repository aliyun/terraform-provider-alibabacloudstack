---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_key_pairs"
sidebar_current: "docs-apsarastack-datasource-key-pairs"
description: |-
    Provides a list of available key pairs that can be used by an Apsarastack Cloud account.
---

# apsarastack\_key\_pairs

This data source provides a list of key pairs in an Apsarastack Cloud account according to the specified filters.

## Example Usage

```
# Declare the data source
resource "apsarastack_key_pair" "default" {
  key_name = "keyPairDatasource"
}

data "apsarastack_key_pairs" "default" {
  name_regex = "${apsarastack_key_pair.default.key_name}"
}

output "key_pairs" {
  value=data.apsarastack_key_pairs.default.*
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the resulting key pairs.
* `ids` - (Optional) A list of key pair IDs.
* `finger_print` - (Optional) A finger print used to retrieve specified key pair.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of key pair names.
* `key_pairs` - A list of key pairs. Each element contains the following attributes:
  * `id` - ID of the key pair.
  * `key_name` - Name of the key pair.
  * `finger_print` - Finger print of the key pair.
  * `instances` - A list of ECS instances that has been bound this key pair.
    * `availability_zone` - The ID of the availability zone where the ECS instance is located.
    * `instance_id` - The ID of the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vswitch_id` - The ID of the VSwitch attached to the ECS instance.
    * `public_ip` - The public IP address or EIP of the ECS instance.
    * `private_ip` - The private IP address of the ECS instance.
  * `tags` - (Optional) A mapping of tags to assign to the resource.