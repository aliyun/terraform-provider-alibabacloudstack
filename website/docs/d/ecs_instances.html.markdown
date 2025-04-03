---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-instances"
description: |- 
  Provides a list of ecs instances owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_instances
-> **NOTE:** Alias name has: `alibabacloudstack_instances`

This data source provides a list of ECS instances in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_instances" "example" {
  name_regex = "web_server"
  status     = "Running"
  vpc_id     = "vpc-1234567890abcdef"
}

output "first_instance_id" {
  value = "${data.alibabacloudstack_ecs_instances.example.instances.0.id}"
}

output "instance_ids" {
  value = "${data.alibabacloudstack_ecs_instances.example.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ECS instance IDs. If specified, the data source will only return instances with matching IDs.
* `name_regex` - (Optional) A regex string to filter results by instance name. This allows you to retrieve instances whose names match the specified pattern.
* `image_id` - (Optional) The image ID used by some ECS instances. Filtering can be done based on this parameter.
* `status` - (Optional) Instance status. Valid values include: "Creating", "Starting", "Running", "Stopping", and "Stopped". If not specified, all statuses are considered.
* `vpc_id` - (Optional) The ID of the VPC linked to the instances. Filtering can be done based on this parameter.
* `vswitch_id` - (Optional) The ID of the VSwitch linked to the instances. Filtering can be done based on this parameter.
* `availability_zone` - (Optional) The availability zone where the instances are located. Filtering can be done based on this parameter.
* `tags` - (Optional) A map of tags assigned to the ECS instances. Filtering can be done based on these tags. For example:
  ```hcl
  data "alibabacloudstack_ecs_instances" "taggedInstances" {
    tags = {
      Environment = "Production",
      Owner      = "TeamA"
    }
  }
  ```
* `ram_role_name` - (Optional) The RAM role name which the instance attaches. Filtering can be done based on this parameter.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ECS instance IDs.
* `names` - A list of instance names.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - The ID of the instance.
  * `region_id` - The region ID where the instance belongs.
  * `availability_zone` - The availability zone where the instance belongs.
  * `status` - The current status of the instance.
  * `name` - The name of the instance.
  * `description` - The description of the instance.
  * `instance_type` - The type of the instance.
  * `instance_charge_type` - The charge type of the instance.
  * `vpc_id` - The ID of the VPC the instance belongs to.
  * `vswitch_id` - The ID of the VSwitch the instance belongs to.
  * `image_id` - The image ID the instance is using.
  * `private_ip` - The private IP address of the instance.
  * `eip` - The EIP address the VPC instance is using.
  * `security_groups` - A list of security group IDs the instance belongs to.
  * `key_name` - The key pair the instance is using.
  * `creation_time` - The creation time of the instance.
  * `internet_max_bandwidth_out` - The maximum output bandwidth for the internet.
  * `tags` - A map of tags assigned to the ECS instance.
  * `disk_device_mappings` - Description of the attached disks.
    * `device` - The device information of the created disk, such as `/dev/xvdb`.
    * `size` - The size of the created disk.
    * `category` - The category of the cloud disk.
    * `type` - The type of the cloud disk: system disk or data disk.
  * `ram_role_name` - The RAM role name attached to the instance.