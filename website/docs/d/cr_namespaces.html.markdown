---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_namespaces"
sidebar_current: "docs-alibabacloudstack-datasource-cr-namespaces"
description: |-
  Provides a list of Container Registry namespaces.
---

# alibabacloudstack_cr_namespaces

This data source provides a list Container Registry namespaces on Alibabacloudstack Cloud.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_cr_namespaces" "my_namespaces" {
  name_regex  = "my-namespace"
  output_file = "my-namespace-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_namespaces.my_namespaces.namespaces}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by namespace name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry namespaces. Its element is a namespace name.
* `names` - A list of namespace names.
* `namespaces` - A list of matched Container Registry namespaces. Each element contains the following attributes:
  * `name` - Name of Container Registry namespace.
  * `auto_create` - Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
  * `default_visibility` - `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.