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

## 参数说明  

支持以下参数：  

* `ids` - (可选) 登录策略的ID列表。
* `name` - (可选) 登录策略名称。
* `name_regex` - (可选) 用于按名称过滤登录策略的正则表达式字符串。
* `description` - (可选) 登录策略描述。
* `rule` - (可选) 登录策略规则。
* `ip_ranges` - (可选) 登录策略的IP范围。
* `ids` - (可选) 登录策略的ID列表。

## 属性说明

导出以下属性：  

* `name` - 登录策略的名称。  
* `policies` - 登录策略列表。每个元素包含以下属性：  
    * `id` - 登录策略的唯一标识符（ID）。  
    * `name` - 登录策略的名称。  
    * `description` - 登录策略的描述信息。  
    * `rule` - 登录策略的规则，取值为 `ALLOW` 或 `DENY`，分别表示允许或拒绝登录。  
    * `ip_range` - 登录策略的IP范围，表示允许或拒绝的IP地址段。  
    * `end_time` - 登录策略的有效结束时间，格式为时间戳或日期时间字符串。  
    * `start_time` - 登录策略的有效开始时间，格式为时间戳或日期时间字符串。  
    * `login_policy_id` - 登录策略的具体登录策略ID，与 `id` 类似但可能有额外的语义区分。  
    * `ip_ranges` - 登录策略的IP范围列表，包含多个IP地址段，用于限制登录来源。  

**补充说明：**  
- `policies`：返回的是一个列表，其中每个元素代表一个登录策略的详细信息。  
- `ip_ranges`：该字段是一个数组，包含多个IP范围，用于定义允许或拒绝登录的具体IP地址段。  
- `start_time` 和 `end_time`：这两个字段定义了登录策略的有效时间段，超出此时间段的登录请求将不受该策略约束。  

请根据实际需求结合上述参数和属性进行配置和查询。