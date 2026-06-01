---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-cspprivate-hsm-groups"
description: |-
  提供CSP Private HSM Group列表。
---

# alibabacloudstack_cspprivate_hsm_groups

该数据源用于获取专有云中可用的CSP Private HSM Group列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_cspprivate_hsm_groups" "example" {
  ids = ["my-hsm-group"]
}

output "hsm_groups" {
  value = data.alibabacloudstack_cspprivate_hsm_groups.example.groups
}
```

## 参数说明

以下参数支持配置：

* `ids` - (可选) 用于按组名过滤结果的HSM组名列表。

## 属性参考

以下属性会被导出：

* `ids` - HSM组ID列表。
* `groups` - CSP Private HSM Group列表。每个元素包含以下属性：
  * `id` - HSM组的ID（等同于组名）。
  * `group_name` - HSM组的名称。
  * `status` - HSM组的当前状态。
  * `unique_id` - HSM组的唯一ID标识符。
  * `security_level_tag` - HSM组的安全级别标签。
  * `create_time` - HSM组的创建时间，ISO 8601格式。
  * `hsm_count` - HSM组中的HSM数量。
  * `vpc_id` - HSM组所属的VPC ID。
  * `update_time` - HSM组的最后更新时间，ISO 8601格式。
  * `zone_ids` - HSM组所在的可用区ID列表。
