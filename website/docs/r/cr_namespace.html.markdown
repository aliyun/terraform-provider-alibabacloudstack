---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack:alibabacloudstack_cr_namespace"
sidebar_current: "docs-alibabacloudstack-resource-container-registry"
description: |-
  Provides a Alibabacloudstack resource to manage Container Registry namespaces.
---

# alibabacloudstack\_cr\_namespace

This resource will help you to manager Container Registry namespaces.


## Example Usage

Basic Usage

```
resource "alibabacloudstack_cr_namespace" "my-namespace" {
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of Container Registry namespace.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Container Registry namespace. The value is same as its name.

