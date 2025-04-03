---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_alias"
sidebar_current: "docs-alibabacloudstack-resource-kms-alias"
description: |-
  编排加密主密钥 (CMK) 的别名
---

# alibabacloudstack_kms_alias

使用Provider配置的凭证在指定的资源集下编排加密主密钥 (CMK) 的别名。


## 示例用法

### 基础用法

```
resource "alibabacloudstack_kms_key" "key" {}

resource "alibabacloudstack_kms_alias" "alias" {
  alias_name = "alias/test_kms_alias"
  key_id     = alibabacloudstack_kms_key.key.id
}
```

## 参数参考

支持以下参数：

* `alias_name` - (必填，变更时重建) CMK 的别名。可以使用别名调用 `Encrypt`、`GenerateDataKey` 和 `DescribeKey`。字符长度(不包括前缀)：最小长度为 1 个字符，最大长度为 255 个字符。必须包含前缀 `alias/`。
* `key_id` - (必填) 密钥的 ID。


-> **注意：** 每个别名只能代表一个主密钥 (CMK)。

-> **注意：** 在同一区域的同一用户内，别名不可重复。

-> **注意：** 可以使用 `UpdateAlias` 更新别名与主密钥 (CMK) 的映射关系。


## 属性参考

* `id` - 别名的 ID。