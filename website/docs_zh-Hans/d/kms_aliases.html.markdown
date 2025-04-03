---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_aliases"
sidebar_current: "docs-alibabacloudstack-datasource-kms-aliases"
description: |-
    查询KMS别名
---

# alibabacloudstack_kms_aliases

根据指定过滤条件列出当前凭证权限可以访问的KMS别名列表。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_kms_aliases" "kms_aliases" {  
  name_regex = "alias/tf-testKmsAlias_123"
}

output "first_key_id" {
  value = "${data.alibabacloudstack_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) KMS别名ID列表。其值与KMS别名名称相同。
* `name_regex` - (可选) 用于按KMS别名名称筛选结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - KMS别名ID列表。其值与KMS别名名称相同。
* `names` - KMS别名名称列表。
* `aliases` - KMS用户别名列表。每个元素包含以下属性：
  * `id` - 别名ID。其值与KMS别名名称相同。
  * `key_id` - 密钥ID。
  * `alias_name` - 别名的唯一标识符。
  * `name_regex` - 用于按KMS别名名称筛选结果的正则表达式字符串。