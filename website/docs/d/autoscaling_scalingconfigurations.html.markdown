---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingconfigurations"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalingconfigurations"
description: |- 
  Provides a list of autoscaling scalingconfigurations owned by an alibabacloudstack account.
---

# alibabacloudstack_autoscaling_scalingconfigurations
-> **NOTE:** Alias name has: `alibabacloudstack_ess_scaling_configurations`

This data source provides a list of autoscaling scalingconfigurations in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_scalingconfigurations" "example" {
  scaling_group_id = "sg-1234567890abcdef"
  ids              = ["sc-abcdefgh12345678", "sc-ijklmnop90123456"]
  name_regex       = "scaling_configuration_example_.*"

  output_file = "scaling_configurations_output.txt"
}

output "first_scaling_configuration_id" {
  value = data.alibabacloudstack_autoscaling_scalingconfigurations.example.configurations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional, ForceNew) The ID of the scaling group to which the scaling configurations belong.
* `name_regex` - (Optional, ForceNew) A regex string to filter resulting scaling configurations by name.
* `ids` - (Optional) A list of scaling configuration IDs to filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling configuration IDs.
* `names` - A list of scaling configuration names.
* `configurations` - A list of scaling configurations. Each element contains the following attributes:
  * `id` - The ID of the scaling configuration.
  * `name` - The name of the scaling configuration.
  * `scaling_group_id` - The ID of the scaling group to which the scaling configuration belongs.
  * `image_id` - The image file ID used when creating ECS instances with this scaling configuration.
  * `instance_type` - The specification of the ECS instance created using this scaling configuration.
  * `security_group_id` - The ID of the security group to which the ECS instance belongs. Instances within the same security group can access each other.
  * `internet_max_bandwidth_in` - The Internet inbound bandwidth value in Mbps (Mega bit per second), with a value range of 1~200. If not specified, it will be automatically set to 200Mbps.
  * `internet_max_bandwidth_out` - The Internet outbound bandwidth value in Mbps for the ECS instance.
  * `system_disk_category` - The category of the system disk used in the scaling configuration.
  * `system_disk_size` - The size of the system disk in GB.
  * `data_disks` - A list of data disks configured in the scaling configuration. Each data disk has the following attributes:
    * `size` - The size of the data disk in GB.
    * `category` - The category of the data disk.
    * `snapshot_id` - The snapshot ID used to create the data disk.
    * `device` - The device attribute of the data disk.
    * `delete_with_instance` - Whether the data disk is deleted along with the instance when the instance is terminated.
  * `lifecycle_state` - The lifecycle state of the scaling configuration.
  * `creation_time` - The creation time of the scaling configuration.