---
subcategory: "Object Storage Service"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_oss_clusters"
description: |-
  Provides a list of OSS Clusters to the user.
---

# alibabacloudstack_oss_clusters

This data source provides the OSS Clusters of the current Alibaba Cloud user.

## Example Usage

```hcl
data "alibabacloudstack_oss_clusters" "example" {
  ids = ["cluster-id-1"]
}

output "first_cluster_id" {
  value = data.alibabacloudstack_oss_clusters.example.clusters.0.id
}
```

```hcl
data "alibabacloudstack_oss_clusters" "filtered" {
  name_regex = "^master-"
}

output "filtered_clusters" {
  value = data.alibabacloudstack_oss_clusters.filtered.clusters
}
```

## Argument Reference

* `ids` - (Optional, ForceNew) Specify a list of cluster IDs to match the cluster information of specific clusters exactly. (*Optional*)
* `name_regex` - (Optional, ForceNew) Filter cluster names by regular expression to select clusters that meet the specified conditions. (*Optional*)

## Attributes Reference

* `ids` - A list of cluster IDs that match.
* `clusters` - A list of cluster information that matches, each element contains the following attributes:
  * `id` - Cluster ID, corresponding to cluster identifier.
  * `cluster` - Cluster identifier.
  * `ha_apsara_stack` - Whether it is a high availability ApsaraStack.
  * `api_zonelocal_endpoint` - Zone local API endpoint.
  * `api_zonelocal_public_endpoint` - Zone local public API endpoint.
  * `oss_public_endpoint` - OSS public endpoint.
  * `real_zone` - Actual zone information.
  * `oss_ha_enable_single_cluster_access` - Whether to enable high availability OSS with single cluster access.
  * `oss_cs_public_endpoint` - OSS container service public endpoint.
  * `oss_unique_domain` - Whether to use a unique domain name.
  * `cluster_name` - Cluster name.
  * `is_master_zone` - Whether it is a master zone.
  * `location` - Geographic location information.
  * `oss_endpoint` - OSS endpoint address.
  * `oss_suffix` - OSS suffix information.
