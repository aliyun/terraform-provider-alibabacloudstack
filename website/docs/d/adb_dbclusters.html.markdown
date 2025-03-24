---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_dbclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-adb-dbclusters"
description: |- 
  Provides a list of adb dbclusters owned by an Alibabacloudstack account.
---

# alibabacloudstack_adb_dbclusters
-> **NOTE:** Alias name has: `alibabacloudstack_adb_clusters` `alibabacloudstack_adb_db_clusters`

This data source provides a list of adb dbclusters in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_adb_dbclusters" "example" {
  description_regex = "example-cluster"
  resource_group_id = "rg-abc1234567890"
  status           = "Running"
}

output "first_adb_dbcluster_id" {
  value = data.alibabacloudstack_adb_dbclusters.example.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the DBCluster. This can be used to filter clusters by their description.
* `description_regex` - (Optional, ForceNew) A regex string to filter results by DBCluster description.
* `enable_details` - (Optional) Default to `false`. Set it to `true` to output more details about resource attributes.
* `ids` - (Optional, ForceNew) A list of DBCluster IDs. This can be used to filter clusters by their unique identifiers.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group. This can be used to filter clusters belonging to a specific resource group.
* `status` - (Optional, ForceNew) The status of the resource. Valid values include `Creating`, `Running`, `Stopping`, `Stopped`, and `Starting`.
* `tags` - (Optional) A map of tags assigned to the cluster. This can be used to filter clusters by their tags.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `descriptions` - A list of DBCluster descriptions.
* `clusters` - A list of ADB DbClusters. Each element contains the following attributes:
  * `commodity_code` - The name of the service associated with the DBCluster.
  * `connection_string` - The endpoint of the cluster.
  * `create_time` - The creation time of the DBCluster.
  * `db_cluster_category` - The category of the DBCluster.
  * `db_cluster_id` - The unique identifier of the DBCluster.
  * `db_cluster_network_type` - The network type of the DBCluster.
  * `network_type` - The network type of the DBCluster.
  * `db_cluster_type` - The type of the DBCluster.
  * `db_cluster_version` - The version of the DBCluster.
  * `db_node_class` - The class of the DB nodes in the DBCluster.
  * `db_node_count` - The number of DB nodes in the DBCluster.
  * `db_node_storage` - The storage size of each DB node in the DBCluster.
  * `description` - The description of the DBCluster.
  * `disk_type` - The type of the disk used by the DBCluster.
  * `dts_job_id` - The ID of the data synchronization task in Data Transmission Service (DTS). This parameter is valid only for analytic instances.
  * `elastic_io_resource` - The elastic I/O resource allocated to the DBCluster.
  * `engine` - The database engine used by the DBCluster.
  * `executor_count` - The number of executor nodes in the DBCluster. These nodes are used for data computing in elastic mode.
  * `expire_time` - The expiration time of the DBCluster.
  * `expired` - Indicates whether the DBCluster has expired.
  * `id` - The ID of the DBCluster.
  * `lock_mode` - The lock mode of the DBCluster.
  * `lock_reason` - The reason why the DBCluster is locked.
  * `maintain_time` - The maintenance window of the DBCluster.
  * `payment_type` - The payment type of the DBCluster.
  * `charge_type` - The charge type of the DBCluster.
  * `port` - The port that is used to access the DBCluster.
  * `rds_instance_id` - The ID of the ApsaraDB RDS instance from which data is synchronized to the DBCluster. This parameter is valid only for analytic instances.
  * `resource_group_id` - The ID of the resource group to which the DBCluster belongs.
  * `security_ips` - A list of IP addresses allowed to access all databases of the DBCluster.
  * `status` - The status of the DBCluster.
  * `storage_resource` - The specifications of storage resources in elastic mode. Increasing these resources can improve the read and write performance of the DBCluster.
  * `tags` - The tags assigned to the DBCluster.
  * `vpc_cloud_instance_id` - The VPC cloud instance ID associated with the DBCluster.
  * `vpc_id` - The VPC ID associated with the DBCluster.
  * `vswitch_id` - The VSwitch ID associated with the DBCluster.
  * `zone_id` - The zone ID of the DBCluster.
  * `region_id` - The region ID of the DBCluster.