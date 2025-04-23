---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance_attachments"
sidebar_current: "docs-Alibabacloudstack-datasource-ots-instance_attachments"
description: |- 
  Provides a list of ots instance attachments owned by an Alibabacloudstack account.
---

# alibabacloudstack_ots_instance_attachments

This data source provides a list of OTS instance attachments in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ots_instance_attachments" "example" {
  instance_name = "sample-instance"
  name_regex    = "testvpc"
  output_file   = "attachments.txt"
}

output "first_ots_attachment_id" {
  value = "${data.alibabacloudstack_ots_instance_attachments.example.attachments.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of the OTS instance.
* `name_regex` - (Optional) A regex string used to filter results by VPC name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of VPC names associated with the OTS instance.
* `vpc_ids` - A list of VPC IDs associated with the OTS instance.
* `attachments` - A list of instance attachments. Each element contains the following attributes:
  * `id` - The resource ID, which is the same as the `instance_name`.
  * `domain` - The domain of the instance attachment.
  * `endpoint` - The access endpoint of the instance attachment.
  * `region` - The region of the instance attachment.
  * `instance_name` - The name of the OTS instance.
  * `vpc_name` - The name of the VPC attached to the instance.
  * `vpc_id` - The ID of the VPC attached to the instance.