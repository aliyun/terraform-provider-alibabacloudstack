---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack:alibabacloudstack_cr_namespace"
sidebar_current: "docs-alibabacloudstack-resource-container-registry"
description: |-
  编排容器镜像仓库命名空间
---

# alibabacloudstack_cr_namespace

使用Provider配置的凭证在指定的资源集下编排容器镜像仓库命名空间。


## 示例用法

### 基础用法

```
resource "alibabacloudstack_cr_namespace" "my-namespace" {
  name               = "my-namespace"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## 参数说明

支持以下参数：

* `name` - (必填，变更时重建) 容器镜像仓库命名空间的名称。
* `auto_create` - (必填) 布尔值，当设置为 true 时，在推送新镜像时会自动创建仓库。如果设置为 false，则在推送之前需要先创建仓库。
* `default_visibility` - (必填) `PUBLIC` 或 `PRIVATE`，此命名空间中默认的仓库可见性。

## 属性说明

导出以下属性：

* `id` - 容器镜像仓库命名空间的唯一标识符。其值与 `name` 参数相同。
```