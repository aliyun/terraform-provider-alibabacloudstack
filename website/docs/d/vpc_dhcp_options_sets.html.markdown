---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_dhcp_options_sets"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-dhcpoptionssets"
description: |-
  Provides a list of vpc dhcpoptionssets owned by an alibabacloudstack account.
---

# alibabacloudstack\_vpc\_dhcp_options_sets

This data source provides a list of vpc dhcpoptionssets in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccVpcDhcpOptionsDataSource-2860522"
}

resource "alibabacloudstack_vpc" "default0" {

  name = "${var.name}0"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc" "default1" {

  name = "${var.name}1"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_dhcp_options_set" "default" {
  	dhcp_options_set_name =        "${var.name}"
	dhcp_options_set_description = "${var.name}"
	domain_name =                  "aliyun.com"
	domain_name_servers =         "10.82.0.12,10.82.0.13"
	associate_vpcs = [
		"${alibabacloudstack_vpc.default0.id}",
		"${alibabacloudstack_vpc.default1.id}"
	]
}

 

data "alibabacloudstack_vpc_dhcp_options_sets" "default" {
  description_regex = "${alibabacloudstack_vpc_dhcp_options_set.default.dhcp_options_set_description}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of dhcp_options_set IDs.
  * `name_regex` - (Optional) - A regex string to filter results by dhcp_options_set name.
  * `description_regex` - (Optional) - A regex string to filter results by dhcp_options_set_description.
  * `domain_name` - (Optional) - The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
  * `dhcp_options_set_name` - (Optional) - The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `dhcp_options_sets` - DHCP options set
    * `id` - The ID of the DHCP options set.
    * `associate_vpcs` - Field 'associate_vpcs' has been deprecated from provider version 1.153.0 and it will be removed in the future version. Please use the new resource 'alicloud_vpc_dhcp_options_set_attachment' to attach DhcpOptionsSet and Vpc.
    * `dhcp_options_set_description` - The description can be blank or contain 1 to 256 characters. It must start with a letter or Chinese character but cannot start with http:// or https://.
    * `dhcp_options_set_id` - The first ID of the resource
    * `dhcp_options_set_name` - The name must be 2 to 128 characters in length and can contain letters, Chinese characters, digits, underscores (_), and hyphens (-). It must start with a letter or a Chinese character.
    * `domain_name` - The root domain, for example, example.com. After a DHCP options set is associated with a Virtual Private Cloud (VPC) network, the root domain in the DHCP options set is automatically synchronized to the ECS instances in the VPC network.
    * `domain_name_servers` - The DNS server IP addresses. Up to four DNS server IP addresses can be specified. IP addresses must be separated with commas (,).Before you specify any DNS server IP address, all ECS instances in the associated VPC network use the IP addresses of the Alibaba Cloud DNS servers, which are 100.100.2.136 and 100.100.2.138.
    * `status` - The status of the resource.
