---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_networkinterfaces"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-networkinterfaces"
description: |- 
  Provides a list of ecs networkinterfaces owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_networkinterfaces
-> **NOTE:** Alias name has: `alibabacloudstack_network_interfaces`

This data source provides a list of ECS network interfaces (ENIs) in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
resource "alibabacloudstack_network_interface" "eni" {
  name              = "test-eni"
  vswitch_id        = alibabacloudstack_vswitch.vsw.id
  security_groups   = [alibabacloudstack_security_group.secgroup.id]
  private_ip        = "192.168.0.2"
  private_ips_count = 1
  description       = "Test Network Interface"
}

data "alibabacloudstack_ecs_networkinterfaces" "enis" {
  ids            = [alibabacloudstack_network_interface.eni.id]
  name_regex     = "test-eni"
  vpc_id         = alibabacloudstack_vpc.vpc.id
  vswitch_id     = alibabacloudstack_vswitch.vsw.id
  private_ip     = "192.168.0.2"
  security_group_id = alibabacloudstack_security_group.secgroup.id
  type           = "Secondary"
  instance_id    = alibabacloudstack_instance.ecs.id
  output_file    = "eni_list.txt"
}

output "eni_name" {
  value = data.alibabacloudstack_ecs_networkinterfaces.enis.interfaces.0.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ENI IDs.
* `name_regex` - (Optional) A regex string to filter results by ENI name.
* `vpc_id` - (Optional) The VPC ID linked to ENIs.
* `vswitch_id` - (Optional) The VSwitch ID linked to ENIs.
* `private_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID linked to ENIs.
* `type` - (Optional) The type of ENIs, only supports "Primary" or "Secondary".
* `instance_id` - (Optional) The ECS instance ID that the ENI is attached to.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ENI names.
* `interfaces` - A list of ENIs. Each element contains the following attributes:
  * `id` - ID of the ENI.
  * `status` - Current status of the ENI.
  * `vpc_id` - ID of the VPC that the ENI belongs to.
  * `vswitch_id` - ID of the VSwitch that the ENI is linked to.
  * `zone_id` - ID of the availability zone that the ENI belongs to.
  * `public_ip` - Public IP of the ENI.
  * `private_ip` - Primary private IP of the ENI.
  * `private_ips` - A list of secondary private IP addresses assigned to the ENI.
  * `mac` - MAC address of the ENI.
  * `security_groups` - A list of security groups that the ENI belongs to.
  * `name` - Name of the ENI.
  * `description` - Description of the ENI.
  * `instance_id` - ID of the instance that the ENI is attached to.
  * `creation_time` - Creation time of the ENI.
  * `tags` - A map of tags assigned to the ENI.