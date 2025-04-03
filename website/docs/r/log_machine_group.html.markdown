---
subcategory: "Log Service (SLS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_log_machine_group"
sidebar_current: "docs-alibabacloudstack-resource-log-machine-group"
description: |-
  Provides a Alibabacloudstack log tail machine group resource.
---

# alibabacloudstack_log_machine_group

Log Service manages all the ECS instances whose logs need to be collected by using the Logtail client in the form of machine groups.
 [Refer to details](https://www.alibabacloud.com/help/doc-detail/28966.htm)

## Example Usage

Basic Usage

```
resource "alibabacloudstack_log_project" "example" {
  name        = "tf-log"
  description = "created by terraform"
}

resource "alibabacloudstack_log_machine_group" "example" {
  project       = alibabacloudstack_log_project.example.name
  name          = "tf-machine-group"
  identify_type = "ip"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.2"]
}
```


## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the machine group belongs. 
* `name` - (Required, ForceNew) The machine group name, which is unique in the same project. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log machine group. It formats of `<project>:<name>`.
* `project` - The project name.
* `name` - The machine group name.
* `identify_type` - The machine identification type.
* `identify_list` - The machine identification.
* `topic` - The machine group topic.

## Import

Log machine group can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_log_machine_group.example tf-log:tf-machine-group
```