---
subcategory: "CBWP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cbwp_commonbandwidthpackageattachment"
sidebar_current: "docs-Alibabacloudstack-cbwp-commonbandwidthpackageattachment"
description: |- 
  Provides a cbwp Commonbandwidthpackageattachment resource.
---

# alibabacloudstack_cbwp_commonbandwidthpackageattachment
-> **NOTE:** Alias name has: `alibabacloudstack_common_bandwidth_package_attachment`

Provides a cbwp Commonbandwidthpackageattachment resource.

## Example Usage

### Basic Usage

```hcl
variable "name" {
    default = "tf-testAccBandwidtchPackage3659"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
    bandwidth   = "2"
    name        = "${var.name}"
    description = "${var.name}_description"
}

resource "alibabacloudstack_eip" "default" {
    name        = "${var.name}"
    bandwidth   = "2"
}

resource "alibabacloudstack_common_bandwidth_package_attachment" "default" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.default.id}"
    instance_id          = "${alibabacloudstack_eip.default.id}"
}
```

### Advanced Usage with Multiple EIPs

```hcl
variable "name" {
    default = "tf-testAccBandwidtchPackageAdvanced"
}

resource "alibabacloudstack_common_bandwidth_package" "advanced" {
    bandwidth   = "10"
    name        = "${var.name}"
    description = "${var.name}_description"
}

resource "alibabacloudstack_eip" "eip1" {
    name        = "${var.name}-eip1"
    bandwidth   = "2"
}

resource "alibabacloudstack_eip" "eip2" {
    name        = "${var.name}-eip2"
    bandwidth   = "3"
}

resource "alibabacloudstack_common_bandwidth_package_attachment" "eip1_attachment" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.advanced.id}"
    instance_id          = "${alibabacloudstack_eip.eip1.id}"
}

resource "alibabacloudstack_common_bandwidth_package_attachment" "eip2_attachment" {
    bandwidth_package_id = "${alibabacloudstack_common_bandwidth_package.advanced.id}"
    instance_id          = "${alibabacloudstack_eip.eip2.id}"
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, ForceNew) The ID of the Internet Shared Bandwidth instance. This field cannot be modified after creation.
* `instance_id` - (Required, ForceNew) The ID of the Elastic IP Address (EIP) that you want to associate with the shared bandwidth package. This field cannot be modified after creation. You can specify up to 50 EIP IDs for association. Separate multiple IDs with commas (,) if needed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the attachment. It is formatted as `<bandwidth_package_id>:<instance_id>`.
* `status` - The status of the attachment. Possible values include `Attached` and `Detached`.
* `creation_time` - The time when the attachment was created. This is useful for tracking the lifecycle of the attachment.