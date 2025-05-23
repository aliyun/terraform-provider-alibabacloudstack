---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-cloud-firewall-control-policies"
description: |- 
  Provides a list of cloud firewall control policies owned by an alibabacloudstack account.
---

# alibabacloudstack_cloud_firewall_control_policies

This data source provides a list of cloud firewall control policies in an alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage:

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

* `acl_action` - (Optional, ForceNew) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `acl_uuid` - (Optional, ForceNew) The unique ID of the access control policy.
* `description` - (Optional, ForceNew) The description of the access control policy.
* `destination` - (Optional, ForceNew) The destination address defined in the access control policy.
* `direction` - (Required, ForceNew) The direction of the traffic to which the access control policy applies. Valid values: `in`, `out`.
* `ip_version` - (Optional, ForceNew) The IP version of the address in the access control policy.
* `lang` - (Optional, ForceNew) The language of the content within the response. Valid values: `en`, `zh`.
* `proto` - (Optional, ForceNew) The type of the protocol in the access control policy. Valid values: If `direction` is `in`, the valid value is `ANY`. If `direction` is `out`, the valid values are `ANY`, `TCP`, `UDP`, `ICMP`.
* `source` - (Optional, ForceNew) The source address in the access control policy.
* `source_ip` - (Removed since v1.213.0) The source IP address of the request. **NOTE:** Field `source_ip` has been removed from provider version 1.213.0.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Control Policy IDs.
* `policies` - A list of Cloud Firewall Control Policies. Each element contains the following attributes:
  * `id` - The ID of the Control Policy. It formats as `<acl_uuid>:<direction>`.
  * `acl_uuid` - The unique ID of the access control policy.
  * `acl_action` - The action that Cloud Firewall performs on the traffic.
  * `application_id` - The application ID in the access control policy.
  * `application_name` - The type of the application that the access control policy supports.
  * `description` - The description of the access control policy.
  * `dest_port` - The destination port in the access control policy.
  * `dest_port_group` - The name of the destination port address book in the access control policy.
  * `dest_port_group_ports` - The ports in the destination port address book.
  * `dest_port_type` - The type of the destination port in the access control policy.
  * `destination` - The destination address in the access control policy.
  * `destination_group_cidrs` - The CIDR blocks in the destination address book.
  * `destination_group_type` - The type of the destination address book in the access control policy.
  * `destination_type` - The type of the destination address in the access control policy.
  * `dns_result` - The DNS resolution result.
  * `dns_result_time` - The timestamp of the DNS resolution result.
  * `hit_times` - The number of hits for the access control policy.
  * `order` - The priority of the access control policy.
  * `proto` - The type of the protocol in the access control policy.
  * `release` - Indicates whether the access control policy is enabled.
  * `source` - The source address in the access control policy.
  * `source_group_cidrs` - The CIDR blocks in the source address book.
  * `source_group_type` - The type of the source address book in the access control policy.
  * `source_type` - The type of the source address in the access control policy.