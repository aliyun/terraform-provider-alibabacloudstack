---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_snatentry"
sidebar_current: "docs-Alibabacloudstack-natgateway-snatentry"
description: |- 
  Provides a natgateway Snatentry resource.
---

# alibabacloudstack_natgateway_snatentry
-> **NOTE:** Alias name has: `alibabacloudstack_snat_entry`

Provides a natgateway Snatentry resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccnat_gatewaysnat_entry66949"
}

data "alibabacloudstack_zones" "default" {
    available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
    name       = "${var.name}"
    cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
    vpc_id            = "${alibabacloudstack_vpc.default.id}"
    cidr_block        = "172.16.0.0/21"
    availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
    name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
    vpc_id        = "${alibabacloudstack_vpc.default.id}"
    specification = "Small"
    name          = "${var.name}"
}

resource "alibabacloudstack_eip" "default" {
    name = "${var.name}"
}

resource "alibabacloudstack_eip_association" "default" {
    allocation_id = "${alibabacloudstack_eip.default.id}"
    instance_id   = "${alibabacloudstack_nat_gateway.default.id}"
}

resource "alibabacloudstack_natgateway_snatentry" "default" {
    snat_table_id     = "${alibabacloudstack_nat_gateway.default.snat_table_ids}"
    source_vswitch_id = "${alibabacloudstack_vswitch.default.id}"
    snat_ip           = "${alibabacloudstack_eip.default.ip_address}"
    source_cidr       = "${alibabacloudstack_vswitch.default.cidr_block}"
}
```

## Argument Reference

The following arguments are supported:

* `snat_table_id` - (Required, ForceNew) The SNAT table ID to which the SNAT entry belongs. This value can be obtained from the `snat_table_ids` attribute of the `alibabacloudstack_nat_gateway` resource.
* `source_vswitch_id` - (Optional, ForceNew) The ID of the VSwitch associated with the SNAT entry. This parameter is mutually exclusive with the `source_cidr` parameter.
* `snat_ip` - (Required) The public IP address used for the SNAT entry. This IP must belong to the Elastic IP (EIP) associated with the NAT Gateway.
* `source_cidr` - (Optional, ForceNew) The private network segment of the ECS instances. This parameter is mutually exclusive with the `source_vswitch_id` parameter.
* `snat_entry_id` - (Required) The unique identifier of the SNAT entry on the server.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the SNAT entry. The format is `<snat_table_id>:<snat_entry_id>`.
* `snat_entry_id` - The unique identifier of the SNAT entry on the server.