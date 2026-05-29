---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_forward_domain"
sidebar_current: "docs-Alibabacloudstack-resource-dns-forward-domain"
description: |-
  云解析全局转发域名
---

# alibabacloudstack_dns_forward_domain

使用Provider配置的凭证在指定的资源集创建云解析全局转发域名。

-> **注意:** 该资源也可以使用以下别名引用：`apsarastack_dns_forward_domain`。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc14020.test."
}


resource "alibabacloudstack_dns_forward_domain" "default" {
  name         = var.name
  remark       = "test1234"
  forward_mode = "FORWARD_FIRST"
  forwarders = [
    "192.168.101.1"
  ]
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 转发域名名称。必须以"."结尾，例如"example.com."。
* `forward_mode` - (必填) 转发模式。取值：
  * `FORWARD_FIRST`: 优先转发模式 - 转发失败降级到互联网递归。
  * `FORWARD_ONLY`: 全部转发模式(建议) - 全部请求只做转发(不作任何递归)。
* `forwarders` - (必填) 转发目的IP列表。例如：`["192.168.101.1"]`。
* `remark` - (可选) 备注信息，用于描述该转发域名的用途。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 转发域名ID。
* `caller_uid` - 系统参数，标识创建该资源的调用者。
* `create_timestamp` - 创建时间戳（秒）。
* `update_timestamp` - 修改时间戳（秒）。

## Import

转发域名可以使用资源 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_dns_forward_domain.example <resource_id>
```