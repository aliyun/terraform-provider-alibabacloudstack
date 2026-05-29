---
subcategory: "裸金属算力平台 BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_keypairs"
sidebar_current: "docs-Alibabacloudstack-datasource-bcmp-keypairs"
description: |- 
  提供 AlibabaCloudStack 账户拥有的 BCMP 密钥对列表。
---

# alibabacloudstack_bcmp_keypairs

此数据源根据指定的过滤器提供 AlibabaCloudStack 账户中的 BCMP 密钥对列表。

## 示例用法

```hcl
# 声明资源
resource "alibabacloudstack_bcmp_keypair" "default" {
  key_pair_name = "exampleKeyPair"
}

# 检索匹配 name_regex 的所有密钥对
data "alibabacloudstack_bcmp_keypairs" "default" {
  name_regex = "${alibabacloudstack_bcmp_keypair.default.key_pair_name}"
}

output "key_pairs" {
  value = data.alibabacloudstack_bcmp_keypairs.default.key_pairs
}
```

## 参数说明

支持以下参数：

* `name_regex` - (选填) 用于按名称过滤结果密钥对的正则表达式字符串。
* `ids` - (选填) 密钥对 ID 列表。如果提供，将只返回具有这些 ID 的密钥对。
* `finger_print` - (选填) 密钥对的指纹。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `names` - 密钥对名称列表。
* `key_pairs` - 密钥对列表。每个元素包含以下属性：
  * `id` - 密钥对的 ID。
  * `key_pair_name` - 密钥对的名称。
  * `finger_print` - 密钥对的指纹。
  * `region` - 密钥对的区域。
  * `create_time` - 密钥对的创建时间。
  * `update_time` - 密钥对的更新时间。
