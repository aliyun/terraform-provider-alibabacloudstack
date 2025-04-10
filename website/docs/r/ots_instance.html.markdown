---
subcategory: "OTS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance"
sidebar_current: "docs-Alibabacloudstack-ots-instance"
description: |- 
  Provides a ots Instance resource.
---

# alibabacloudstack_ots_instance

Provides a ots Instance resource.

## Example Usage

```hcl
# Create an OTS instance
resource "alibabacloudstack_ots_instance" "foo" {
  name          = "my-ots-instance"
  description   = "This is a test OTS instance"
  accessed_by   = "Vpc" # Options: Any, Vpc, ConsoleOrVpc. Default: Any
  instance_type = "Capacity" # Options: Capacity, HighPerformance. Default: HighPerformance
  tags = {
    Created = "TF"
    For     = "Building table"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the OTS instance. Changing this forces a new resource to be created.
* `accessed_by` - (Optional) The network limitation for accessing the OTS instance. Valid values:
  * `Any` - Allows all networks to access the instance.
  * `Vpc` - Only allows access from the attached VPC.
  * `ConsoleOrVpc` - Allows access from the web console or the attached VPC.
  
  Default value: `Any`.

* `instance_type` - (Optional, ForceNew) The type of the OTS instance. Valid values:
  * `Capacity` - Suitable for scenarios with large data volumes and high throughput requirements.
  * `HighPerformance` - Suitable for scenarios requiring lower latency and higher performance.
  
  Default value: `HighPerformance`. Changing this forces a new resource to be created.

* `description` - (Optional, ForceNew) A brief description of the OTS instance. This field cannot be modified after creation. Changing this forces a new resource to be created.
* `tags` - (Optional) A mapping of tags to assign to the OTS instance.
* `propreties` - (Optional) Additional properties for the OTS instance. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the OTS instance. It has the same value as the `name`.
* `name` - The name of the OTS instance.
* `description` - The description of the OTS instance.
* `accessed_by` - The network limitation for accessing the OTS instance.
* `instance_type` - The type of the OTS instance.
* `tags` - A mapping of tags assigned to the OTS instance.
* `propreties` - The computed properties of the OTS instance.

## Import

OTS instance can be imported using the instance id or name, e.g.

```bash
$ terraform import alibabacloudstack_ots_instance.foo "my-ots-instance"
```