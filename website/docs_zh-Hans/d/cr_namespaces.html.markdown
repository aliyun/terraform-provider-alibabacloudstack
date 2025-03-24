---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_namespaces"
sidebar_current: "docs-alibabacloudstack-datasource-cr-namespaces"
description: |-
  查询容器镜像仓库命名空间
---

# alibabacloudstack_cr_namespaces

根据指定过滤条件列出当前凭证权限可以访问的容器镜像仓库命名空间列表。



## 示例用法

```
# 声明数据源
data "alibabacloudstack_cr_namespaces" "my_namespaces" {
  name_regex  = "my-namespace"
  output_file = "my-namespace-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_namespaces.my_namespaces.namespaces}"
}
```

## 参数参考

支持以下参数：

* `name_regex` - (可选) 用于通过命名空间名称过滤结果的正则表达式字符串。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配的容器镜像仓库命名空间列表。其元素是一个命名空间名称。
* `names` - 命名空间名称列表。
* `namespaces` - 匹配的容器镜像仓库命名空间列表。每个元素包含以下属性：
  * `name` - 容器镜像仓库命名空间的名称。
  * `auto_create` - 布尔值，当设置为 true 时，在推送新镜像时会自动创建仓库。如果设置为 false，则需要在推送之前创建镜像仓库。
  * `default_visibility` - `PUBLIC` 或 `PRIVATE`，此命名空间中的默认仓库可见性。