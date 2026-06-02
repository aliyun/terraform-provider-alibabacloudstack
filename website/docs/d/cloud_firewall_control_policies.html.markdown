---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policies"
description: |- 
  Provides a list of cloud firewall control policies owned by an alibabacloudstack account.
---

# Data Source: alibabacloudstack_cloud_firewall_control_policies

This data source provides a list of cloud firewall control policies in an alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_cloud_firewall_control_policies" "example" {
  direction = "in"
  acl_action = "accept"
  source = "192.168.0.0/16"
  destination = "10.0.0.0/8"
  proto = "TCP"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required) The direction of the traffic to which the access control policy applies. Valid values: `in`, `out`.
* `acl_action` - (Optional) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `acl_uuid` - (Optional) The unique ID of the access control policy.
* `description` - (Optional) The description of the access control policy.
* `destination` - (Optional) The destination address defined in the access control policy.
* `proto` - (Optional) The type of the protocol in the access control policy. Valid values: `TCP`, `UDP`, `ANY`, `ICMP`.
* `source` - (Optional) The source address in the access control policy.
* `source_ip` - (Optional) The source IP address of the request.
* `output_file` - (Optional, Deprecated) The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Control Policy IDs. Each ID formats as `<acl_uuid>:<direction>`.
* `policies` - A list of Cloud Firewall Control Policies. Each element contains the following attributes:
  * `id` - The ID of the Control Policy. It formats as `<acl_uuid>:<direction>`.
  * `acl_uuid` - The unique ID of the access control policy.
  * `acl_action` - The action that Cloud Firewall performs on the traffic.
  * `application_id` - The application ID in the access control policy.
  * `application_name` - The application name in the access control policy.
  * `description` - The description of the access control policy.
  * `dest_port` - The destination port in the access control policy.
  * `dest_port_group` - The name of the destination port address book in the access control policy.
  * `dest_port_group_ports` - A list of ports in the destination port address book.
  * `dest_port_type` - The type of the destination port in the access control policy.
  * `destination` - The destination address in the access control policy.
  * `destination_group_cidrs` - A list of CIDR blocks in the destination address book.
  * `destination_group_type` - The type of the destination address book in the access control policy.
  * `destination_type` - The type of the destination address in the access control policy.
  * `direction` - The direction of the traffic to which the access control policy applies.
  * `hit_times` - The number of hits for the access control policy.
  * `order` - The priority of the access control policy.
  * `proto` - The type of the protocol in the access control policy.
  * `release` - Indicates whether the access control policy is enabled.
  * `source` - The source address in the access control policy.
  * `source_group_cidrs` - A list of CIDR blocks in the source address book.
  * `source_group_type` - The type of the source address book in the access control policy.
  * `source_type` - The type of the source address in the access control policy.
