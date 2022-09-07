---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_key_pair_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ecs-key-pair-attachment"
description: |-
  Provides a Alibabacloudstack ECS Key Pair Attachment resource.
---


# alibabacloudstack\_ecs\_key\_pair\_attachment

Provides a ECS Key Pair Attachment resource.

For information about ECS Key Pair Attachment and how to use it, see [What is Key Pair Attachment](https://help.aliyun.com/apsara/enterprise/v_3_16_0_20220117/ecs/enterprise-developer-guide/RebootInstance.html?spm=a2c4g.14484438.10001.266).

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_ecs_key_pair_attachment" "example" {
  key_pair_name = "key_pair_name"
  instance_ids  = [i-gw80pxxxxxxxxxx]
}

```

## Argument Reference

The following arguments are supported:

* `key_pair_name` - (Required, ForceNew) The name of key pair used to bind.
* `force` - (Optional, ForceNew) Set it to true and it will reboot instances which attached with the key pair to make key pair affect immediately.
* `instance_ids` - (Required, ForceNew) The list of ECS instance's IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Key Pair Attachment. The value is formatted `<key_pair_name>:<instance_ids>`.

## Import

ECS Key Pair Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_key_pair_attachment.example <key_pair_name>:<instance_ids>
```
