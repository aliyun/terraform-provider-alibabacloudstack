---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_logon_policies"
sidebar_current: "docs-alibabacloudstack-ascm-logon-policies"
description: |-
  查询登录策略列表
---
# alibabacloudstack_ascm_logon_policies

根据指定过滤条件列出当前凭证权限可以查看的登录策略列表。

## 示例用法

```
resource "alibabacloudstack_ascm_logon_policy" "default" {
  name="Test_login_policy"
  description="testing policy"
  rule="ALLOW"
}
output "login" {
  value = alibabacloudstack_ascm_logon_policy.default.id
}
data "alibabacloudstack_ascm_logon_policies" "default"{
  name = alibabacloudstack_ascm_logon_policy.default.name
}
output "policies" {
  value = data.alibabacloudstack_ascm_logon_policies.default.*
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 登录策略的ID列表。
* `name` - (可选) 登录策略名称。
* `name_regex` - (可选) 用于按名称过滤登录策略的正则表达式字符串。
* `description` - (可选) 登录策略描述。
* `rule` - (可选) 登录策略规则。
* `ip_ranges` - (可选) 登录策略的IP范围。
* `ids` - (可选) 登录策略的ID列表。

## 属性参考

导出以下属性：

* `name` - 登录策略的名称。
* `policies` - 登录策略列表。每个元素包含以下属性：
    * `id` - 登录策略的ID。
    * `name` - 登录策略的名称。
    * `description` - 登录策略的描述。
    * `rule` - 登录策略的规则。
    * `ip_range` - 登录策略的IP范围。
    * `end_time` - 登录策略的结束时间。
    * `start_time` - 登录策略的开始时间。
    * `login_policy_id` - 登录策略的登录策略ID。
    * `ip_ranges` - 登录策略的IP范围列表。