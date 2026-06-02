---
subcategory: "Web Application Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_waf_instances"
description: |-
  provides a list of WAF instances in Alibaba Cloud Stack based on the provided filters.
---

# alibabacloudstack_waf_instances
-> **NOTE:** Alias name has: `alibabacloudstack_waf_instances`

This data source provides a list of WAF instances in Alibaba Cloud Stack based on the provided filters.

## Example Usage
hcl
data "alibabacloudstack_waf_instances" "example" {
  ids = ["waf-instance-1", "waf-instance-2"]
  output_file = "output.json"
}

output "instances" {
  value = data.alibabacloudstack_waf_instances.example.ids
}
## Argument Reference
The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of WAF instance IDs to filter the results.
* `name` - (Computed) The name of the WAF instance.
* `instance_status` - (Computed) The current status of the WAF instance.
* `instance_make_status` - (Computed) The creation or provisioning status of the WAF instance.
* `output_file` - (Optional) File name to save the results to in JSON format.
* `Deprecated`: This field is deprecated and will be removed in version 3.19.0. Use the local_file provider instead.
* `vpc_vswitch` - (Computed) A list containing VPC and vSwitch configuration for the WAF instance with the following structure:
* `vswitch_name` - Name of the vSwitch.
* `vswitch` - ID of the vSwitch.
* `cidr_block` - CIDR block of the vSwitch.
* `available_zone` - Availability zone where the vSwitch resides.
* `vpc` - ID of the associated VPC.
* `vpc_name` - Name of the associated VPC.
* `detector_specs` - (Computed) Specifications of the detection engine.
* `detector_version` - (Computed) Version or level of the detection engine.
* `detector_nodenum` - (Computed) Number of detection engine nodes deployed in a single availability zone.
## Argument Reference
In addition to all the above arguments, the following attributes are exported:

* `ids` - A list of WAF instance IDs matching the criteria.
* `instances` - A list of WAF instances with their detailed attributes. Each entry contains:
* `name` - (Computed) The name of the WAF instance.
* `instance_status` - (Computed) The current status of the WAF instance.
* `instance_make_status` - (Computed) The creation or provisioning status of the WAF instance.
* `vpc_vswitch` - (Computed) A list containing VPC and vSwitch configuration for the WAF instance with the following structure:
* `detector_specs` - (Computed) Specifications of the detection engine.
* `detector_version` - (Computed) Version or level of the detection engine.
* `detector_nodenum` - (Computed) Number of detection engine nodes deployed in a single 