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

## 参数说明

支持以下参数：

* `name` - (必填) 登录策略的名称。  
* `description` - (可选) 登录策略的描述信息，用于说明该策略的具体用途或规则。  
* `rule` - (可选) 登录策略的规则。有效值为：  
  * `ALLOW`：允许登录。
  * `DENY`：拒绝登录。


## 属性说明

导出以下属性：

* `name` - 登录策略的名称。  
* `description` - 登录策略的描述信息。  
* `rule` - 登录策略的规则，表示该策略是允许 (`ALLOW`) 还是拒绝 (`DENY`) 登录。  
* `policy_id` - 创建的登录策略的唯一标识符（ID），可用于后续管理和引用该策略。