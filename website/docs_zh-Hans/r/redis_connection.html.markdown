---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_connection"
sidebar_current: "docs-Alibabacloudstack-redis-connection"
description: |- 
  编排Redis互联网连接字符串
---

# alibabacloudstack_redis_connection
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_kvstore_connection`

使用Provider配置的凭证在指定的资源集编排Redis互联网连接字符串。

## 示例用法

### 基础用法：

```terraform
variable "name" {
    default = "tf-testaccredisconnection21482"
}

resource "alibabacloudstack_redis_connection" "default" {
  connection_string_prefix = "testprefix" # 连接字符串的前缀，长度为8到64个字符，包含小写字母和数字，必须以小写字母开头。
  instance_id              = "r-8vb6ces3yk5huhxoek" # Redis实例的ID，一旦设置后不可更改。
  port                     = "6379" # Redis服务端口，有效范围为1024到65535。
}
```

## 参数说明

支持以下参数：
  * `connection_string_prefix` - (必填) 实例的连接地址前缀。前缀长度可以为8到64个字符，可以包含小写字母和数字。必须以小写字母开头。
  * `instance_id` - (必填，变更时重建) Redis实例的ID。一旦设置，此值不能更改。
  * `port` - (必填) Redis实例的服务端口号。有效范围是从`1024`到`65535`。
  * `connection_string` - (可选) Redis实例的连接字符串。如果未提供，系统将自动生成。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - Redis实例的ID，与`instance_id`相同。
  * `connection_string` - Redis实例的完整连接字符串，包含前缀、实例ID和端口号。

### 超时时间

`timeouts`块允许您为某些操作指定[超时时间](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)：

* `create` - (默认为2分钟)用于创建Redis连接(直到达到初始`正常`状态)。
* `update` - (默认为2分钟)用于更新Redis连接(直到达到初始`正常`状态)。
* `delete` - (默认为2分钟)用于删除Redis连接(直到达到初始`正常`状态)。

## 导入

Redis连接可以使用id导入，例如：

```bash
$ terraform import alibabacloudstack_redis_connection.example r-abc12345678
``` 