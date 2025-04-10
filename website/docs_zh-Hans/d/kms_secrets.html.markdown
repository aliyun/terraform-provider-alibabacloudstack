---
subcategory: "KMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kms_secrets"
sidebar_current: "docs-alibabacloudstack-datasource-kms-secrets"
description: |-
    查询KMS密钥
---

# alibabacloudstack_kms_secrets

根据指定过滤条件列出当前凭证权限可以访问的KMS密钥列表。

## 示例用法

```
# 声明数据源
data "alibabacloudstack_kms_secrets" "kms_secrets_ds" {
  fetch_tags = true
  name_regex = "name_regex"
  tags = {
    "k-aa" = "v-aa",
    "k-bb" = "v-bb"
  }
}

output "first_secret_id" {
  value = "${data.alibabacloudstack_kms_secrets.kms_secrets_ds.secrets.0.id}"
}
```

## 参数说明

支持以下参数：

* `fetch_tags` - (可选) 是否在返回值中包含预定义的资源标签。默认为 `false`。
* `ids` - (可选) KMS 密钥 ID 列表。其值与 KMS `secret_name` 相同。
* `name_regex` - (可选) 用于通过 KMS `secret_name` 过滤结果的正则表达式字符串。
* `tags` - (可选) 分配给资源的标签映射。
* `names` - (可选) KMS 密钥名称列表。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - KMS 密钥 ID 列表。其值与 KMS `secret_name` 相同。
* `names` - KMS 密钥名称列表。
* `secrets` - KMS 密钥列表。每个元素包含以下属性：
  * `id` - KMS 密钥 ID。其值与 KMS `secret_name` 相同。
  * `secret_name` - KMS 密钥名称。
  * `planned_delete_time` - 计划删除时间。如果密钥未设置计划删除，则该值为空。
  * `tags` - 分配给资源的标签映射（计算得出）。