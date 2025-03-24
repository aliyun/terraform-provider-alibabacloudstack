---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_repo"
sidebar_current: "docs-alibabacloudstack-resource-cr-ee-repo"
description: |-
  编排容器镜像企业版仓库
---

# alibabacloudstack_cr_ee_repo
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cr_repository`

使用Provider配置的凭证在指定的资源集下编排容器镜像企业版仓库。

有关容器镜像企业版仓库的更多信息以及如何使用它，请参阅 [创建仓库](https://www.alibabacloud.com/help/doc-detail/145291.htm)。



-> **注意**：在使用此资源之前，您需要在容器镜像企业版控制台中设置您的注册表密码。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alibabacloudstack_cr_ee_repo" "my-repo" {
  instance_id = alibabacloudstack_cr_ee_namespace.my-namespace.instance_id
  namespace   = alibabacloudstack_cr_ee_namespace.my-namespace.name
  name        = "my-repo"
  summary     = "this is summary of my new repo"
  repo_type   = "PUBLIC"
  detail      = "this is a public repo"
}
```

## 参数参考

支持以下参数：

* `instance_id` - (必填，变更时重建) 容器镜像企业版实例的 ID。
* `namespace` - (必填，变更时重建) 容器镜像企业版命名空间的名称。该命名空间下包含仓库。它可以包含 2 到 30 个字符。
* `name` - (必填，变更时重建) 容器镜像企业版仓库的名称。它可以包含 2 到 64 个字符。
* `summary` - (必填) 仓库的一般信息。它可以包含 1 到 100 个字符。
* `repo_type` - (必填) `PUBLIC` 或 `PRIVATE`，表示仓库的可见性类型。
* `detail` - (可选) 仓库的具体信息。支持 MarkDown 格式，长度限制为 2000。
* `repo_id` - (可选) 容器镜像企业版仓库的 uuid。

## 属性参考

导出以下属性：

* `id` - 容器镜像企业版仓库的资源 ID。格式为 `{instance_id}:{namespace}:{repository}`。
* `repo_id` - 容器镜像企业版仓库的 uuid。

## 导入

可以使用 `{instance_id}:{namespace}:{repository}` 导入容器镜像企业版仓库，例如：

```bash
$ terraform import alibabacloudstack_cr_ee_repo.default `cri-xxx:my-namespace:my-repo`
```