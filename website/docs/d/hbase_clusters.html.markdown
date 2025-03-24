---
subcategory: "HBase"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hbase_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-hbase-clusters"
description: |- 
  Provides a list of hbase clusters owned by an alibabacloudstack account.
---

# alibabacloudstack_hbase_clusters
-> **NOTE:** Alias name has: `alibabacloudstack_hbase_instances`

This data source provides a list of HBase clusters in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_hbase_clusters" "hbase" {
  name_regex        = "tf_testAccHBase"
  availability_zone = "cn-shenzhen-b"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the cluster name. This allows filtering clusters based on their names using regular expressions.
* `ids` - (Optional) The IDs list of HBase clusters. This can be used to filter clusters by their unique identifiers.
* `availability_zone` - (Optional) The availability zone where the HBase clusters reside. Use this parameter to filter clusters within a specific zone.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - The list of HBase cluster IDs.
* `names` - The list of HBase cluster names.
* `instances` - A list of HBase clusters. Each element contains the following attributes:
  * `id` - The ID of the HBase cluster.
  * `name` - The name of the HBase cluster.
  * `region_id` - The region ID where the cluster is located.
  * `zone_id` - The availability zone ID where the cluster is located.
  * `engine` - The engine type of the cluster (e.g., HBase).
  * `engine_version` - The version of the engine used by the cluster.
  * `network_type` - The network type of the cluster, such as `Classic` or `VPC`.
  * `master_instance_type` - The instance type for master nodes (e.g., `hbase.sn2.2xlarge`).
  * `master_node_count` - The number of master nodes in the cluster.
  * `core_instance_type` - The instance type for core nodes (e.g., `hbase.sn2.4xlarge`).
  * `core_node_count` - The number of core nodes in the cluster.
  * `core_disk_type` - The disk type for core nodes, such as `Cloud_SSD` or `Cloud_Efficiency`.
  * `core_disk_size` - The disk size (in GB) for core nodes.
  * `vpc_id` - The VPC ID associated with the cluster.
  * `vswitch_id` - The VSwitch ID associated with the cluster.
  * `pay_type` - The billing method of the cluster. Possible values include `PostPaid` (Pay-As-You-Go) and `PrePaid` (yearly or monthly subscription).
  * `status` - The current status of the cluster.
  * `backup_status` - The backup status of the cluster.
  * `created_time` - The creation time of the cluster.
  * `expire_time` - The expiration time of the cluster (if applicable).
  * `deletion_protection` - Indicates whether deletion protection is enabled for the cluster.
  * `tags` - A mapping of tags assigned to the cluster.