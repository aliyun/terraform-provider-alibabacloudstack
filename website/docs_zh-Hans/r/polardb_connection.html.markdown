---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbconnection"
sidebar_current: "docs-Alibabacloudstack-polardb-dbconnection"
description: |-
  编排polardb数据库连接
---

# alibabacloudstack_polardb_dbconnection

使用Provider配置的凭证在指定的资源集编排polardb数据库连接。

## 示例用法

```hcl
resource "alibabacloudstack_polardb_dbconnection" "default" {
  instance_id       = "your_polardb_instance_id"
  connection_prefix = "your_connection_prefix"
  port              = "3306"
}
```

## 参数说明

以下是支持的参数：

* `instance_id` -(必填，变更时重建) - PolarDB实例的ID。
* `connection_prefix` -(可选，变更时重建) - 连接字符串的前缀。它的长度必须为1到31个字符，并且可以包含数字、字母、下划线（_）和连字符（-）。它必须以字母、数字或中文字符开头。如果未指定，默认为实例ID加上`tf`。
* `port` -(可选) - 连接的端口号。默认值为`3306`，有效值范围为`1024`到`65535`。


## 属性说明

除了上述列出的参数外，还导出了以下属性：

* `connection_string` - PolarDB实例的完整连接字符串，用于客户端连接数据库。
* `ip_address` - 数据库连接的IP地址，表示PolarDB实例的访问入口地址。

## 导入
PolarDB数据库连接可以通过id导入，例如：

```sh
$ terraform import alibabacloudstack_polardb_dbconnection.example <instance_id>:<connection_prefix>
```

## 示例：

```sh
$ terraform import alibabacloudstack_polardb_dbconnection.example polardb-instance-123456:my_connection_prefix
```