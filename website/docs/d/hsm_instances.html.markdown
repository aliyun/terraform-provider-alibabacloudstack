---
subcategory: "Hardware Security Module (HSM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-hsm-instances"
description: |-
  Query Alibaba Cloud Hardware Security Module (HSM) instances.
---

# alibabacloudstack_hsm_instances

Query Alibaba Cloud Hardware Security Module (HSM) instances.

> `NOTE` This data source is used to query HSM instances and does not support create, modify, or delete operations.

## Example Usage

```hcl
variable "name" {
  default = "test-tf-hsm-instance78983"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_hsm_instance" "default" {
  product_code = "jnta.SJJ1528"
  vendor_code  = "jnta"
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id   ="${alibabacloudstack_vpc_vswitch.default.id}"
  ip           = "172.16.1.100"
  white_list   = "192.168.1.0/24"
}

data "alibabacloudstack_hsm_instances" "default" {
  status     = alibabacloudstack_hsm_instance.default.status
  name_regex = "test-tf-hsm-instance"
}
```

## Argument Reference

The following arguments are supported for filtering query results:

* `cluster_id` (String, Optional) - The ID of the cluster to which the HSM instance belongs.

* `ids` (List, Optional) - A list of instance IDs used to filter query results.

* `instance_id` (String, Optional) - The ID of the HSM instance.

* `name_regex` (String, Optional) - A regular expression used to filter results by instance name (Remark).

* `status` (String, Optional) - The status of the HSM instance. Multiple statuses can be separated by commas. Valid values: 1 (Not initialized), 3 (Released), 4 (Production failed), 5 (Running), 6 (Synchronizing), 7 (Resetting), 8 (Disabled).

* `vpc_ip` (String, Optional) - The classic network IP address registered for the HSM instance.

* `vsm_type` (String, Optional) - The device type of the HSM instance. Valid values: evsm (Financial Data HSM), gvsm (General Server HSM), svsm (Signature Verification Server HSM).

* `zone_no` (String, Optional) - The zone ID where the HSM instance is located.

* `enable_details` (Boolean, Optional, Default: false) - Whether to retrieve detailed information for each instance.

## Attributes Reference

The following attributes are exported:

* `id` (String) - The resource ID.

* `cluster_id` (String) - The ID of the cluster to which the HSM instance belongs.

* `cluster_name` (String) - The name of the cluster to which the HSM instance belongs.

* `hsm_id` (String) - The physical HSM ID associated with the instance.

* `hsm_status` (Integer) - The status of the HSM instance. Valid values: 1 (Not initialized), 3 (Released), 4 (Production failed), 5 (Running), 6 (Synchronizing), 7 (Resetting), 8 (Disabled).

* `ip` (String) - The classic network IP address associated with the HSM instance.

* `is_master` (Integer) - Whether the current HSM is the master HSM in the cluster. Valid values: 0 (No), 1 (Yes).

* `product_code` (String) - The device model code used by the HSM instance.

* `product_name` (String) - The device model name used by the HSM instance.

* `release_protection` (Integer) - Whether release protection is enabled. Valid values: 0 (Disabled), 1 (Enabled).

* `remark` (String) - The alias or remark of the HSM instance.

* `show_create_cluster` (Boolean) - Whether creating a cluster for this HSM instance is allowed.

* `vpc_id` (String) - The VPC ID assigned to the HSM instance.

* `vendor_code` (String) - The device vendor code used by the HSM instance.

* `vendor_name` (String) - The device vendor name used by the HSM instance.

* `vsm_type` (String) - The device type of the HSM instance. Valid values: evsm (Financial Data HSM), gvsm (General Server HSM), svsm (Signature Verification Server HSM).

* `vswitch_id` (String) - The vSwitch ID assigned to the HSM instance.

* `white_list` (List) - The whitelist IP list of the HSM instance. If the instance belongs to a cluster, the whitelist will be consistent with the cluster.

* `zone_no` (String) - The zone ID where the HSM instance is located.