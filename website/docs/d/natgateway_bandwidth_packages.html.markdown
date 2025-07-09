---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_bandwidthpackages"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-bandwidthpackages"
description: |-
  Provides a list of natgateway bandwidthpackages owned by an alibabacloudstack account.
---

# alibabacloudstack\_natgateway\_bandwidthpackages

This data source provides a list of natgateway bandwidthpackages in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccNatGatewaysBandwidthPackagesDatasource19756"
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
    name = "${var.name}"
	bandwidth = "5"
    natgateway_id = "${alibabacloudstack_nat_gateway.default.id}"
	description = "${var.name}"
	ip_count = "2"
}

data "alibabacloudstack_natgateway_bandwidth_packages" "default" {
	name_regex = "${alibabacloudstack_natgateway_bandwidth_package.default.name}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of natgateway bandwidthpackage IDs.
  * `name_regex` - (Optional) - A name Regex of natgateway bandwidthpackage.
  * `description_regex` - (Optional) - A description Regex of natgateway bandwidthpackage.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `bandwidth_packages` - A list of natgateway bandwidthpackage.
    * `id` - The ID of the natgateway bandwidthpackage.
    * `bandwidth` - The peak bandwidth of the shared bandwidth. Unit: Mbps.
    * `bandwidth_package_id` - The ID of the Internet shared bandwidth.
    * `business_status` - The status of the Internet Shared Bandwidth instance. Value:-**Normal**: Normal.-**Financialized**: Arrears.-**Unactivated**: not activated.
    * `creation_time` - The create time of the shared bandwidth.
    * `name` - The name of the shared bandwidth.
    * `description` - The description of the shared bandwidth.
    * `instance_charge_type` - The billing type of the Internet Shared Bandwidth instance. Value:<props = "china">**PostPaid**: Pay-as-you-go. </props><props = "china">**PrePaid**: Package year and month. </props><props = "intl">**PostPaid**: Pay-as-you-go. </props>
    * `internet_charge_type` - The billing method of the shared bandwidth. Value:-**PayByBandwidth**: Pay by bandwidth.-**PayBy95**: Pay by 95peak.
    * `natgateway_id` - The ID of the NAT gateway bound to the shared bandwidth.
    * `ip_count` - The number of EIPs bound to the shared bandwidth.
    * `isp` - Line type, value:-**BGP**:BGP (multi-line) line.-**BGP_PRO**:BGP (multi-line) boutique line.
    * `public_ip_addresses` - The public IP address of the Internet shared bandwidth instance.
    * `status` - The status of the Internet Shared Bandwidth instance. Default value: **Available * *.
