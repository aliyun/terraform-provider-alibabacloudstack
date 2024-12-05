---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb"
sidebar_current: "docs-alibabacloudstack-resource-slb"
description: |-
  Provides an Application Load Balancer resource.
---

# alibabacloudstack\_slb

Provides an Application Load Balancer resource.

## Example Usage

```
variable "name" {
  default = "terraformtestslbconfig"
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

resource "alibabacloudstack_slb" "default" {
  name          = "${var.name}"
  vswitch_id    = "${alibabacloudstack_vswitch.default.id}"
  specification = "slb.s2.small"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SLB. This name must be unique within your alibabacloudstack account, can have a maximum of 80 characters,
must contain only alphanumeric characters or hyphens, such as "-","/",".","_", and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `address_type` - (Optional, ForceNew) The network type of the SLB instance. Valid values: ["internet", "intranet"]. If load balancer launched in VPC, this value must be "intranet".
    - internet: After an Internet SLB instance is created, the system allocates a public IP address so that the instance can forward requests from the Internet.
    - intranet: After an intranet SLB instance is created, the system allocates an intranet IP address so that the instance can only forward intranet requests.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in. If `address_type` is internet, it will be ignore.
* `specification` - (Optional)  The specification of the Server Load Balancer instance. Default to empty string indicating it is "Shared-Performance" instance.it is must be specified and it valid values are: `slb.s1.small`, `slb.s2.small`, `slb.s2.medium`,
  `slb.s3.small`, `slb.s3.medium`, `slb.s3.large` and `slb.s4.large`.
* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.

