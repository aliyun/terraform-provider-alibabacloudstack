---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_dbcluster"
sidebar_current: "docs-Alibabacloudstack-adb-dbcluster"
description: |- 
  Provides a adb Dbcluster resource.
---

# alibabacloudstack_adb_dbcluster
-> **NOTE:** Alias name has: `alibabacloudstack_adb_cluster` `alibabacloudstack_adb_db_cluster`

Provides a adb Dbcluster resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccadbCluster73485"
}

data "alibabacloudstack_ascm_resource_groups" "default" {
  name_regex = ""
}

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "172.16.0.0/16"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "ADB"
}

data "alibabacloudstack_vswitches" "default" {
  vpc_id   = alibabacloudstack_vpc.default.id
  zone_id  = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_vswitch" "default" {
  name                 = "tf_testAccAdb_vpc"
  vpc_id               = alibabacloudstack_vpc.default.id
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  cidr_block          = "172.16.0.0/24"
}

resource "alibabacloudstack_adb_db_cluster" "default" {
  vswitch_id             = alibabacloudstack_vswitch.default.id
  db_cluster_category   = "Basic"
  description           = var.name
  db_node_storage       = "200"
  mode                  = "reserver"
  cpu_type              = "intel"
  db_cluster_version    = "3.0"
  db_node_count         = "2"
  db_node_class         = "C8"
  cluster_type          = "analyticdb"
  payment_type          = "PayAsYouGo"
  maintain_time         = "23:00Z-00:00Z"
  security_ips          = ["10.168.1.12", "10.168.1.11"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when `payment_type` is `Subscription`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`. Default to `1`.
* `compute_resource` - (Optional) The specifications of computing resources in elastic mode. The increase of resources can speed up queries. AnalyticDB for MySQL automatically scales computing resources. For more information, see [ComputeResource](https://www.alibabacloud.com/help/en/doc-detail/144851.htm).
* `db_cluster_category` - (Required) The DB cluster category. Valid values: `Basic`, `Cluster`, `MixedStorage`.
* `db_cluster_class` - (Deprecated) It duplicates with attribute `db_node_class` and is deprecated from version 1.121.2.
* `storage_resource` - (Optional) The specifications of storage resources in elastic mode. The resources are used for data read and write operations. The increase of resources can improve the read and write performance of your cluster. For more information, see [Specifications](https://www.alibabacloud.com/help/en/doc-detail/144851.htm).
* `storage_type` - (Optional) Reserved parameters, not involved.
* `db_cluster_version` - (Optional, ForceNew) The DB cluster version. Value options: `3.0`. Default to `3.0`.
* `cluster_type` - (Required) The cluster type of the resource. Valid values: `analyticdb`, `AnalyticdbOnPanguHybrid`.
* `cpu_type` - (Required) The CPU type of the resource. Valid values: `intel`.
* `db_node_class` - (Optional) The DB node class. For more information, see [DBClusterClass](https://help.aliyun.com/document_detail/190519.html).
* `db_node_count` - (Optional) The number of DB nodes.
* `executor_count` - (Optional) The number of nodes. The node resources are used for data computing in elastic mode.
* `db_node_storage` - (Optional) The storage capacity of each DB node.
* `description` - (Optional) The description of the DBCluster.
* `maintain_time` - (Optional) The maintenance window of the cluster. Format: `hh:mmZ-hh:mmZ`.
* `mode` - (Required) The mode of the cluster. Valid values: `reserver`, `flexible`.
* `modify_type` - (Optional) The modify type.
* `payment_type` - (Optional, ForceNew) The payment type of the resource. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `pay_type` - (Optional, ForceNew) Deprecated field. Use `payment_type` instead.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when `payment_type` is `Subscription`. Valid values: [1~9], 12, 24, 36.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`. Default to `NotRenewal`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include `0.0.0.0/0`, `10.23.12.24` (IP), and `10.23.12.24/24` (CIDR mode. `/24` represents the length of the prefix in an IP address. The range of the prefix length is `[1,32]`).
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the resource.

* `instance_inner_connection` - (Optional, Deprecated) Field `instance_inner_connection` is deprecated and will be removed in a future release. Please use new field `connection_string` instead.
* `instance_inner_port` - (Optional, Deprecated) Field `instance_inner_port` is deprecated and will be removed in a future release. Please use new field `port` instead.

### Removing `alibabacloudstack_adb_cluster` from your configuration

The `alibabacloudstack_adb_cluster` resource allows you to manage your ADB cluster, but Terraform cannot destroy it if your cluster type is pre-paid (post-paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the ADB Console.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `storage_resource` - The specifications of storage resources in elastic mode.
* `storage_type` - Reserved parameters, not involved.
* `db_node_class` - The DB node class.
* `db_node_count` - The number of DB nodes.
* `executor_count` - The number of nodes.
* `db_node_storage` - The storage capacity of each DB node.
* `description` - The description of the DBCluster.
* `maintain_time` - The maintenance window of the cluster.
* `payment_type` - The payment type of the resource.
* `pay_type` - Deprecated field. Use `payment_type` instead.
* `resource_group_id` - The ID of the resource group.
* `security_ips` - Security IPs.
* `status` - The status of the resource.
* `instance_inner_connection` - The endpoint of the cluster.
* `instance_inner_port` - The internal port of the cluster.
* `instance_vpc_id` - The VPC ID.
* `connection_string` - The connection string of the cluster.
* `port` - The port of the cluster.
* `zone_id` - The zone ID of the resource.

## Import

AnalyticDB for MySQL (ADB) DBCluster can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_adb_db_cluster.example <id>
```