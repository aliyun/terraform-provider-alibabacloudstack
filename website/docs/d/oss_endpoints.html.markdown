---
subcategory: "OSS"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_oss_endpoints"
sidebar_current: "docs-alibabacloudstack-datasource-oss-endpoints"
description: |-
  Provides a list of OSS Endpoints to the user.
---

# alibabacloudstack_oss_endpoints

This data source provides the OSS Endpoints of the current Alibaba Cloud user.

## Example Usage

```hcl
data "alibabacloudstack_oss_endpoints" "example" {
  ids = ["cn-beijing"]
}

output "first_endpoint_id" {
  value = data.alibabacloudstack_oss_endpoints.example.endpoints.0.id
}
```

```hcl
data "alibabacloudstack_oss_endpoints" "filtered" {
  name_regex = "^cn-(beijing|shanghai)"
}

output "filtered_endpoints" {
  value = data.alibabacloudstack_oss_endpoints.filtered.endpoints
}
```

## Argument Reference

* `ids` - (Optional, ForceNew) Specify a list of endpoint IDs to match the endpoint information of specific regions exactly. (*Optional*)
* `name_regex` - (Optional, ForceNew) Filter endpoint names by regular expression to select endpoints that meet the specified conditions. (*Optional*)
* `region_id` - (Optional, ForceNew) Specify the region ID of the endpoint. (*Optional*)

## Attributes Reference

* `ids` - A list of endpoint IDs that match.
* `endpoints` - A list of endpoint information that matches, each element contains the following attributes:
  * `id` - Endpoint ID, corresponding to cluster identifier.
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