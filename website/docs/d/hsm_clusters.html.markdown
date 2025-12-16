---
subcategory: "Hardware Security Module (HSM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_clusters"
sidebar_current: "docs-alibabacloudstack-datasource-hsm-clusters"
description: |-
    Provides information about Alibaba Cloud HSM Clusters.
---

# alibabacloudstack_hsm_clusters

This data source provides information about Alibaba Cloud HSM Clusters.

## Example Usage

```hcl
variable "name" {
  default = "tf_hsm_cluster6821719"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  provider     = alibabacloudstack-common
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
}

data "alibabacloudstack_hsm_vendors" "default" {
}

resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

resource "alibabacloudstack_hsm_instance" "default" {
  product_code = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
  vendor_code  = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
}

resource "alibabacloudstack_hsm_cluster" "default" {
  cluster_name       = var.name
  master_instance_id = alibabacloudstack_hsm_instance.default.id
  vpc_id             = alibabacloudstack_vpc.default.id
  vswitch_ids        = alibabacloudstack_vswitch.default.id
  zone_nos           = data.alibabacloudstack_zones.default.zones.0.id
  ip_white_list      = "123.12.13.1/16"
  password           = random_password.password.0.result
}

data "alibabacloudstack_hsm_clusters" "default" {
  name_regex = alibabacloudstack_hsm_cluster.default.cluster_name
}
```

## Argument Reference

The following arguments are supported for filtering query results:

* `vsm_type` (String, ForceNew) - Used to filter clusters of specific types. Valid values: `evsm` (Financial Data HSM), `gvsm` (General Server HSM), `svsm` (Signature Verification HSM).

* `ids` (List of Strings, Optional) - Used to filter clusters by specific IDs. If specified, only clusters with IDs in the list will be returned.

* `name_regex` (String, Optional) - Used to filter cluster names by regular expression. Only clusters with names matching the regex will be returned.

## Attributes Reference

The following attributes are exported:

* `id` (String) - The ID of the cluster, serving as the unique identifier for the data source.

* `abnormal_type` (String) - The abnormal type of the cluster. This parameter has a value only when the cluster status is abnormal; otherwise, it is empty. Valid values: `3` (HSM data inconsistency), `4` (All HSMs inaccessible), `5` (HSM count inconsistency).

* `cluster_id` (String) - The ID of the cluster, serving as the unique identifier.

* `cluster_master` (String) - The ID of the master HSM instance in the cluster.

* `cluster_name` (String) - The user-defined name of the cluster.

* `cluster_size` (String) - The size of the cluster, i.e., the number of HSM instances under the cluster.

* `cluster_zones` (List of Objects) - The dataset of zones and vswitches under the cluster. Each vswitch belongs to only one zone.
  * `vswitch_id` (String) - The ID of the vswitch.
  * `zone_no` (String) - The zone ID.

* `gmt_create` (Integer) - The creation time of the cluster (Unix timestamp in milliseconds).

* `hsm_cluster_items` (List of Objects) - The dataset of sub-items in the cluster, where each sub-item is an HSM instance.
  * `id` (Integer) - The primary key ID of the sub-item in the cluster.
  * `instance_id` (String) - The ID of the HSM instance.
  * `ip` (String) - The classic network IP address registered for the HSM instance.
  * `is_master` (Integer) - Indicates whether the current HSM is the master HSM in the cluster (1 for yes, 0 for no).
  * `status` (Integer) - The status of the HSM instance. Valid values: `1` (Not initialized), `3` (Released), `5` (Running), `6` (Synchronizing), `7` (Resetting), `8` (Disabled).
  * `vsm_type` (String) - The device type of the HSM instance.
  * `zone_no` (String) - The zone ID where the HSM instance is located.

* `ip_white_list` (String) - The whitelist IP addresses of the cluster, separated by commas.

* `master_ip` (String) - The classic network IP address registered for the master HSM instance in the cluster.

* `product_code` (String) - The model code of the device used by the HSM instance.

* `vendor_code` (String) - The vendor code of the device used by the HSM instance.

* `vpc_id` (String) - The VPC instance ID of the cluster. Each cluster has only one VPC instance.

* `vsm_type` (String) - The device type of the master HSM instance in the cluster, which represents the device type for all HSM instances in the cluster. Valid values: `evsm` (Financial Data HSM), `gvsm` (General Server HSM), `svsm` (Signature Verification HSM).