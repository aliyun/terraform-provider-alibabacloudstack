---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterface"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-networkinterface"
description: |- 
  Provides a ecs Networkinterface resource.
---

# alibabacloudstack_ecs_networkinterface
-> **NOTE:** This resource can also be referred to by the following aliases: `alibabacloudstack_network_interface`

Provides a ecs Networkinterface resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccecsnetwork_interface15831"
}

resource "alibabacloudstack_vpc" "vpc" {
  name       = "${var.name}"
  cidr_block = "10.0.0.0/16"
}

resource "alibabacloudstack_vswitch" "vsw" {
  name       = "${var.name}"
  vpc_id     = alibabacloudstack_vpc.vpc.id
  cidr_block = "10.0.0.0/24"
  availability_zone = "cn-beijing-b"
}

resource "alibabacloudstack_security_group" "secgroup" {
  name        = "${var.name}"
  description = "Security Group"
  vpc_id      = alibabacloudstack_vpc.vpc.id
}

resource "alibabacloudstack_ecs_networkinterface" "default" {
  network_interface_name = "${var.name}-eni"
  vswitch_id             = alibabacloudstack_vswitch.vsw.id
  security_groups        = [alibabacloudstack_security_group.secgroup.id]
  primary_ip_address     = "10.0.0.10"
  private_ips_count      = 2
  description            = "Test ENI"
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_name` - (Optional, Computed) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `vswitch_id` - (Required, ForceNew) The VSwitch ID to create the ENI in.
* `security_groups` - (Required) A list of security group IDs to associate with the ENI. At least one security group must be specified.
* `primary_ip_address` - (Optional, ForceNew, Computed) The primary private IP address of the ENI. If not specified, Alibaba Cloud will automatically assign one within the CIDR block of the VSwitch.
* `private_ips` - (Optional, Computed) List of secondary private IPs to assign to the ENI. Maximum of 10 secondary IPs. Do not use both `private_ips` and `private_ips_count` in the same ENI resource block.
* `private_ips_count` - (Optional, Computed) Number of secondary private IPs to assign to the ENI. Value must be between 0 and 10. Do not use both `private_ips` and `private_ips_count` in the same ENI resource block.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters. It cannot begin with http:// or https://. Default is an empty string.
* `tags` - (Optional) A mapping of tags to assign to the resource.

**Deprecated Arguments:**

* `name` - (Deprecated) Deprecated in favor of `network_interface_name`. This field will be removed in a future release.
* `private_ip` - (Deprecated, ForceNew) Deprecated in favor of `primary_ip_address`. This field will be removed in a future release. Changing this parameter will force a new resource to be created.
* `mac` - (Deprecated) Deprecated in favor of `mac_address`. This field will be removed in a future release.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ENI ID.
* `mac_address` - The MAC address of the ENI.
* `network_interface_name` - Name of the ENI.
* `primary_ip_address` - The primary private IP address of the ENI.
* `private_ips` - List of all private IP addresses assigned to the ENI, including the primary IP and any secondary IPs.
* `private_ips_count` - Total number of private IPs assigned to the ENI, including the primary IP and any secondary IPs.