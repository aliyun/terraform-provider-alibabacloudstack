---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfirewall_controlpolicy"
sidebar_current: "docs-Alibabacloudstack-cloudfirewall-controlpolicy"
description: |- 
  Provides a cloudfirewall Controlpolicy resource.
---

# alibabacloudstack_cloudfirewall_controlpolicy
-> **NOTE:** Alias name has: `alibabacloudstack_cloud_firewall_control_policy`

Provides a cloudfirewall Controlpolicy resource.

## Example Usage

### Basic Usage

```terraform
variable "name" {
    default = "tf-testacccloud_firewallcontrol_policy46819"
}

resource "alibabacloudstack_cloudfirewall_controlpolicy" "default" {
  source           = "0.0.0.0/0"
  proto            = "ANY"
  destination      = "0.0.0.0/0"
  application_name = "ANY"
  acl_action       = "accept"
  dest_port_type   = "port"
  release          = "true"
  description      = "test"
  direction        = "in"
  source_type      = "net"
  dest_port        = "80"
  destination_type = "net"
}
```

### Advanced Usage with `acl_uuid` and `ip_version`

```terraform
resource "alibabacloudstack_cloudfirewall_controlpolicy" "example" {
  acl_uuid         = "example-acl-uuid"
  ip_version       = "ipv4"
  source           = "192.168.1.0/24"
  proto            = "TCP"
  destination      = "10.0.0.0/24"
  application_name = "HTTP"
  acl_action       = "drop"
  dest_port_type   = "port"
  release          = "false"
  description      = "Advanced example"
  direction        = "out"
  source_type      = "net"
  dest_port        = "8080"
  destination_type = "net"
}
```

## Argument Reference

The following arguments are supported:

* `acl_action` - (Required) The action that Cloud Firewall performs on the traffic. Valid values: `accept`, `drop`, `log`.
* `acl_uuid` - (ForceNew, Optional) The unique ID of the access control policy. If not specified, Terraform will automatically generate one.
* `application_name` - (Required) The application type that the access control policy supports. If `direction` is `in`, the valid value is `ANY`. If `direction` is `out`, the valid values are `ANY`, `HTTP`, `HTTPS`, `MQTT`, `Memcache`, `MongoDB`, `MySQL`, `RDP`, `Redis`, `SMTP`, `SMTPS`, `SSH`, `SSL`, `VNC`.
* `description` - (Required) The description of the access control policy.
* `dest_port` - (Optional) The destination port defined in the access control policy. Required if `dest_port_type` is set to `port`.
* `dest_port_group` - (Optional) The destination port address book defined in the access control policy. Required if `dest_port_type` is set to `group`.
* `dest_port_type` - (Optional) The destination port type defined in the access control policy. Valid values: `group`, `port`.
* `destination` - (Required) The destination address defined in the access control policy.
* `destination_type` - (Required) The type of the destination address. Valid values:
  * If `direction` is `in`: `net`, `group`
  * If `direction` is `out`: `net`, `group`, `domain`, `location`
* `direction` - (Required, ForceNew) The direction of the traffic. Valid values: `in`, `out`.
* `ip_version` - (Optional) The IP version. Valid values: `ipv4`, `ipv6`.
* `lang` - (Optional) The language for the description. Valid values: `en`, `zh`.
* `proto` - (Required) The protocol used in the access control policy. Valid values: `TCP`, `UDP`, `ANY`, `ICMP`.
* `release` - (Optional) Specifies whether the access control policy is enabled. By default, an access control policy is enabled after it is created. Valid values: `true`, `false`.
* `source` - (Required) The source address defined in the access control policy.
* `source_ip` - (Optional) The source IP address.
* `source_type` - (Required) The type of the source address. Valid values:
  * If `direction` is `in`: `net`, `group`, `location`
  * If `direction` is `out`: `net`, `group`

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the Cloud Firewall Control Policy. The format is `<acl_uuid>:<direction>`.
* `acl_uuid` - The unique ID of the access control policy.
* `dest_port` - The destination port defined in the access control policy.
* `dest_port_group` - The destination port address book defined in the access control policy.
* `dest_port_type` - The destination port type defined in the access control policy.
* `release` - Specifies whether the access control policy is enabled.
* `source_ip` - The source IP address.

## Import

Cloud Firewall Control Policy can be imported using the `id`, e.g.

```
$ terraform import alibabacloudstack_cloudfirewall_controlpolicy.example <acl_uuid>:<direction>
```