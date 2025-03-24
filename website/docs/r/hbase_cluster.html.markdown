---
subcategory: "HBase"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hbase_cluster"
sidebar_current: "docs-Alibabacloudstack-hbase-cluster"
description: |- 
  Provides a hbase Cluster resource.
---

# alibabacloudstack_hbase_cluster
-> **NOTE:** Alias name has: `alibabacloudstack_hbase_instance`

Provides a hbase Cluster resource.

## Example Usage

```hcl
variable "name" {
	default = "tf-testAccVpc1175381"
}
variable "password" {
}

data "alibabacloudstack_zones" "default" {}

data "alibabacloudstack_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alibabacloudstack_vpc" "default" {
	name       = var.name
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = "${alibabacloudstack_vpc.default.id}"
	cidr_block        = "172.16.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.ids.0
	name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
	count   = 2
	vpc_id  = alibabacloudstack_vpc.default.id
	name    = var.name
}

resource "alibabacloudstack_hbase_instance" "default" {
	name                  = var.name
	ip_white             = "192.168.1.2"
	vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
	account              = "adminu"
	zone_id              = "${data.alibabacloudstack_zones.default.zones.0.id}"
	security_groups      = ["${alibabacloudstack_security_group.default.0.id}", "${alibabacloudstack_security_group.default.1.id}"]
	core_disk_size       = "480"
	engine_version       = "2.0"
	cold_storage_size    = "900"
	deletion_protection  = "false"
	core_instance_type   = "hbase.sn1.large"
	master_instance_type = "hbase.sn1.large"
	immediate_delete_flag = "true"
	maintain_start_time  = "14:00Z"
	password             = var.password
	maintain_end_time    = "16:00Z"
	core_disk_type       = "cloud_efficiency"
	tags = {
		Created = "TF-update"
		For     = "acceptance test 123"
	}
}
```

This example demonstrates how to create an HBase instance with various configurations such as network settings, security groups, and storage options.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HBase cluster. It must be 2-128 characters long and can include Chinese characters, English letters, numbers, periods (`.`), underscores (`_`), or dashes (`-`).
* `zone_id` - (Optional, ForceNew) The Zone ID where the HBase instance will be launched. If `vswitch_id` is specified, this field can either be empty or consistent with the VSwitch's zone.
* `engine` - (Optional, ForceNew) The engine type of the cluster. Valid values are `hbase`, `hbaseue`, or `bds`. Starting from version 1.73.0, `hbaseue` and `bds` are supported.
* `engine_version` - (Required, ForceNew) The major version of HBase. Valid values:
  - For `hbase`: `1.1` or `2.0`
  - For `hbaseue`: `2.0`
  - For `bds`: `1.0`
* `master_instance_type` - (Required, ForceNew) The specification of the master node. Refer to [Instance Specifications](https://help.aliyun.com/document_detail/53532.html) or use the `describeInstanceType` API.
* `core_instance_type` - (Required, ForceNew) The specification of the core node. Refer to [Instance Specifications](https://help.aliyun.com/document_detail/53532.html) or use the `describeInstanceType` API.
* `core_instance_quantity` - (Optional) The number of core nodes. Default is `2`, and the range is `[1-200]`.
* `core_disk_type` - (Optional, ForceNew) The type of disk for core nodes. Valid values:
  - `cloud_ssd`
  - `cloud_essd_pl1`
  - `cloud_efficiency`
  - `local_hdd_pro`
  - `local_ssd_pro`
  When `engine=bds`, no need to set disk type (or leave it as an empty string).
* `core_disk_size` - (Optional) The size of the storage for one core node in GB. Valid when `engine=hbase/hbaseue`. Not required for `bds`. Value ranges:
  - Custom storage space: `[20, 64000]`
  - Cluster: `[400, 64000]`, increments by 40GB.
  - Single: `[20-500]`, increments by 1GB.
* `pay_type` - (Optional) The payment type. Valid values are `PrePaid` or `PostPaid`. Defaults to `PostPaid`. You can convert `PostPaid` to `PrePaid` or vice versa starting from version 1.115.0.
* `duration` - (Optional, ForceNew) The subscription duration in months. Valid values: `1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36`. Only valid when `pay_type=PrePaid`. Values `12, 24, 36` represent 1, 2, and 3 years respectively.
* `auto_renew` - (Optional, ForceNew) Whether auto-renewal is enabled. Valid values are `true` or `false`. Defaults to `false`. Only valid when `pay_type=PrePaid`.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. If specified, the network type is `vpc`. If not specified, the network type is `classic`. International sites do not support classic networks.
* `cold_storage_size` - (Optional, ForceNew) The size of cold storage in GB. Valid values: `0` or `[800, 1000000]`, increments by 10GB. `0` means cold storage is disabled.
* `maintain_start_time` - (Optional) The start time of the maintenance period in UTC format `HH:mmZ`. Example: `02:00Z`.
* `maintain_end_time` - (Optional) The end time of the maintenance period in UTC format `HH:mmZ`. Example: `04:00Z`.
* `deletion_protection` - (Optional) Whether deletion protection is enabled. Valid values are `true` or `false`. Defaults to `false`.
* `immediate_delete_flag` - (Optional) Whether immediate deletion is enabled. Valid values are `true` or `false`. Defaults to `false`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `account` - (Optional) The account for the cluster web UI. Length must be between `0-128` characters.
* `password` - (Optional) The password for the cluster web UI account. Length must be between `0-128` characters.
* `ip_white` - (Optional) The whitelist of IP addresses for the cluster.
* `security_groups` - (Optional) The list of security group IDs associated with the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the HBase cluster.
* `zone_id` - The Zone ID of the cluster.
* `master_instance_quantity` - The number of master nodes in the cluster.
* `ui_proxy_conn_addrs` - The Web UI proxy connection information list.
  * `net_type` - The access type of the connection address. Return value:
    - `2`: Intranet access.
    - `0`: Public network access.
  * `conn_addr` - The connection address.
  * `conn_addr_port` - The connection port.
* `zk_conn_addrs` - The Zookeeper connection information list.
  * `net_type` - The access type of the connection address. Return value:
    - `2`: Intranet access.
    - `0`: Public network access.
  * `conn_addr` - The connection address.
  * `conn_addr_port` - The connection port.
* `slb_conn_addrs` - The SLB connection information list.
  * `net_type` - The access type of the connection address. Return value:
    - `2`: Intranet access.
    - `0`: Public network access.
  * `conn_addr` - The connection address.
  * `conn_addr_port` - The connection port. 

## Import

HBase clusters can be imported using the `id`, e.g.

```bash
$ terraform import alicloud_hbase_instance.example hb-wz96815u13k659fvd
```