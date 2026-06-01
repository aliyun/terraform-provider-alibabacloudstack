---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxy_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-proxy-types"
description: |-
  Provides a list of PolarDB cluster proxy types.
---

# alibabacloudstack_polardb_cluster_proxy_types

This data source provides a list of PolarDB cluster proxy types in an Alibaba Cloud Stack environment.


## Example Usage

```hcl
data "alibabacloudstack_polardb_cluster_proxy_types" "default" {
  db_type    = "MySQL"
  db_version = "5.7"
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Optional, ForceNew) The database engine type. Valid values: `MySQL`, `PostgreSQL`, `Oracle`.
* `db_version` - (Optional) The database engine version.
* `ids` - (Optional, ForceNew) A list of proxy type IDs to filter results.
* `core_count` - (Optional) The number of CPU cores.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of proxy type IDs.
* `proxy_classes` - A list of proxy classes. Each element contains the following attributes:
  * `id` - The ID of the proxy type.
  * `core_count` - The number of CPU cores.
