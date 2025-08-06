---
subcategory: "CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_ceninstance"
sidebar_current: "docs-Alibabacloudstack-cen-ceninstance"
description: |-
  Provides a cen Ceninstance resource.
---

# alibabacloudstack\_cen\_ceninstance

Provides a cen Ceninstance resource.

## Example Usage
```
variable "name" {
}


resource "alibabacloudstack_cen_instance" "default" {
  description = var.name
  cen_instance_name = var.name
}
```

## Argument Reference

The following arguments are supported:
  * `cen_instance_name` - (Optional) - cen instance name.
  * `description` - (Optional) - cen instance description.
  * `protection_level` - (Optional) - cen instance protection level.
  * `status` - (Optional) - cen instance status.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `cen_id` - cen instance id.
  * `create_time` - cen instance create time.
  * `status` - cen instance status.
