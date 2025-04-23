---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_loadbalancer"
sidebar_current: "docs-Alibabacloudstack-slb-loadbalancer"
description: |- 
  Provides a slb Loadbalancer resource.
---

# alibabacloudstack_slb_loadbalancer
-> **NOTE:** Alias name has: `alibabacloudstack_slb`

Provides a slb Loadbalancer resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccslbload_balancer19164"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name        = "${var.name}_vsw"
  vpc_id      = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_slb_loadbalancer" "default" {
  address_type = "intranet"
  name         = "rdk_test_name"
  specification = "slb.s1.small"
  vswitch_id   = "${alibabacloudstack_vpc_vswitch.default.id}"

  tags = {
    Environment = "Test"
    CreatedBy   = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SLB. This name must be unique within your AlibabaCloudStack account, can have a maximum of 80 characters, must contain only alphanumeric characters or hyphens (`-`), and must not begin or end with a hyphen. If not specified, Terraform will autogenerate a name beginning with `tf-lb`.
* `address_type` - (Optional, ForceNew) The network type of the SLB instance. Valid values: ["internet", "intranet"]. If the load balancer is launched in VPC, this value must be "intranet".
  * `internet`: After an Internet SLB instance is created, the system allocates a public IP address so that the instance can forward requests from the Internet.
  * `intranet`: After an intranet SLB instance is created, the system allocates an intranet IP address so that the instance can only forward intranet requests.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The ID of the VSwitch to launch the SLB in. If `address_type` is set to "internet", this field will be ignored.
* `specification` - (Optional) The specification of the Server Load Balancer instance. Default to an empty string indicating it is a "Shared-Performance" instance. Valid values include:
  * `slb.s1.small`
  * `slb.s2.small`
  * `slb.s2.medium`
  * `slb.s3.small`
  * `slb.s3.medium`
  * `slb.s3.large`
  * `slb.s4.large`
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `address` - (Optional, ForceNew) The service address of the load balancing instance. This field is automatically assigned by the system based on the `address_type`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.
* `address_type` - The address type of the load balancing instance.
* `specification` - The specification of the Server Load Balancer instance.
* `vswitch_id` - The ID of the VSwitch associated with the load balancer.