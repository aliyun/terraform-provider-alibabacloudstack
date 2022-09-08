---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ots-instance-attachment"
description: |-
  Provides an OTS (Open Table Service) resource to attach VPC to instance.
---

# alibabacloudstack\_ots\_instance\_attachment

This resource will help you to bind a VPC to an OTS instance.

## Example Usage

```
# Create an OTS instance
resource "alibabacloudstack_ots_instance" "foo" {
  name        = "my-ots-instance"
  description = "for table"
  accessed_by = "Vpc"
  tags = {
    Created = "TF"
    For     = "Building table"
  }
}

data "alibabacloudstack_zones" "foo" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  cidr_block = "172.16.0.0/16"
  name       = "for-ots-instance"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = alibabacloudstack_vpc.foo.id
  vswitch_name      = "for-ots-instance"
  cidr_block        = "172.16.1.0/24"
  zone_id           = data.alibabacloudstack_zones.foo.zones[0].id
}

resource "alibabacloudstack_ots_instance_attachment" "foo" {
  instance_name = alibabacloudstack_ots_instance.foo.name
  vpc_name      = "attachment1"
  vswitch_id    = alibabacloudstack_vswitch.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance.
* `vpc_name` - (Required, ForceNew) The name of attaching VPC to instance.
* `vswitch_id` - (Required, ForceNew) The ID of attaching VSwitch to instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID. The value is same as "instance_name".
* `instance_name` - The instance name.
* `vpc_name` - The name of attaching VPC to instance.
* `vswitch_id` - The ID of attaching VSwitch to instance.
* `vpc_id` - The ID of attaching VPC to instance.


