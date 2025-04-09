---
subcategory: "CloudFirewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfirewall_controlpolicyorder"
sidebar_current: "docs-Alibabacloudstack-cloudfirewall-controlpolicyorder"
description: |- 
  Provides a cloudfirewall Controlpolicyorder resource.
---

# alibabacloudstack_cloudfirewall_controlpolicyorder
-> **NOTE:** Alias name has: `alibabacloudstack_cloud_firewall_control_policy_order`

Provides a cloudfirewall Controlpolicyorder resource.

For information about Cloud Firewall Control Policy Order and how to use it, see [What is Control Policy Order](https://www.alibabacloud.com/help/doc-detail/138867.htm).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alibabacloudstack_cloud_firewall_control_policy" "default" {
  direction        = "in"
  application_name = "ANY"
  description      = var.name
  acl_action       = "accept"
  source           = "127.0.0.1/32"
  source_type      = "net"
  destination      = "127.0.0.2/32"
  destination_type = "net"
  proto            = "ANY"
}

resource "alibabacloudstack_cloudfirewall_controlpolicyorder" "default" {
  acl_uuid  = alibabacloudstack_cloud_firewall_control_policy.default.acl_uuid
  direction = alibabacloudstack_cloud_firewall_control_policy.default.direction
  order     = 1
}
```

## Argument Reference

The following arguments are supported:

* `acl_uuid` - (Required, ForceNew) The unique identifier of the security access control policy. This is the ID assigned to the policy when it is created.
* `direction` - (Required, ForceNew) The direction of the traffic flow for which the access control policy applies. Valid values: `in`, `out`. Specifies whether the policy controls inbound or outbound traffic.
* `order` - (Optional, Int) <!--  AI CREATE  --> The priority level of the security access control policy. The priority value starts from 1, with smaller numbers indicating higher priority. A value of `-1` indicates the lowest priority. **NOTE:** From version 1.227.1, this field must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource in Terraform. It is formatted as `<acl_uuid>:<direction>`.