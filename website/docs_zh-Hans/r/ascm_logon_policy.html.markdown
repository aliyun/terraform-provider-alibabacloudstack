---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_logon_policy"
sidebar_current: "docs-alibabacloudstack_ascm_logon_policy"
description: |-
  编排登录策略资源
---
# alibabacloudstack_ascm_logon_policy

使用Provider配置的凭证编排登录策略资源。

### 基础用法

```
resource "alibabacloudstack_ascm_logon_policy" "login" {
  name="test_foo"
  description="testing purpose"
  rule="ALLOW"
}
```

## 参数参考

支持以下参数：

* `name` - (必填) 登录策略的名称。

* `description` - (可选) 登录策略的描述。

* `rule` - (可选) 登录策略的规则。有效值：Allow 和 Deny。


## 属性参考

导出以下属性：

* `name` - 登录策略的名称。
* `description` - 登录策略的描述。
* `rule` - 登录策略的规则。
* `policy_id` - 创建的登录策略的ID。