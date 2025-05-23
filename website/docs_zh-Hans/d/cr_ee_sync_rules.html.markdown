---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_sync_rules"
sidebar_current: "docs-alibabacloudstack-datasource-cr-ee-sync-rules"
description: |-
  查询容器镜像企业版同步规则
---

# alibabacloudstack_cr_ee_sync_rules

根据指定过滤条件列出当前凭证权限可以访问的容器镜像企业版同步规则列表。



## 示例用法

```
# 声明数据源
data "alibabacloudstack_cr_ee_sync_rules" "my_sync_rules" {
  instance_id = "cri-xxx"
  namespace_name = "test-namespace"
  repo_name = "test-repo"
  target_instance_id = "cri-yyy"
  name_regex = "test-rule"
}

output "output" {
  value = data.alibabacloudstack_cr_ee_sync_rules.my_sync_rules.rules.*.id
}
```

## 参数参考

支持以下参数：

* `instance_id` - (必填) Container Registry Enterprise Edition 本地实例的 ID。
* `namespace_name` - (可选) Container Registry Enterprise Edition 本地命名空间的名称。
* `repo_name` - (可选) Container Registry Enterprise Edition 本地仓库的名称。
* `target_instance_id` - (可选) Container Registry Enterprise Edition 目标实例的 ID。
* `name_regex` - (可选) 用于通过同步规则名称过滤结果的正则表达式字符串。
* `ids` - (可选) 用于通过同步规则 ID 过滤结果的 ID 列表。
* `names` - (可选) 用于通过同步规则名称过滤结果的名称列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 匹配的 Container Registry Enterprise Edition 同步规则列表。其元素是一个同步规则的 UUID。
* `names` - 同步规则名称列表。
* `rules` - 匹配的 Container Registry Enterprise Edition 同步规则列表。每个元素包含以下属性：
  * `id` - Container Registry Enterprise Edition 同步规则的 ID。
  * `name` - Container Registry Enterprise Edition 同步规则的名称。
  * `region_id` - Container Registry Enterprise Edition 本地实例所在的区域。
  * `instance_id` - Container Registry Enterprise Edition 本地实例的 ID。
  * `namespace_name` - Container Registry Enterprise Edition 本地命名空间的名称。
  * `repo_name` - Container Registry Enterprise Edition 本地仓库的名称。
  * `target_region_id` - Container Registry Enterprise Edition 目标实例所在的区域。
  * `target_instance_id` - Container Registry Enterprise Edition 目标实例的 ID。
  * `target_namespace_name` - Container Registry Enterprise Edition 目标命名空间的名称。
  * `target_repo_name` - Container Registry Enterprise Edition 目标仓库的名称。
  * `tag_filter` - 用于在源存储库中筛选要同步的镜像标签的正则表达式。
  * `sync_direction` - `FROM` 或 `TO`，同步的方向。`FROM` 表示本地实例是源实例。`TO` 表示本地实例是要同步的目标实例。
  * `sync_scope` - `REPO` 或 `NAMESPACE`，同步规则适用的范围。
  * `sync_trigger` - `PASSIVE` 或 `INITIATIVE`，配置的触发同步规则的策略。