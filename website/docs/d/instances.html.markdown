---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instances"
sidebar_current: "docs-alibabacloudstack-datasource-instances"
description: |-
    Provides a list of ECS instances to the user.
---

# alibabacloudstack\_instances

The Instances data source list ECS instance resources according to their ID, name regex, image id, status and other fields.

## Example Usage

```
data "alibabacloudstack_instances" "instances_ds" {
  name_regex = "web_server"
  status     = "Running"
}

output "first_instance_id" {
  value = "${data.alibabacloudstack_instances.instances_ds.instances.id}"
}

output "instance_ids" {
  value = "${data.alibabacloudstack_instances.instances_ds.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ECS instance IDs.
* `name_regex` - (Optional) A regex string to filter results by instance name.
* `image_id` - (Optional) The image ID of some ECS instance used.
* `status` - (Optional) Instance status. Valid values: "Creating", "Starting", "Running", "Stopping" and "Stopped". If undefined, all statuses are considered.
* `vpc_id` - (Optional) ID of the VPC linked to the instances.
* `vswitch_id` - (Optional) ID of the VSwitch linked to the instances.
* `availability_zone` - (Optional) Availability zone where instances are located.
* `tags` - (Optional) A map of tags assigned to the ECS instances. It must be in the format:
```
  data "alicloud_instances" "taggedInstances" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
```
* `ram_role_name` - (Optional, ForceNew) The RAM role name which the instance attaches.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of ECS instance IDs.
* `names` - A list of instances names. 
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `region_id` - Region ID the instance belongs to.
  * `availability_zone` - Availability zone the instance belongs to.
  * `status` - Instance current status.
  * `name` - Instance name.
  * `description` - Instance description.
  * `instance_type` - Instance type.
  * `instance_charge_type` - Instance charge type.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the VSwitch the instance belongs to.
  * `image_id` - Image ID the instance is using.
  * `private_ip` - Instance private IP address.
  * `eip` - EIP address the VPC instance is using.
  * `security_groups` - List of security group IDs the instance belongs to.
  * `key_name` - Key pair the instance is using.
  * `creation_time` - Instance creation time.
  * `internet_max_bandwidth_out` - Max output bandwidth for internet.
  * `tags` - A map of tags assigned to the ECS instance.
  * `disk_device_mappings` - Description of the attached disks.
    * `device` - Device information of the created disk: such as /dev/xvdb.
    * `size` - Size of the created disk.
    * `category` - Cloud disk category.
    * `type` - Cloud disk type: system disk or data disk.
  * `ram_role_name` - The Ram role name.
