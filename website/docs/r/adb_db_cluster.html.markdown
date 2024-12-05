---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_db_cluster"
sidebar_current: "docs-alibabacloudstack-resource-adb-db-cluster"
description: |-
  Provides a Alibabacloudstack AnalyticDB for MySQL (ADB) DBCluster resource.
---

# alibabacloudstack\_adb\_db\_cluster

Provides a AnalyticDB for MySQL (ADB) DBCluster resource.

For information about AnalyticDB for MySQL (ADB) DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/help/en/doc-detail/190519.htm).

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = var.creation
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id       = alibabacloudstack_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alibabacloudstack_adb_db_cluster" "this" {
  db_cluster_category = "Basic"
  db_cluster_class    = "C8"
  db_node_count       = "2"
  db_node_storage     = "200"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alibabacloudstack_vswitch.default.id
  description         = "Test new adb again."
  maintain_time       = "23:00Z-00:00Z"
  security_ips      = ["10.168.1.12", "10.168.1.11"]
  cluster_type      = "analyticdb"
  cpu_type          = "intel"
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when `payment_type` is `Subscription`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`. Default to `1`.
* `compute_resource` - (Optional) The specifications of computing resources in elastic mode. The increase of resources can speed up queries. AnalyticDB for MySQL automatically scales computing resources. For more information, see [ComputeResource](https://www.alibabacloud.com/help/en/doc-detail/144851.htm)
* `db_cluster_category` - (Required) The db cluster category. Valid values: `Basic`, `Cluster`, `MixedStorage`.
* `db_cluster_class` - (Deprecated) It duplicates with attribute db_node_class and is deprecated from 1.121.2.
* `db_cluster_version` - (Optional, ForceNew) The db cluster version. Value options: `3.0`, Default to `3.0`.
* `db_node_class` - (Optional, Computed) The db node class. For more information, see [DBClusterClass](https://help.aliyun.com/document_detail/190519.html)
* `db_node_count` - (Optional) The db node count.
* `db_node_storage` - (Optional) The db node storage.
* `description` - (Optional, Computed) The description of DBCluster.
* `executor_count` - The number of nodes. The node resources are used for data computing in elastic mode.
* `instance_inner_port` - (Optional, ForceNew)The endpoint's port of the cluster.
* `maintain_time` - (Optional, Computed) The maintenance window of the cluster. Format: hh:mmZ-hh:mmZ.
* `mode` - (Required) The mode of the cluster. Valid values: `reserver`, `flexible`.
* `modify_type` - (Optional) The modify type.
* `pay_type` - (Optional, Computed, ForceNew) Field `pay_type` has been deprecated. New field `payment_type` instead.
* `payment_type` - (Optional, Computed, ForceNew) The payment type of the resource. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when `payment_type` is `Subscription`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_ips` - (Optional, Computed) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `storage_type` - (Optional) The type of storage media that is used for the instance.
* `storage_resource` - The specifications of storage resources in elastic mode. The resources are used for data read and write operations. The increase of resources can improve the read and write performance of your cluster. For more information, see [Specifications](https://www.alibabacloud.com/help/en/doc-detail/144851.htm).
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the resource.
* `cluster_type` - (Required) The cluster type of the resource. Valid values: `analyticdb`, `AnalyticdbOnPanguHybrid`.
* `cpu_type` - (Required) The cpu type of the resource.Valid values: `intel`.
-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~30 minutes. Please make full preparation before changing them.

### Removing alibabacloudstack_adb_cluster from your configuration
 
The alibabacloudstack_adb_cluster resource allows you to manage your adb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the adb Console.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.
* `instance_inner_connection` - The endpoint of the cluster.
* `instance_vpc_id` - The Vpc id.
* `status` - The status of the resource.

## Import

AnalyticDB for MySQL (ADB) DBCluster can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_adb_db_cluster.example <id>
```
