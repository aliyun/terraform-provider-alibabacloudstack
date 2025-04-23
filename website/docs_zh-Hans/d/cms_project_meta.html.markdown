---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_project_meta"
sidebar_current: "docs-alibabacloudstack-datasource-cms-project-meta"
description: |-
    查询云监控项目元数据
---

# alibabacloudstack_cms_project_meta

根据指定过滤条件列出当前凭证权限可以访问的云监控项目元数据列表

## 示例用法

基本用法

```
data "alibabacloudstack_cms_project_meta" "default" {
  name_regex = "OSS"
}

output "project_meta" {
  value = data.alibabacloudstack_cms_project_meta.default.*
}
```

## 参数参考

支持以下参数：

* `name_regex` - (可选，变更时重建) 一个正则表达式字符串，用于通过项目元数据描述过滤结果。
* `resources` - (可选) CMS 项目元数据列表。每个元素包含以下属性：

## 属性参考

导出以下属性：

* `resources` - CMS 项目元数据列表。每个元素包含以下属性：
    * `description` - 项目元数据的描述。
    * `labels` - CMS 项目元数据的标签。度量标准的标签作为触发该度量标准警报的特殊标记使用。格式为 `[{"name":"标签名称","value":"标签值"}, {"name":"标签名称","value":"标签值"}]`。
        * `name` - 标签的名称。
        * `value` - 标签的值。
    * `namespace` - 服务的命名空间，用于区分不同的服务。通常，其值为 acs_ + 服务名称缩写 的格式。