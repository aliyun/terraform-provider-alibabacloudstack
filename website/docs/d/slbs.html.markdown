---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slbs"
sidebar_current: "docs-alibabacloudstack-datasource-slbs"
description: |-
    Provides a list of server load balancers to the user.
---

# alibabacloudstack\_slbs

This data source provides the server load balancers of the current Alibabacloudstack Cloud user.

## Example Usage

```
data "alibabacloudstack_slbs" "slbs_ds" {
  name_regex = "sample_slb"
}

output "first_slb_id" {
  value = "${data.alibabacloudstack_slbs.slbs_ds.slbs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLBs IDs.
* `name_regex` - (Optional) A regex string to filter results by SLB name.
* `master_availability_zone` - (Optional) Master availability zone of the SLBs.
* `slave_availability_zone` - (Optional) Slave availability zone of the SLBs.
* `network_type` - (Optional) Network type of the SLBs. Valid values: `vpc` and `classic`.
* `vpc_id` - (Optional) ID of the VPC linked to the SLBs.
* `vswitch_id` - (Optional) ID of the VSwitch linked to the SLBs.
* `address` - (Optional) Service address of the SLBs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of slb IDs.
* `names` - A list of slb names.
* `slbs` - A list of SLBs. Each element contains the following attributes:
  * `id` - ID of the SLB.
  * `region_id` - Region ID the SLB belongs to.
  * `master_availability_zone` - Master availability zone of the SLBs.
  * `slave_availability_zone` - Slave availability zone of the SLBs.
  * `name` - SLB name.
  * `network_type` - Network type of the SLB. Possible values: `vpc` and `classic`.
  * `vpc_id` - ID of the VPC the SLB belongs to.
  * `vswitch_id` - ID of the VSwitch the SLB belongs to.
  * `address` - Service address of the SLB.
  * `creation_time` - SLB creation time.