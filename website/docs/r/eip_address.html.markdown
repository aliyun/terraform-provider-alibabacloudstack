---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_address"
sidebar_current: "docs-Alibabacloudstack-eip-address"
description: |- 
  Provides a eip Address resource.
---

# alibabacloudstack_eip_address
-> **NOTE:** Alias name has: `alibabacloudstack_eip`

Provides a EIP Address resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAcceEipName5478"
}

resource "alibabacloudstack_eip_address" "default" {
  name        = var.name
  description = "This is a test EIP address."
  bandwidth   = "5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the EIP instance. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with `http://` or `https://`.
* `description` - (Optional) Description of the EIP instance. This description can have a string of 2 to 256 characters. It cannot begin with `http://` or `https://`. Default value is null.
* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). If this value is not specified, it defaults to **5** Mbps.  
  - When `payment_type` is set to `PayAsYouGo` and `internet_charge_type` is set to `PayByBandwidth`, valid values for `bandwidth` are **1** to **500**.
  - When `payment_type` is set to `PayAsYouGo` and `internet_charge_type` is set to `PayByTraffic`, valid values for `bandwidth` are **1** to **200**.
  - When `payment_type` is set to `Subscription`, valid values for `bandwidth` are **1** to **1000**.
* `ip_address` - (Optional, ForceNew) The elastic IP address. It must be a valid IP address. Supports a maximum of 50 EIPs.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the EIP.
* `status` - The state of the EIP. Valid values:
  - **Associating**
  - **Unassociating**
  - **InUse**
  - **Available**
  - **Releasing**
* `ip_address` - The elastic IP address assigned to the EIP.