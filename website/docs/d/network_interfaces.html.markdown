---
subcategory: "ECS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_network_interfaces"
sidebar_current: "docs-apsarastack-datasource-network-interfaces"
description: |-
  Provides a data source to get a list of elastic network interfaces according to the specified filters.
---

# apsarastack\_network_interfaces

Use this data source to get a list of elastic network interfaces according to the specified filters in an ApsaraStack account.

## Example Usage

```
resource "apsarastack_network_interface" "NetInterface" {
  name              = "net_interface"
  vswitch_id        = apsarastack_vswitch.vsw.id
  security_groups   = [apsarastack_security_group.secgroup.id]
  private_ip        = "192.168.0.2"
  private_ips_count = 1
  description = "Hello Network Interface"
}

data "apsarastack_network_interfaces" "NetInterfaces" {
  ids = [
    apsarastack_network_interface.NetInterface.id
  ]
  name_regex = apsarastack_network_interface.NetInterface.name
  vswitch_id = apsarastack_vswitch.vsw.id
}

output "eni_name" {
    value = "${data.apsarastack_network_interface.NetInterface.interfaces.0.name}"
}
```

###  Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of ENI IDs.
* `name_regex` - (Optional) A regex string to filter results by ENI name.
* `vswitch_id` - (Optional) The VSwitch ID linked to ENIs.
* `private_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID linked to ENIs.
* `name` - (Optional) The name of the ENIs.
* `type` - (Optional) The type of ENIs, Only support for "Primary" or "Secondary".
* `instance_id` - (Optional) The ECS instance ID that the ENI is attached to.
* `tags` - (Optional) A map of tags assigned to ENIs.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `interfaces` - A list of ENIs. Each element contains the following attributes:
    * `id` - ID of the ENI.
    * `status` - Current status of the ENI.
    * `vswitch_id` - ID of the VSwitch that the ENI is linked to.
    * `zone_id` - ID of the availability zone that the ENI belongs to.
    * `public_ip` - Public IP of the ENI.
    * `private_ip` - Primary private IP of the ENI.
    * `private_ips` - A list of secondary private IP address that is assigned to the ENI.
    * `security_groups` - A list of security group that the ENI belongs to.
    * `name` - Name of the ENI.
    * `description` - Description of the ENI.
    * `instance_id` - ID of the instance that the ENI is attached to.
    * `creation_time` - Creation time of the ENI.
    * `tags` - A map of tags assigned to the ENI.
