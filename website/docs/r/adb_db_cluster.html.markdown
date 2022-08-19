---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_adb_db_cluster"
sidebar_current: "docs-apsarastack-resource-adb-db-cluster"
description: |-
  Provides a Apsarastack AnalyticDB for MySQL (ADB) DBCluster resource.
---

# apsarastack\_adb\_db\_cluster

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

data "apsarastack_zones" "default" {
  available_resource_creation = var.creation
}

resource "apsarastack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id       = apsarastack_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.apsarastack_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "apsarastack_adb_db_cluster" "this" {
  db_cluster_category = "Basic"
  db_cluster_class    = "C8"
  db_node_count       = "2"
  db_node_storage     = "200"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = apsarastack_vswitch.default.id
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
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the resource.
* `cluster_type` - (Required) The cluster type of the resource. Valid values: `analyticdb`, `AnalyticdbOnPanguHybrid`.
* `cpu_type` - (Required) The cpu type of the resource.Valid values: `intel`.

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~30 minutes. Please make full preparation before changing them.

### Removing apsarastack_adb_cluster from your configuration
 
The apsarastack_adb_cluster resource allows you to manage your adb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the adb Console.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.
* `connection_string` - The endpoint of the cluster.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when create the DBCluster.
* `delete` - (Defaults to 50 mins) Used when delete the DBCluster.
* `update` - (Defaults to 72 mins) Used when update the DBCluster.

## Import

AnalyticDB for MySQL (ADB) DBCluster can be imported using the id, e.g.

```
$ terraform import apsarastack_adb_db_cluster.example <id>
```
