---
subcategory: "对象存储 OSS"
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

## 参数说明

* `ids` - (可选, 变更后新建) 指定端点ID列表，用于精确匹配特定区域的端点信息。(*可选*)
* `name_regex` - (可选, 变更后新建) 通过正则表达式过滤端点名称，用于筛选符合条件的端点。(*可选*)
* `region_id` - (可选, 变更后新建) 指定端点所在地域ID。(*可选*)

## 属性导出

* `ids` - 匹配到的端点ID列表。
* `endpoints` - 匹配到的端点信息列表，每个元素包含以下属性：
  * `id` - 端点ID，对应集群标识。
  * `cluster` - 集群标识。
  * `ha_apsara_stack` - 是否为高可用ApsaraStack。
  * `api_zonelocal_endpoint` - Zone本地API端点。
  * `api_zonelocal_public_endpoint` - Zone本地公共API端点。
  * `oss_public_endpoint` - OSS公共端点。
  * `real_zone` - 实际Zone信息。
  * `oss_ha_enable_single_cluster_access` - 是否启用单集群访问的高可用OSS。
  * `oss_cs_public_endpoint` - OSS容器服务公共端点。
  * `oss_unique_domain` - 是否使用唯一域名。
  * `cluster_name` - 集群名称。
  * `is_master_zone` - 是否为主Zone。
  * `location` - 地理位置信息。
  * `oss_endpoint` - OSS端点地址。
  * `oss_suffix` - OSS后缀信息。