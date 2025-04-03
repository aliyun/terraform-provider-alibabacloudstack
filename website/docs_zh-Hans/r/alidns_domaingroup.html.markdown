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

## 参数参考

支持以下参数：

* `name` - (必填) 域名组的名称。    

* `child` - (可选) 域名组的子属性。 

## 属性参考

导出以下属性：

* `id` - 组ID。
* `name` - 组名称。

* `propreties` - 域名组的属性。