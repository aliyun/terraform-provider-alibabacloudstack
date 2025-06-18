---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_dhcpoptionsset"
sidebar_current: "docs-Alibabacloudstack-vpc-dhcpoptionsset"
description: |-
  Provides a vpc Dhcpoptionsset resource.
---

# alibabacloudstack\_vpc\_dhcpoptionsset

Provides a vpc Dhcpoptionsset resource.

## Example Usage
```
variable "name" {
  default = "tf-testAccVpcDhcpoptionssetBasic"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default2" {

  name = "${var.name}2"
  cidr_block = "172.16.0.0/16"
}



resource "alibabacloudstack_vpc_dhcp_options_set" "default" {
  domain_name = "aliyun.com"
  domain_name_servers = "10.82.0.12,10.82.0.13"
  associate_vpcs = [
                     "${alibabacloudstack_vpc.default0.id}",
                     "${alibabacloudstack_vpc.default1.id}",
                     "${alibabacloudstack_vpc.default2.id}"
                   ]
  dhcp_options_set_name = "${var.name}"
  dhcp_options_set_description = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:
  * `associate_vpcs` - (Optional) - A list of DhcpOptionsSet associated vpc IDs.
  * `dhcp_options_set_description` - (Optional) - The description can be blank or contain 1 to 256 characters. It must start with a letter or Chinese character but cannot start with http:// or https://.
  * `dhcp_options_set_name` - (Optional) - The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.
  * `domain_name` - (Optional) - The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
  * `domain_name_servers` - (Optional) - The DNS server IP addresses. Up to four DNS server IP addresses can be specified. IP addresses must be separated with commas (,).Before you specify any DNS server IP address, all ECS instances in the associated VPC network use the IP addresses of the Alibaba Cloud DNS servers, which are 100.100.2.136 and 100.100.2.138.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `associate_vpcs` - Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.
  * `dhcp_options_set_id` - The first ID of the resource
  * `status` - The status of the resource.
