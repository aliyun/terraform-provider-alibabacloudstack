---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_namespaces"
sidebar_current: "docs-alibabacloudstack-datasource-cr-ee-namespaces"
description: |-
  Provides a list of Container Registry Enterprise Edition namespaces.
---

# alibabacloudstack\_cr\_ee\_namespaces

This data source provides a list Container Registry Enterprise Edition namespaces on Alibaba Cloud.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_cr_ee_namespaces" "my_namespaces" {
  instance_id = "cri-xxx"
  name_regex  = "my-namespace"
  output_file = "my-namespace-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_ee_namespaces.my_namespaces.namespaces}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of Container Registry Enterprise Edition instance.
* `ids` - (Optional) A list of ids to filter results by namespace id.
* `name_regex` - (Optional) A regex string to filter results by namespace name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Container Registry Enterprise Edition namespaces. Its element is a namespace uuid.
* `names` - A list of namespace names.
* `namespaces` - A list of matched Container Registry Enterprise Edition namespaces. Each element contains the following attributes:
  * `instance_id` - ID of Container Registry Enterprise Edition instance.
  * `id` - ID of Container Registry Enterprise Edition namespace.
  * `name` - Name of Container Registry Enterprise Edition namespace.
  * `auto_create` - Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
  * `default_visibility` - `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

