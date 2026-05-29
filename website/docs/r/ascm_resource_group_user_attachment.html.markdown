---
subcategory: "Application"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_resource_group_user_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ascm-resource-group-user-attachment"
description: |-
  Provides a ASCM resource group user attachment resource.
---

# alibabacloudstack_ascm_resource_group_user_attachment

Provides a ASCM resource group user attachment resource. This resource is used to attach a user to a resource group in ASCM.

-> **NOTE:** Modifying any argument will force a new resource to be created.

## Example Usage

```
resource "alibabacloudstack_ascm_resource_group_user_attachment" "default" {
  user_id = alibabacloudstack_ascm_user.user.user_id
  rg_id   = alibabacloudstack_ascm_resource_group.default.rg_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, ForceNew) The ID of the user to be attached to the resource group. Modifying this argument will force a new resource to be created.
* `rg_id` - (Optional, ForceNew) The ID of the resource group. Modifying this argument will force a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource group user attachment. The value follows format `<rg_id>:<user_id>`.

## Import

ASCM Resource Group User Attachment can be imported using the resource group ID and user ID separated by colon, e.g.

```
$ terraform import alibabacloudstack_ascm_resource_group_user_attachment.example 12345:67890
```
