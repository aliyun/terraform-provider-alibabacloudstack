---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_group"
sidebar_current: "docs-alibabacloudstack-resource-ons-group"
description: |-
  Provides a alibabacloudstack ONS Group resource.
---

# alibabacloudstack_ons_group

Provides an ONS group resource.


## Example Usage

Basic Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "group_id" {
  default = "GID-onsGroupDatasourceName"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = var.name
  remark = "Ons Instance"
}

resource "alibabacloudstack_ons_group" "default" {
  group_id = var.group_id
  instance_id = alibabacloudstack_ons_instance.default.id
  remark = "dafault_ons_group_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id` - (Required) Name of the group. Two groups on a single instance cannot have the same name. A `group_id` starts with "GID_" or "GID-", and contains letters, numbers, hyphens (-), and underscores (_).
* `remark` - (Optional) This attribute is a concise description of group. The length cannot exceed 256. 
* `read_enable` - (Optional) This attribute is used to set the message reading enabled or disabled. It can only be set after the group is used by the client.

## Attributes Reference

The following attributes are exported:

* `id` - GroupID and InstanceID of the ONS Group. The value is in format `GroupID:InstanceID`.