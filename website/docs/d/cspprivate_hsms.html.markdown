---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsms"
description: |-
  Provides a list of CSP Private HSM instances available to Alibaba Cloud Stack.
---

# alibabacloudstack\_cspprivate_hsms

This data source provides a list of CSP Private HSM (Hardware Security Module) instances in an Alibaba Cloud Stack zone.


## Example Usage

```hcl
data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsms" "default" {
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].code
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].products[0].code
  zone_id      = data.alibabacloudstack_zones.default.zones[0].id
}

output "first_hsm_id" {
  value = data.alibabacloudstack_cspprivate_hsms.default.hsms[0].id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The ID of the zone where the HSM instances are deployed.
* `vendor_code` - (Required) The code of the HSM vendor.
* `product_code` - (Required) The code of the HSM product.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of CSP Private HSM instance IDs.
* `hsms` - A list of CSP Private HSM instances. Each element contains the following attributes:
  * `id` - The ID of the CSP Private HSM instance.
