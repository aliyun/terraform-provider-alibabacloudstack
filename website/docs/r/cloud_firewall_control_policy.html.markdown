---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policy"
sidebar_current: "docs-Alibabacloudstack-resource-cloud-firewall-control-policy"
description: |-
  Manages Cloud Firewall control policy.
---

# alibabacloudstack_cloud_firewall_control_policy

Manages Cloud Firewall control policy, used to define access control rules for traffic that passes through Cloud Firewall.

-> **Note:** This resource can also be referred to by the following aliases:
-> - `alibabacloudstack_cloudfw_controlpolicy`

## Example Usage

### Basic Usage

```hcl
resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
  acl_action       = "accept"
  application_name = "HTTP"
  description      = "tf-testacc-control-policy"
  destination      = "192.168.0.0/24"
  destination_type = "net"
  direction        = "out"
  proto            = "TCP"
  source           = "10.0.0.0/8"
  source_type      = "net"
}
```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Required) The action to perform on traffic that matches the access control policy. Valid values: `accept` (allow), `drop` (deny), `log` (observe).
* `application_name` - (Required) The application type that the access control policy supports. Valid values: `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
* `description` - (Required) The description of the access control policy.
* `destination` - (Required) The destination address in the access control policy. The value depends on `destination_type`:
  - When `destination_type` is `net`, specify a CIDR block (e.g., `192.168.0.0/24`).
  - When `destination_type` is `group`, specify an address book name (e.g., `db_group`).
  - When `destination_type` is `domain`, specify a domain name (e.g., `*.example.com`).
  - When `destination_type` is `location`, specify a region code (e.g., `BJ11`, `ZB`).
* `destination_type` - (Required) The type of the destination address. Valid values: `group` (address book), `location` (region), `net` (CIDR block), `domain` (domain name).
* `direction` - (Required) The direction of the traffic to which the access control policy applies. Valid values: `in` (inbound), `out` (outbound). Modifying this parameter will force a new resource.
* `proto` - (Required) The protocol type of the traffic in the access control policy. Valid values: `ANY`, `TCP`, `UDP`, `ICMP`.
* `source` - (Required) The source address in the access control policy. The value depends on `source_type`:
  - When `source_type` is `net`, specify a CIDR block (e.g., `192.168.0.0/24`).
  - When `source_type` is `group`, specify an address book name (e.g., `db_group`).
  - When `source_type` is `location`, specify a region code (e.g., `BJ11`, `ZB`).
* `source_type` - (Required) The type of the source address. Valid values: `group` (address book), `location` (region), `net` (CIDR block).
* `dest_port` - (Optional) The destination port in the access control policy. Required when `dest_port_type` is `port`.
* `dest_port_group` - (Optional) The destination port address book name in the access control policy. Required when `dest_port_type` is `group`.
* `dest_port_type` - (Optional) The type of the destination port. Valid values: `group` (port group), `port` (port).
* `lang` - (Optional) The language of the request. Valid values: `en` (English), `zh` (Chinese).
* `release` - (Optional) The status of the access control policy. By default, the policy is enabled after it is created. Valid values: `true` (enabled), `false` (disabled).
* `source_ip` - (Optional) The source IP address of the request.

-> **Note:** The parameters `dest_port` and `dest_port_group` are mutually exclusive based on the value of `dest_port_type`. When `dest_port_type` is `port`, set `dest_port`. When `dest_port_type` is `group`, set `dest_port_group`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the resource. The format is `<acl_uuid>:<direction>`.
* `acl_uuid` - The unique identifier of the access control policy.

## Import

Cloud Firewall Control Policy can be imported using the `acl_uuid` and `direction`, e.g.

```
$ terraform import alibabacloudstack_cloud_firewall_control_policy.example acl-12345678:out
```
