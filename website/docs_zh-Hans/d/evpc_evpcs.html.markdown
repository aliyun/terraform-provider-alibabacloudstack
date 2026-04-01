---
subcategory: "EasyAI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_evpc_evpcs"
sidebar_current: "docs-Alibabacloudstack-datasource-evpc-evpcs"
description: |- 
  获取EVPC列表
---

# alibabacloudstack_evpc_evpcs

该数据源根据指定的过滤器获取AlibabacloudStack账户中的EVPC（Elastic Virtual Private Cloud，弹性虚拟私有云）列表。

## 示例用法

```hcl
# 声明数据源
data "alibabacloudstack_evpc_evpcs" "example" {
  evpc_name = "my-evpc"
}

output "evpc_ids" {
  value = data.alibabacloudstack_evpc_evpcs.example.ids
}

output "evpc_names" {
  value = data.alibabacloudstack_evpc_evpcs.example.names
}
```

## 参数说明

支持以下参数：

* `evpc_name` - (选填) 用于过滤结果的EVPC名称。
* `status` - (选填) 用于过滤结果的EVPC状态。
* `name_regex` - (选填) 用于按EVPC名称过滤结果的正则表达式字符串。
* `ids` - (选填) 用于过滤结果的EVPC ID列表。如果未指定，将考虑所有EVPC。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `names` - 匹配的EVPC名称列表。
* `ids` - 匹配的EVPC ID列表。
* `evpcs` - 匹配的EVPC列表。每个元素包含以下属性：
  * `evpc_id` - EVPC的ID。
  * `evpc_name` - EVPC的名称。
  * `status` - EVPC的状态。
  * `description` - EVPC的描述。
  * `cidr` - EVPC的CIDR块。
  * `tenant_id` - EVPC的租户ID。
  * `department` - EVPC的部门ID。
  * `department_name` - EVPC的部门名称。
  * `region_id` - EVPC的地域ID。
  * `resource_group` - EVPC的资源组ID。
  * `resource_group_name` - EVPC的资源组名称。
  * `cluster_id` - EVPC的集群ID。
  * `ascm_create_user` - 创建EVPC的用户。
  * `create_time` - EVPC的创建时间。
  * `update_time` - EVPC的最后更新时间。
