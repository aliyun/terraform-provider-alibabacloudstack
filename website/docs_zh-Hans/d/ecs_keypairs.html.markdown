---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_keypairs"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-keypairs"
description: |- 
  查询云服务器密钥对
---

# alibabacloudstack_ecs_keypairs
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_key_pairs`

根据指定过滤条件列出当前凭证权限可以访问的云服务器密钥对列表。

## 示例用法

```hcl
# 创建一个密钥对
resource "alibabacloudstack_key_pair" "default" {
  key_name = "exampleKeyPair"
}

# 检索与 name_regex 匹配的所有密钥对
data "alibabacloudstack_ecs_keypairs" "default" {
  name_regex = "${alibabacloudstack_key_pair.default.key_name}"
  ids        = ["${alibabacloudstack_key_pair.default.id}"]
  finger_print = "MD5FingerprintOfKeyPair"
  tags = {
    Environment = "Test"
  }
}

output "key_pairs" {
  value = data.alibabacloudstack_ecs_keypairs.default.key_pairs
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 用于通过名称筛选结果密钥对的正则表达式字符串。可以用于模糊匹配密钥对名称。
* `ids` - (可选) 密钥对 ID 列表。如果提供，仅返回这些 ID 的密钥对。
* `finger_print` - (可选) 密钥对的指纹。基于 RFC 4716 中定义的公钥指纹格式使用消息摘要算法 5 (MD5)。有关更多信息，请参见 [RFC 4716](https://tools.ietf.org/html/rfc4716)。
* `tags` - (可选) 要分配给资源的标签映射。可以通过标签进一步筛选密钥对。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 密钥对名称列表。
* `key_pairs` - 密钥对列表。每个元素包含以下属性：
  * `id` - 密钥对的 ID。
  * `key_name` - 密钥对的名称。
  * `finger_print` - 密钥对的指纹。基于 RFC 4716 中定义的公钥指纹格式使用消息摘要算法 5 (MD5)。有关更多信息，请参见 [RFC 4716](https://tools.ietf.org/html/rfc4716)。
  * `instances` - 已绑定到此密钥对的 ECS 实例列表。每个实例包括以下属性：
    * `availability_zone` - ECS 实例所在的可用区 ID。
    * `instance_id` - ECS 实例的 ID。
    * `instance_name` - ECS 实例的名称。
    * `vswitch_id` - 连接到 ECS 实例的交换机 ID。
    * `public_ip` - ECS 实例的公网 IP 地址或 EIP。
    * `private_ip` - ECS 实例的私网 IP 地址。
  * `tags` - (可选) 分配给密钥对的标签映射。