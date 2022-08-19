---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_eip"
sidebar_current: "docs-apsarastack-resource-eip"
description: |-
  Provides a ECS EIP resource.
---

# apsarastack\_eip

Provides an elastic IP resource.

## Example Usage

```
# Create a new EIP.
resource "apsarastack_eip" "example" {
  bandwidth            = "10"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the EIP instance. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://.
* `description` - (Optional) Description of the EIP instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). If this value is not specified, then automatically sets it to 5 Mbps.
* `isp` - (Optional, ForceNew) The line type of the Elastic IP instance. Default to `BGP`. Other type of the isp need to open a whitelist.

## Attributes Reference

The following attributes are exported:

* `id` - The EIP ID.
* `bandwidth` - The elastic public network bandwidth.
* `status` - The EIP current status.
* `ip_address` - The elastic ip address

