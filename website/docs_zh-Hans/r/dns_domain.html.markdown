---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain"
sidebar_current: "docs-alibabacloudstack-resource-dns-domain"
description: |-
  编排DNS域名
---

# alibabacloudstack_dns_domain

使用Provider配置的凭证在指定的资源集下编排DNS域名。

-> **注意：** 您要添加的域名必须已经注册，并且没有被其他账户添加。每个域名只能存在于一个唯一的组中。

## 示例用法

```
# 添加一个新的域名。
resource "alibabacloudstack_dns_domain" "default" {
  domain_name     = "starmove."
  remark   =  "测试域名"
}
output "dns" {
  value = alibabacloudstack_dns_domain.default.*
}
```

## 参数参考

支持以下参数：

* `domain_name` - (必填，变更时重建) 域名名称。此名称(不包括后缀)可以包含1到63个字符(域名主体，不包括后缀)，只能包含字母数字字符或“-”，并且不能以“-”开头或结尾，“-”不能同时出现在第3和第4个字符位置。后缀 `.sh` 和 `.tel` 不被支持。
* `group_id` - (可选) 域名将要添加的组的ID。如果不提供，则使用默认组。
* `resource_group_id` - (可选，变更时重建) DNS域名所属的资源组ID。
* `lang` - (可选) 用户语言。
* `remark` - (可选) 您的域名的备注信息。
* `domain_name` - (必填) 域名名称。

## 属性参考

导出以下属性：

* `id` - 此资源的ID。该值设置为 `domain_name`。
* `domain_id` - 域名ID。
* `dns_servers` - DNS服务器名称列表。
* `domain_name` - 域名名称。