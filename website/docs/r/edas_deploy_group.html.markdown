---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_deploy_group"
sidebar_current: "docs-alibabacloudstack-resource-edas-deploy-group"
description: |-
  Provides an EDAS deploy group resource.
---

# alibabacloudstack\_edas\_deploy\_group

Provides an EDAS deploy group resource.




## Example Usage

Basic Usage

```
resource "alibabacloudstack_edas_deploy_group" "default" {
  app_id     = var.app_id
  group_name = var.group_name
}

```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the application that you want to deploy.
* `group_name` - (Required, ForceNew) The name of the instance group that you want to create. 

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<app_id>:<group_name>:<group_id>`.
* `group_type` - The type of the instance group that you want to create. Valid values: 0: Default group. 1: Phased release is disabled for traffic management. 2: Phased release is enabled for traffic management.

## Import

EDAS deploy group can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_edas_deploy_group.group app_id:group_name:group_id
```
