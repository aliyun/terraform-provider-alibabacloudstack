---
subcategory: "Datahub Service"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_datahub_project"
sidebar_current: "docs-apsarastack-resource-datahub-project"
description: |-
  Provides a Apsarastack datahub project resource.
---

# apsarastack\_datahub\_project

The project is the basic unit of resource management in Datahub Service and is used to isolate and control resources. It contains a set of Topics. You can manage the datahub sources of an application by using projects. [Refer to details](https://help.aliyun.com/document_detail/47440.html).

-> **NOTE:** Currently Datahub service only can be supported in the regions: cn-beijing, cn-hangzhou, cn-shanghai, cn-shenzhen,  ap-southeast-1.

## Example Usage

Basic Usage

```
resource "apsarastack_datahub_project" "example" {
  name    = "tf_datahub_project"
  comment = "created by terraform"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the datahub project. Its length is limited to 3-32 and only characters such as letters, digits and '_' are allowed. It is case-insensitive.
* `comment` - (Optional) Comment of the datahub project. It cannot be longer than 255 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the datahub project. It is the same as its name.
* `create_time` - Create time of the datahub project. It is a human-readable string rather than 64-bits UTC.
* `last_modify_time` - Last modify time of the datahub project. It is the same as *create_time* at the beginning. It is also a human-readable string rather than 64-bits UTC.

## Import

Datahub project can be imported using the *name* or ID, e.g.

```
$ terraform import apsarastack_datahub_project.example tf_datahub_project
```
