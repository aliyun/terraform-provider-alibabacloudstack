---
subcategory: "ACK"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_account"
sidebar_current: "docs-Alibabacloudstack-datasource-account"
description: |- 
  查询阿里云账户的ID
---

# alibabacloudstack_account

当前Provider中所配置的阿里云账户的唯一标识符查询。

## 示例用法

```hcl
data "alibabacloudstack_account" "current" {}

output "account_id" {
  value = data.alibabacloudstack_account.current.id
}
```

## 参数参考
此数据源没有可配置的参数。它基于已验证的账户检索信息。

## 属性参考
以下属性被导出：

`id` - 阿里巴巴云账户的唯一标识符。

## 常见注意事项
确保您的提供商配置正确并包含必要的凭据。
此数据源可用于获取有关您的阿里巴巴云账户的各种详细信息。

## 注意事项
此数据源主要用于获取当前已认证用户的账户ID。
确保提供商已正确配置有效的凭据以访问阿里巴巴云API。

## 导入
alibabacloudstack_account 数据源不支持导入，因为它直接基于已验证会话检索账户信息。