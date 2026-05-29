---
subcategory: "云数据库 Redis 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_connection"
sidebar_current: "docs-Alibabacloudstack-resource-redis-connection"
description: |- 
  在指定的资源集中编排Redis互联网连接字符串。
---

# alibabacloudstack_redis_connection

在指定的资源集中编排Redis互联网连接字符串。

> **注意：** 该资源也可以使用以下别名进行引用：
> - `alibabacloudstack_kvstore_connection`

## 示例用法

### 基础用法

```terraform
variable "name" {
    default = "tf-testaccredisconnection"
}

resource "alibabacloudstack_redis_connection" "default" {
  connection_string_prefix = var.name
  instance_id              = local.instance_id
  port                     = "6379"
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填，ForceNew) Redis实例的ID。修改此参数将强制重新创建资源。
* `connection_string_prefix` - (必填) 连接字符串的前缀。前缀长度可以为8到64个字符，可以包含小写字母和数字。必须以小写字母开头。
* `port` - (必填) Redis实例的服务端口号。有效范围：`1024` 到 `65535`。
* `connection_string` - (可选，Computed) Redis实例的连接字符串。此属性由API返回，无法手动设置。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - Redis实例的ID，与 `instance_id` 相同。
* `connection_string` - Redis实例的完整连接字符串。

### 超时配置

`timeouts` 块允许您为某些操作指定[超时时间](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)：

* `create` - (默认2分钟) 用于创建Redis连接（直到达到初始`正常`状态）。
* `update` - (默认2分钟) 用于更新Redis连接（直到达到初始`正常`状态）。
* `delete` - (默认2分钟) 用于删除Redis连接（直到达到初始`正常`状态）。

## 导入

Redis连接可以使用实例ID进行导入，例如：

```bash
$ terraform import alibabacloudstack_redis_connection.example r-abc12345678
```
