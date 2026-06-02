---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_instances"
description: |-
  Query Alibaba Cloud Hardware Security Module (HSM) instances.
---

# alibabacloudstack_cspprivate_hsm_instances

Query Alibaba Cloud Hardware Security Module (HSM) instances.

> **NOTE:** This data source is used to query HSM instances and does not support create, modify, or delete operations.

## Example Usage

```hcl
variable "name" {
  default = "test-tf-cspprivate-hsm"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
  vendor_code  = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
  vsm_type     = "gvsm"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  alias_name   = "${var.name}"
}

data "alibabacloudstack_cspprivate_hsm_instances" "default" {
  name_regex = "${alibabacloudstack_cspprivate_hsm_instance.default.alias_name}"
  ids        = ["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs to filter results by.
* `name_regex` - (Optional) A regex string to filter results by instance name (Remark).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source.
* `ids` - A list of instance IDs.
* `names` - A list of instance names (Remark) corresponding to the returned instances.
* `instances` - A list of HSM instances matching the filter criteria. Each element contains the following attributes:
  * `id` - The ID of the resource.
  * `instance_id` - The ID of the HSM instance.
  * `hsm_status` - The status of the HSM instance. 1: uninitialized, 3: released, 4: failed, 5: running, 6: syncing, 7: resetting, 8: disabled.
  * `vpc_id` - The VPC ID assigned to the HSM instance.
  * `vswitch_id` - The vSwitch ID assigned to the HSM instance.
  * `ip` - The classic network IP address associated with the HSM instance.
  * `alias_name` - The alias or alias_name of the HSM instance.
  * `product_code` - The product model code of the device used by the HSM instance.
  * `vendor_code` - The vendor code of the device used by the HSM instance.
  * `vsm_type` - The type of the HSM instance. evsm: financial data HSM, gvsm: general server HSM, svsm: signature verification server HSM.
  * `zone_id` - The zone ID where the HSM instance is located.