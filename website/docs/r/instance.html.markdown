---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance"
sidebar_current: "docs-alibabacloudstack-resource-instance"
description: |-
  Provides an ECS instance resource.
---

# alibabacloudstack\_instance

Provides a ECS instance resource.

## Example Usage

```

# Create a new ECS instance for VPC
resource "alibabacloudstack_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "${var.cidr_block}"
  enable_ipv6    = true
}

resource "alibabacloudstack_vswitch" "vsw" {
  vpc_id            = "${alibabacloudstack_vpc.vpc.id}"
  cidr_block        = "${var.cidr_block}"
  availability_zone = "${var.availability_zone}"
  ipv6_cidr_block   = "${var.ipv6_cidr_block}"
}

resource "alibabacloudstack_security_group" "group" {
  name   = "new-group"
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_instance" "instance" {
  image_id             = "ubuntu_18_04_64_20G_alibase_20190624.vhd"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_security_group.group.id]
  instance_name        = "test_apsara_instance"
  vswitch_id           = alibabacloudstack_vswitch.vsw.id
  ipv6_cidr_block      = "${var.ipv6_cidr_block}"
  enable_ipv6          = true
  ipv6_address_count   = 3
}


```


## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The Image to use for the instance. ECS instance's image can be replaced via changing 'image_id'. When it is changed, the instance will reboot to make the change take effect.
* `instance_type` - (Required) The type of instance to start. When it is changed, the instance will reboot to make the change take effect.
* `security_groups` - (Required)  A list of security group ids to associate with.
* `availability_zone` - (Optional) The Zone to start the instance in. It is ignored and will be computed when set `vswitch_id`.
* `instance_name` - (Optional) The name of the ECS. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. If not specified, 
Terraform will autogenerate a default name is `ECS-Instance`.
* `system_disk_category` - (Optional) Valid values are `ephemeral_ssd`, `cloud_efficiency`, `cloud_ssd`, `cloud_essd`, `cloud`. `cloud` only is used to some none I/O optimized instance. Default to `cloud_efficiency`.
* `system_disk_size` - (Optional) Size of the system disk, measured in GiB. Value range: [20, 500]. The specified value must be equal to or greater than max{20, Imagesize}. Default value: max{40, ImageSize}. ECS instance's system disk can be reset when replacing system disk. When it is changed, the instance will reboot to make the change take effect.
* `system_disk_name` - (Optional) Name of the system disk. The name must be 2 to 128 characters in length. It must start with a letter and can contain letters, digits, colons (:), underscores (_), and hyphens (-). It cannot start with http:// or https://. If not specified, this parameter is null. Default value: null
* `system_disk_description` - (Optional) Description of the system disk. The description must be 2 to 256 characters in length. It cannot start with http:// or https://. If not specified, this parameter is null. Default value: null
* `description` - (Optional) Description of the instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `internet_max_bandwidth_in` - (Optional) Maximum incoming bandwidth from the public network, measured in Mbps (Mega bit per second). Value range: [1, 200]. If this value is not specified, then automatically sets it to 200 Mbps.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 100]. Default to 0 Mbps.
* `host_name` - (Optional) Host name of the ECS, which is a string of at least two characters. “hostname” cannot start or end with “.” or “-“. In addition, two or more consecutive “.” or “-“ symbols are not allowed. On Windows, the host name can contain a maximum of 15 characters, which can be a combination of uppercase/lowercase letters, numerals, and “-“. The host name cannot contain dots (“.”) or contain only numeric characters. When it is changed, the instance will reboot to make the change take effect.
On other OSs such as Linux, the host name can contain a maximum of 30 characters, which can be segments separated by dots (“.”), where each segment can contain uppercase/lowercase letters, numerals, or “_“. When it is changed, the instance will reboot to make the change take effect.
* `password` - (Optional, Sensitive) Password to an instance is a string of 8 to 30 characters. It must contain uppercase/lowercase letters and numerals, but cannot contain special symbols. When it is changed, the instance will reboot to make the change take effect./if you want to use random password See [Random Password](random_password.html.markdown).
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to an instance. If the `password` is filled in, this field will be ignored. When it is changed, the instance will reboot to make the change take effect.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating an instance with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set. When it is changed, the instance will reboot to make the change take effect.

* `vswitch_id` - (Optional) The virtual switch ID to launch in VPC. This parameter must be set unless you can create classic network instances. When it is changed, the instance will reboot to make the change take effect.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `user_data` - (Optional) User-defined data to customize the startup behaviors of an ECS instance and to pass data into an ECS instance. If updated, the instance will reboot to make the change take effect. Note: Not all of changes will take effect and it depends on [cloud-init module type](https://cloudinit.readthedocs.io/en/latest/topics/modules.html).
* `key_name` - (Optional, Force new resource) The name of key pair that can login ECS instance successfully without password. If it is specified, the password would be invalid.
* `role_name` - (Optional, Force new resource) Instance RAM role name. The name is provided and maintained by RAM. You can use `alibabacloudstack_ram_role` to create a new one.
* `private_ip` - (Optional) Instance private IP address can be specified when you creating new instance. It is valid when `vswitch_id` is specified. When it is changed, the instance will reboot to make the change take effect.
    Default to NoSpot. Note: Currently, the spot instance only supports domestic site account.
    Default to false.
* `security_enhancement_strategy` - (Optional, ForceNew) The security enhancement strategy.
    - Active: Enable security enhancement strategy, it only works on system images.
    - Deactive: Disable security enhancement strategy, it works on all images.
* `data_disks` - (Optional, ForceNew) The list of data disks created with instance.
    * `name` - (Optional, ForceNew) The name of the data disk.
    * `size` - (Required, ForceNew) The size of the data disk.
        - cloud：[5, 2000]
        - cloud_efficiency：[20, 32768]
        - cloud_ssd：[20, 32768]
        - cloud_essd：[20, 32768]
        - ephemeral_ssd: [5, 800]
    * `category` - (Optional, ForceNew) The category of the disk:
        - `cloud`: The general cloud disk.
        - `cloud_efficiency`: The efficiency cloud disk.
        - `cloud_ssd`: The SSD cloud disk.
        - `ephemeral_ssd`: The local SSD disk.
        Default to `cloud_efficiency`.
    * `encrypted` -(Optional, Bool, ForceNew) Encrypted the data in this disk.

        Default to false
    * `snapshot_id` - (Optional, ForceNew) The snapshot ID used to initialize the data disk. If the size specified by snapshot is greater that the size of the disk, use the size specified by snapshot as the size of the data disk.
    * `description` - (Optional, ForceNew) The description must be 2 to 256 characters in length.
    * `delete_with_instance` - (Optional, ForceNew) Delete this data disk when the instance is destroyed. It only works on cloud, cloud_efficiency, cloud_essd, cloud_ssd disk. If the category of this data disk was ephemeral_ssd, please don't set this param.

        Default to true
* `enable_ipv6` - (Optional, ForceNew) Specifies whether to enable the IPv6 block. Valid values: `false` (Default): disables IPv6 blocks. `true`: enables IPv6 blocks. 
* `ipv6_cidr_block` - (Optional) The ipv6 cidr block of VPC.
* `ipv6_address_count` - (Optional) The count of ipv6_address requested for allocation. If `enable_ipv6` is true. `ipv6_address_count` must be greater than 0.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the instance (until it reaches the initial `Running` status). 
`Note`: There are extra at most 2 minutes used to retry to aviod some needless API errors and it is not in the timeouts configure.
* `update` - (Defaults to 10 mins) Used when stopping and starting the instance when necessary during update - e.g. when changing instance type, password, image, vswitch and private IP.
* `delete` - (Defaults to 20 mins) Used when terminating the instance. `Note`: There are extra at most 5 minutes used to retry to aviod some needless API errors and it is not in the timeouts configure.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `status` - The instance status.
* `private_ip` - The instance private ip.

