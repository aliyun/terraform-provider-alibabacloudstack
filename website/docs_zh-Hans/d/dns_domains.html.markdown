---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domains"
sidebar_current: "docs-alibabacloudstack-datasource-dns-domains"
description: |-
    查询DNS域名
---

# alibabacloudstack_dns_domains

根据指定过滤条件列出当前凭证权限可以访问的域名列表。

## 示例用法

```
resource "alibabacloudstack_dns_domain" "default" {
  domain_name = "domaintest."
  remark = "testing Domain"
}
data "alibabacloudstack_dns_domains" "default"{
  domain_name   = alibabacloudstack_dns_domain.default.domain_name
}
output "domains" {
  value = data.alibabacloudstack_dns_domains.default.*
}
```

## 参数说明

支持以下参数：

* `domain_name` - （可选）按域名过滤结果的正则表达式字符串。  
* `ids` - （可选）域名ID列表，用于精确匹配特定域名。  
* `resource_group_id` - （可选，变更时重建）DNS所属的资源组ID，用于筛选属于某个资源组的域名。  
* `dns_servers` - （可选）分析系统中域名的DNS服务器列表，用于进一步过滤域名。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 域名ID列表，包含所有符合条件的域名的唯一标识符。  
* `names` - 域名名称列表，包含所有符合条件的域名的名称。  
* `domains` - 域名详细信息列表，每个元素包含以下属性：  
  * `domain_id` - 域名的唯一标识符（ID）。  
  * `domain_name` - 域名的名称。  
  * `dns_servers` - 分析系统中该域名对应的DNS服务器列表。  
  * `resource_group_id` - 该域名所属的资源组ID。  