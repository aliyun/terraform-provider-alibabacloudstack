---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_namespace"
sidebar_current: "docs-alibabacloudstack-resource-cr-ee-namespace"
description: |-
  编排容器镜像企业版命名空间
---

# alibabacloudstack_cr_ee_namespace

使用Provider配置的凭证在指定的资源集下编排容器镜像企业版命名空间。

有关容器镜像企业版命名空间的信息以及如何使用它，请参阅 [创建命名空间](https://www.alibabacloud.com/help/doc-detail/145483.htm)。



-> **注意：** 在使用此资源之前，您需要在容器镜像企业版控制台中设置您的注册表密码。

## 示例用法

### 基础用法

```
resource "alibabacloudstack_cr_ee_namespace" "my-namespace" {
  instance_id        = "cri-xxx"
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## 参数说明

支持以下参数：

* `instance_id` - (必填，变更时重建) 容器镜像企业版实例的 ID。
* `name` - (必填，变更时重建) 容器镜像企业版命名空间的名称。它可以包含 2 到 30 个字符。
* `auto_create` - (必填) 布尔值，当设置为 true 时，在推送新镜像时会自动创建仓库。如果设置为 false，则在推送前需要创建仓库以存储镜像。
* `default_visibility` - (必填) 命名空间内存储库的默认可见性设置。有效值为 `PUBLIC` 或 `PRIVATE`。

## 属性说明

导出以下属性：

* `id` - 容器镜像企业版命名空间的 ID。格式为 `{instance_id}:{namespace}`。

## 导入

容器镜像企业版命名空间可以使用 `{instance_id}:{namespace}` 导入，例如：

```
$ terraform import alibabacloudstack_cr_ee_namespace.default cri-xxx:my-namespace
```