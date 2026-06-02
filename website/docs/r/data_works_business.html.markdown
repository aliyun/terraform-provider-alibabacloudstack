---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_business"
description: |-
  Provides a DataWorks Business resource.
---

# alibabacloudstack_data_works_business

Provides a DataWorks Business resource.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_business" "example" {
  project_id  = "12345"
  name        = "tf-testaccdataworksbusiness"
  description = "This is a test business"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the business.

* `project_id` - (Optional, ForceNew) The ID of the DataWorks project where the business will be created. Changing this parameter will force a new resource to be created.

* `description` - (Optional) The description of the business.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the business resource. The value is formatted as `<project_id>:<business_id>`.

* `business_id` - The unique identifier of the business within the DataWorks project.

## Import

DataWorks Business can be imported using the project_id and business_id separated by a colon, e.g.

```
$ terraform import alibabacloudstack_data_works_business.example 12345:67890
```
