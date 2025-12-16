---
subcategory: "Hsm"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hsms"
sidebar_current: "docs-alibabacloudstack-datasource-hsms"
description: |-
  Provides a list of Hsm instances available to Alicloudstck.
---

# alibabacloudstack\_hsms

This data source provides a list of HSM instances in an Alibaba Cloud Stack zone.

## Example Usage

```hcl
data "alibabacloudstack_hsm_vendors" "default" {
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_hsms" "default" {
  vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].code}"
  product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].products[0].code}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
output "first_hsm_id" {
  value = data.alibabacloudstack_hsms.this.hsms.0.id
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required) The Zone of the HSM instance.
* `vendor_code` - (Required) The code of the HSM vendor.
* `product_code` - (Required) The code of the HSM product.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of HSM instance IDs.
* `hsms` - A list of HSM instances. Each element contains the following attributes:
  * `id` - The ID of the HSM instance.
