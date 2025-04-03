---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_accesscontrollist"
sidebar_current: "docs-Alibabacloudstack-slb-accesscontrollist"
description: |- 
  Provides a slb Accesscontrollist resource.
---

# alibabacloudstack_slb_accesscontrollist
-> **NOTE:** Alias name has: `alibabacloudstack_slb_acl`

Provides a slb Accesscontrollist resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccslbaccess_control_list69490"
}

resource "alibabacloudstack_slb_accesscontrollist" "default" {
  acl_name            = "Rdk_test_name01"
  address_ip_version  = "ipv4"

  entry_list {
    entry   = "10.10.10.0/24"
    comment = "first"
  }

  entry_list {
    entry   = "168.10.10.0/24"
    comment = "second"
  }

  tags = {
    CreatedBy = "Terraform"
    Purpose   = "Testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_name` - (Required) The name of the access control list.
* `address_ip_version` - (Optional, ForceNew) The IP version. Valid values: `ipv4` and `ipv6`. Our plugin provides a default value of `ipv4`.
* `entry_list` - (Optional) A list of entries (IP addresses or CIDR blocks) to be added. At most 50 entries can be supported in one resource. Each entry contains:
  * `entry` - (Required) An IP address or CIDR block.
  * `comment` - (Optional) The comment for the entry.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the access control list.
* `acl_name` - The name of the access control list.
* `address_ip_version` - The IP version. Valid values: `ipv4` and `ipv6`.