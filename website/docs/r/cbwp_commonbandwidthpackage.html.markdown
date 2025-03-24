---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackage"
sidebar_current: "docs-Alibabacloudstack-cbwp-commonbandwidthpackage"
description: |- 
  Provides a cbwp Commonbandwidthpackage resource.
---

# alibabacloudstack_cbwp_commonbandwidthpackage
-> **NOTE:** Alias name has: `alibabacloudstack_common_bandwidth_package`

Provides a cbwp Commonbandwidthpackage resource.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf-testAccCommonBandwidthPackage481904"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
  internet_charge_type = "PayByTraffic"
  bandwidth            = "10"
  name                = var.name
  description         = "Test common bandwidth package"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The maximum bandwidth of the Internet Shared Bandwidth instance. Unit: Mbit/s. Valid values: **1** to **1000**. Default value: **1**.
* `internet_charge_type` - (Optional, ForceNew) The billing method of the Internet Shared Bandwidth instance. Set the value to **PayByTraffic**, which specifies the pay-by-data-transfer billing method.
* `ratio` - (Optional, ForceNew) The percentage of the minimum bandwidth commitment. Set the parameter to **20**. > This parameter is available only on the Alibaba Cloud China site.
* `name` - (Optional) The name of the Internet Shared Bandwidth instance. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (`_`), and hyphens (`-`). The name must start with a letter.
* `description` - (Optional) The description of the Internet Shared Bandwidth instance. The description must be 2 to 256 characters in length and start with a letter. The description cannot start with `http://` or `https://`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the common bandwidth package instance.
* `bandwidth_package_name` - The name of the Internet Shared Bandwidth instance. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (`_`), and hyphens (`-`). The name must start with a letter.