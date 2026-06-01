---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_domain"
sidebar_current: "docs-Alibabacloudstack-resource-dns-domain"
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
  domain_name = "tfacc-test."
  remark      = "测试域名"
}
```

## 参数说明

支持以下参数：

* `domain_name` - (必填，变更时重建) 域名名称。修改此参数将强制创建新资源。
* `remark` - (可选) 域名的备注信息。

## 属性说明

导出以下属性：

* `id` - 资源的ID。格式：`domain_name:domain_id`。
* `domain_id` - 域名ID。
* `domain_name` - 域名名称。

## Import

DNS域名可以通过 `domain_name` 和 `domain_id` 用冒号分隔进行导入，例如：

```
$ terraform import alibabacloudstack_dns_domain.example example.com:12345678
```
