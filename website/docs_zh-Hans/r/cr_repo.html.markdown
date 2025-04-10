---
subcategory: "容器镜像服务 (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_repo"
sidebar_current: "docs-alibabacloudstack-resource-container-registry"
description: |-
  编排容器镜像服务的存储库

---

# alibabacloudstack_cr_repo

使用Provider配置的凭证在指定的资源集下编排容器镜像服务的存储库。


## 示例用法

### 基础用法

```
resource "alibabacloudstack_cr_namespace" "my-namespace" {
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alibabacloudstack_cr_repo" "my-repo" {
  namespace = alibabacloudstack_cr_namespace.my-namespace.name
  name      = "my-repo"
  summary   = "this is summary of my new repo"
  repo_type = "PUBLIC"
  detail    = "this is a public repo"
}
```

## 参数说明

以下是支持的参数：

* `namespace` - (必填，变更时重建) 容器镜像服务命名空间名称，存储库所在位置。
* `name` - (必填，变更时重建) 容器镜像服务存储库名称。
* `summary` - (必填) 存储库的一般信息描述。它可以包含1到80个字符。
* `repo_type` - (必填) 存储库的可见性类型，取值为 `PUBLIC` 或 `PRIVATE`。
* `detail` - (可选) 存储库的具体信息描述。支持MarkDown格式，长度限制为2000。
* `summary` - (必填) 存储库的一般信息描述。它可以包含1到100个字符。

## 属性说明

以下属性将被导出：

* `id` - 容器镜像服务存储库的ID。格式为 `命名空间/存储库`。
* `domain_list` - 存储库域名列表，包含以下字段：
  * `public` - 公网访问域名。
  * `internal` - 内网访问域名，仅在部分区域可用。
  * `vpc` - VPC内网访问域名。