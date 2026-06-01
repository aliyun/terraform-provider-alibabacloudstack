---
subcategory: "云数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_readwrite_splitting_connection"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-readwrite-splitting-connection"
description: |-
  提供PolarDB读写分离连接资源。
---

# alibabacloudstack\_polardb\_readwrite\_splitting\_connection

提供PolarDB读写分离连接资源，允许您为PolarDB集群配置读写分离。

-> **注意：** 版本 v3.20.0+ 可用。

## 示例用法

```hcl
resource "alibabacloudstack_polardb_readwrite_splitting_connection" "example" {
  instance_id       = "pc-xxxxxxxxxxxxx"
  connection_id     = "pe-xxxxxxxxxxxxx"
  distribution_type = "Standard"
  max_delay_time    = 30
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填, ForceNew) PolarDB集群ID。
* `connection_id` - (必填, ForceNew) 代理终端ID。
* `distribution_type` - (必填) 读请求分发类型。取值：`Standard`、`Custom`。
  - `Standard`：系统根据只读节点的权重自动分发读请求。
  - `Custom`：您可以自定义每个只读节点的权重。
* `weight` - (可选, 由API返回) 每个只读节点的权重。当 `distribution_type` 设置为 `Custom` 时，此参数为必填。值是一个映射，键为只读节点ID，值为权重（0-100）。
* `max_delay_time` - (可选) 只读节点的最大延迟时间阈值。单位：秒。默认值：30。如果只读节点的延迟超过此阈值，读请求将不会分发到该节点。

## 属性说明

导出以下属性：

* `id` - 资源ID，与 `instance_id` 相同。
* `connection_string` - 读写分离终端的连接地址。
* `port` - 读写分离终端的端口号。

## Import

PolarDB读写分离连接可以通过实例ID导入，例如：

```
$ terraform import alibabacloudstack_polardb_readwrite_splitting_connection.example pc-xxxxxxxxxxxxx
```
