---
subcategory: "API Gateway V2"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_cascade_instances"
description: |-
  Provides a list of Api Gateway V2 Cascade Instances to the user.
---

# alibabacloudstack_api_gateway_v2_cascade_instances

This data source provides the Api Gateway V2 Cascade Instances of the current Alibaba Cloud user.

## Example Usage

Basic Usage

```hcl
data "alibabacloudstack_api_gateway_v2_cascade_instances" "example" {
  
}

output "first_api_gateway_v2_cascade_instance_id" {
  value = "${data.alibabacloudstack_api_gateway_v2_cascade_instances.example.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Cascade Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by Cascade Instance name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of Cascade Instance IDs.
* `names` - A list of Cascade Instance names.
* `instances` - A list of Cascade Instances. Each element contains the following attributes:
  * `id` - The ID of the Cascade Instance.
  * `instance_type` - The type of the Cascade Instance.
  * `instance_name` - The name of the Cascade Instance.
  * `cascade_instance_id` - The cascade instance ID.
  * `create_time` - The creation time of the Cascade Instance.
  * `status` - The status of the Cascade Instance.