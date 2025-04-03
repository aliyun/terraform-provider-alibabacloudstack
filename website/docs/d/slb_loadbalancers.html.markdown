---
subcategory: "SLB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_loadbalancers"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-loadbalancers"
description: |- 
  Provides a list of slb loadbalancers owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_loadbalancers
-> **NOTE:** Alias name has: `alibabacloudstack_slbs`

This data source provides a list of SLB load balancers in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_loadbalancers" "slbs_ds" {
  name_regex = "sample_slb"
}

output "first_slb_id" {
  value = "${data.alibabacloudstack_slb_loadbalancers.slbs_ds.slbs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of SLB Load Balancer IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by SLB name.
* `master_availability_zone` - (Optional, ForceNew) Master availability zone of the SLBs.
* `slave_availability_zone` - (Optional, ForceNew) Slave availability zone of the SLBs.
* `network_type` - (Optional, ForceNew) Network type of the SLBs. Valid values: `vpc` and `classic`.
* `vpc_id` - (Optional, ForceNew) ID of the VPC linked to the SLBs.
* `vswitch_id` - (Optional, ForceNew) ID of the VSwitch linked to the SLBs.
* `address` - (Optional, ForceNew) Service address of the SLBs.
* `tags` - (Optional, ForceNew) A map of tags assigned to the SLB instances. The tags can have a maximum of 5 key-value pairs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of SLB names.
* `slbs` - A list of SLBs. Each element contains the following attributes:
  * `id` - ID of the SLB.
  * `region_id` - Region ID the SLB belongs to.
  * `master_availability_zone` - Master availability zone of the SLB.
  * `slave_availability_zone` - Slave availability zone of the SLB.
  * `name` - Name of the SLB.
  * `network_type` - Network type of the SLB. Possible values: `vpc` and `classic`.
  * `vpc_id` - ID of the VPC the SLB belongs to.
  * `vswitch_id` - ID of the VSwitch the SLB belongs to.
  * `address` - Service address of the SLB.
  * `creation_time` - Creation time of the SLB.
  * `tags` - Tags assigned to the SLB.