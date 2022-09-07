---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_network_interface_attachment"
sidebar_current: "docs-alibabacloudstack-resource-network-interface-attachment"
description: |-
  Provides an alibabacloudstack ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.
---

# alibabacloudstack\_network\_interface\_attachment

Provides an alibabacloudstack ECS Elastic Network Interface Attachment as a resource to attach ENI to or detach ENI from ECS Instances.

## Example Usage

Bacis Usage

```

resource "alibabacloudstack_instance" "apsarainstance" {
  image_id              = "gj2j1g3-45h3nnc-454hj5g"
  instance_type        = "ecs.n4.large"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["Grp-1"]
  instance_name        = "apsarainstance"
  vswitch_id           = "vsw-abc1345"
}

resource "alibabacloudstack_network_interface" "NetInterface" {
  name              = "ENI"
  vswitch_id        = alibabacloudstack_instance.apsarainstance.vswitch_id
  security_groups   = alibabacloudstack_instance.apsarainstance.security_groups
  private_ip        = "192.168.0.2"
  private_ips_count = 3
}

resource "alibabacloudstack_network_interface_attachment" "NetIntAttachment" {
  count                = alibabacloudstack_network_interface.NetInterface.private_ips_count
  instance_id          = alibabacloudstack_instance.apsarainstance.id
  network_interface_id = alibabacloudstack_network_interface.NetInterface.*.id
}
```

## Argument Reference

The following argument are supported:

* `instance_id` - (Required, ForceNew) The instance ID to attach.
* `network_interface_id` - (Required, ForceNew) The ENI ID to attach.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource, formatted as `<network_interface_id>:<instance_id>`.
