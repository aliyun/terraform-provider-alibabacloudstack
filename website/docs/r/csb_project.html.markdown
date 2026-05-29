---
subcategory: "Cloud Service Bus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_csb_project"
sidebar_current: "docs-Alibabacloudstack-resource-csb-project"
description: |-
  Provides a Alibabacloudstack resource to manage CSB Project.
---

# alibabacloudstack_csb_project

This resource will help you to manage CSB Project.

For information about CSB Project and how to use it, see [Create a Project](https://help.aliyun.com/apsara/enterprise/v_3_18_0_30393230/csb/apsarastack-developer-guide/obtains-information-about-a-single-service-group.html?spm=a2c4g.14484438.10001.97)

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_csb_project" "project" {
  csb_id            = "your-csb-id"
  project_name      = "example-project"
  owner_name        = "project-owner"
  owner_email       = "owner@example.com"
  owner_phone_num   = "13800138000"
  description       = "Example CSB project"
}
```

## Argument Reference

The following arguments are supported:

* `csb_id` - (Required, ForceNew) The ID of the CSB instance. Modifying this parameter will force the creation of a new resource.
* `project_name` - (Required) The name of the CSB project. Length constraint: 1 to 128 characters.
* `owner_name` - (Required) The name of the project owner.
* `owner_email` - (Optional) The email address of the project owner.
* `owner_phone_num` - (Optional) The phone number of the project owner.
* `description` - (Optional) The description of the project.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier of the resource, in the format `csb_id:project_name`.
* `csb_id` - The ID of the CSB instance.
* `project_name` - The name of the CSB project.
* `project_id` - The internal ID of the CSB project.
* `owner_name` - The name of the project owner.
* `owner_email` - The email address of the project owner.
* `owner_phone_num` - The phone number of the project owner.
* `description` - The description of the project.
* `owner_id` - The owner ID of the CSB project.
* `api_num` - The number of APIs published in the CSB project.

## Import

CSB Project can be imported using the combination of `csb_id` and `project_name`, separated by a colon, e.g.

```
$ terraform import alibabacloudstack_csb_project.example <csb_id>:<project_name>
```
