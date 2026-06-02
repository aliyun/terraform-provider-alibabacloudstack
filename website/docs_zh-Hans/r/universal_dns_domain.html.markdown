---
subcategory: "跨云域名服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_domain"
sidebar_current: "docs-Alibabacloudstack-resource-universal-dns-domain"
description: |-
  跨云解析域名资源
---

# alibabacloudstack_universal_dns_domain

创建和管理跨云解析域名资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc35847"
}


resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "tf-testacc35847.example."
  remark = "Created by Terraform"
}
```

## 参数说明

支持以下参数：

* `name` - (必填, 变更时重建) 跨云共享域名名称（必须以"."结尾）。
* `remark` - (可选) 域名备注信息。

## 属性说明

以下属性被导出：

* `id` - 域名ID。
* `create_timestamp` - 创建时间戳（秒）。
* `record_count` - 解析记录集总数。
* `update_timestamp` - 修改时间戳（秒）。