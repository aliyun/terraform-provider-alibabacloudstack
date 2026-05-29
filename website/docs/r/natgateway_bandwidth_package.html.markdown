---
subcategory: "Network Address Translation Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_bandwidthpackage"
sidebar_current: "docs-Alibabacloudstack-natgateway-bandwidthpackage"
description: |-
  Provides a natgateway Bandwidthpackage resource.
---

# alibabacloudstack\_natgateway\_bandwidthpackage

Provides a natgateway Bandwidthpackage resource.

## Example Usage
```
variable "name" {
	default = "tf-testaccnat_gatewaybandwi59570"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	name = "${var.name}"
}


resource "alibabacloudstack_natgateway_bandwidth_package" "default" {
  name = "tf-testaccnat_gatewaybandwi59570"
  bandwidth = "5"
  natgateway_id = "${alibabacloudstack_nat_gateway.default.id}"
  description = "tf-testaccnat_gatewaybandwi59570"
  ip_count = "2"
}
```

## Argument Reference

The following arguments are supported:
  * `name` - (Optional) - The name of the shared bandwidth.
  * `bandwidth` - (Required) - The peak bandwidth of the shared bandwidth. Unit: Mbps.
  * `natgateway_id` - (Required) - The ID of the NAT gateway.
  * `description` - (Optional) - The description of the shared bandwidth.
  * `ip_count` - (Required) - The number of EIPs that can be attached to the shared bandwidth.
  * `status` - (Optional) - The status of the Internet Shared Bandwidth instance. Default value: **Available * *.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `name` - The name of the shared bandwidth.
  * `bandwidth_package_id` - The ID of the Internet shared bandwidth.
  * `business_status` - The status of the Internet Shared Bandwidth instance. Value:-**Normal**: Normal.-**Financialized**: Arrears.-**Unactivated**: not activated.
  * `description` - The description of the shared bandwidth.
  * `instance_charge_type` - The billing type of the Internet Shared Bandwidth instance. Value:<props = "china">**PostPaid**: Pay-as-you-go. </props><props = "china">**PrePaid**: Package year and month. </props><props = "intl">**PostPaid**: Pay-as-you-go. </props>
  * `isp` - Line type, value:-**BGP**:BGP (multi-line) line.-**BGP_PRO**:BGP (multi-line) boutique line.
  * `public_ip_addresses` - The public IP address of the Internet shared bandwidth instance.
    * `allocation_id` - The ID of the instance of the public IP address.
    * `ip_address` - Public IP address.
  * `status` - The status of the Internet Shared Bandwidth instance. Default value: **Available * *.
