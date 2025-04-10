---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_group"
sidebar_current: "docs-alibabacloudstack-resource-dns-group"
description: |-
  编排DNS组资源。
---

# alibabacloudstack_dns_group
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_alidns_domaingroup`

使用Provider配置的凭证在指定的资源集下编排DNS组资源。

## 示例用法

```
# 添加一个新的域名组。
resource "alibabacloudstack_dns_group" "group" {
  name = "testgroup"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 域名组的名称。  
* `child` - (可选) 域名组的子属性，用于定义该组下的子组或关联关系。  
* `propreties` - (必填) 域名组的属性集合，用于描述该组的详细配置信息。

## 属性说明

导出以下属性：

* `id` - 域名组的唯一标识符（ID）。  
* `name` - 域名组的名称，与输入参数中的 `name` 相同。  
* `propreties` - 域名组的属性集合，包含该组的详细配置信息。