---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_metering_query_ecs"
sidebar_current: "docs-alibabacloudstack-ascm-metering-query-ecs"
description: |-
  Provides a list of metering data for ECS instances.
---
# alibabacloudstack_ascm_metering_query_ecs

This data source provides a list of metering data for ECS instances.

## Example Usage

```hcl
data "alibabacloudstack_ascm_metering_query_ecs" "example" {
  start_time = "2023-01-01T00:00:00Z"
  end_time   = "2023-01-31T23:59:59Z"
  product_name = "ECS"
}

output "ecs_metering_data" {
  value = data.alibabacloudstack_ascm_metering_query_ecs.example.data
}
```

## Argument Reference
The following arguments are supported:

* `start_time` - (Required) The start time for the metering query in ISO 8601 format (e.g., "2023-01-01T00:00:00Z").
* `end_time` - (Required) The end time for the metering query in ISO 8601 format (e.g., "2023-01-31T23:59:59Z").
* `org_id` - (Optional) The organization ID for which to retrieve metering data.
* `product_name` - (Required) The name of the product for which to retrieve metering data (e.g., "ECS").
* `is_parent_id` - (Optional) Indicates whether the organization ID is a parent ID. 
* `ins_id` - (Optional) The instance ID for which to retrieve metering data. 
* `region` - (Optional) The region for which to retrieve metering data. 
* `resource_group_id` - (Optional) The resource group ID for which to retrieve metering data. 
* `name_regex` - (Optional) A regex pattern to filter the results by instance name.

## Attributes Reference
The following attributes are exported:

* `data` - A list of ECS metering data. Each element contains the following attributes:
    * `private_ip_address` - The private IP address of the ECS instance.
    * `instance_type_family` - The instance type family of the ECS instance.
    * `memory` - The memory size of the ECS instance in GB.
    * `cpu` - The number of CPUs in the ECS instance.
    * `os_name` - The operating system name of the ECS instance.
    * `org_name` - The name of the organization.
    * `instance_network_type` - The network type of the ECS instance.
    * `eip_address` - The EIP address of the ECS instance.
    * `resource_g_name` - The name of the resource group.
    * `instance_type` - The instance type of the ECS instance.
    * `status` - The status of the ECS instance.
    * `sys_disk_size` - The system disk size of the ECS instance in GB.
    * `gpu_amount` - The number of GPUs in the ECS instance.
    * `instance_name` - The name of the ECS instance.
    * `vpc_id` - The VPC ID of the ECS instance.
    * `start_time` - The start time of the metering data.
    * `end_time` - The end time of the metering data.
    * `create_time` - The creation time of the ECS instance.
    * `data_disk_size` - The data disk size of the ECS instance in GB.
    * `is_parent_id` - Indicates whether the organization ID is a parent ID. 
    * `ins_id` - The instance ID for which to retrieve metering data. 
    * `region` - The region for which to retrieve metering data. 
    * `resource_group_id` - The resource group ID for which to retrieve metering data. 