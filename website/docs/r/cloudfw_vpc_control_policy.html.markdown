---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_vpc_control_policy"
description: |-
  Manages Cloud Firewall VPC control policies.
---

# alibabacloudstack_cloudfw_vpc_control_policy

Manages Cloud Firewall VPC control policies, used to define access control rules for VPC network traffic.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testacc_vpc_control_policy92894"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
  group_type   = "ip"
  group_name   = var.name
  address_list = ["100.100.100.100/30"]
  description  = "test address book"
}

resource "alibabacloudstack_cloudfw_vpc_control_policy" "default" {
  dest_port        = "33/33"
  acl_action       = "log"
  release          = true
  destination      = "0.0.0.0/16"
  proto            = "UDP"
  application_id   = "0"
  application_name = "ANY"
  source_type      = "net"
  destination_type = "net"
  source           = "0.0.0.0/16"
  description      = var.name
  dest_port_type   = "port"
}
```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Required) The action of the access control policy. Valid values: `accept` (allow), `drop` (deny), `log` (log).
* `application_id` - (Required) The application ID. For example, `0` indicates ANY (all applications).
* `application_name` - (Required) The application name. For example, `ANY` indicates all applications.
* `description` - (Required) The description of the policy, used to identify the policy purpose.
* `destination` - (Required) The destination address. Supports CIDR format (e.g., `0.0.0.0/0`) or address book name.
* `destination_type` - (Required) The destination address type. Valid values: `net` (network segment), `group` (address book).
* `proto` - (Required) The protocol type. Valid values: `TCP`, `UDP`, `ANY`.
* `source` - (Required) The source address. Supports CIDR format (e.g., `192.168.1.0/24`) or address book name.
* `source_type` - (Required) The source address type. Valid values: `net` (network segment), `group` (address book).
* `dest_port` - (Optional) The destination port range. Format: `start_port/end_port` (e.g., `80/80`).
* `dest_port_type` - (Optional) The destination port type. Valid values: `port` (port), `group` (port group).
* `new_order` - (Optional) The new policy order. Default value: `-1` (indicating adding to the end of the policy list).
* `release` - (Optional) Whether to publish the policy. Valid values: `true` (publish), `false` (do not publish).
* `vpc_firewall_id` - (Optional) The VPC firewall instance ID. If not specified, the default firewall instance is used.
* `direction` - (Optional, Computed) The policy direction. Valid values: `inout` (bidirectional traffic), `in` (inbound traffic), `out` (outbound traffic).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The policy ID, formatted as `<acl_uuid>:<direction>`.
* `acl_uuid` - The unique identifier of the policy (AclUuid).
* `dest_port_group` - The destination port group ID (returned by the API).
* `dest_port_group_ports` - The destination port group port list (returned when `dest_port_type` is `group`).
* `destination_group_cidrs` - The destination address group CIDR list (returned when `destination_type` is `group`).
* `hit_times` - The number of times the policy has been matched (the number of matching traffic since creation).
* `order` - The current order value of the policy (used for policy priority sorting).
* `source_group_cidrs` - The source address group CIDR list (returned when `source_type` is `group`).

## Import

Cloud Firewall VPC Control Policy can be imported using the ID (format: `<acl_uuid>:<direction>`), e.g.

```
$ terraform import alibabacloudstack_cloudfw_vpc_control_policy.example acl-12345678:inout
```