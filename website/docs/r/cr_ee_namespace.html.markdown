---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_namespace"
sidebar_current: "docs-alibabacloudstack-resource-cr-ee-namespace"
description: |-
  Provides a Alibabacloudstack resource to manage Container Registry Enterprise Edition namespaces.
---

# alibabacloudstack\_cr\_ee\_namespace

This resource will help you to manager Container Registry Enterprise Edition namespaces.

For information about Container Registry Enterprise Edition namespaces and how to use it, see [Create a Namespace](https://www.alibabacloud.com/help/doc-detail/145483.htm)

-> **NOTE:** Available in v1.86.0+.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition instance.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition namespace. It can contain 2 to 30 characters.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - ID of Container Registry Enterprise Edition namespace. The value is in format `{instance_id}:{namespace}` .

## Import

Container Registry Enterprise Edition namespace can be imported using the `{instance_id}:{namespace}`, e.g.

```
$ terraform import alibabacloudstack_cr_ee_namespace.default cri-xxx:my-namespace
```
