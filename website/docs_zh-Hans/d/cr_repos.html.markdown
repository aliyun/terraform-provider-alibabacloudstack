---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_repos"
sidebar_current: "docs-alibabacloudstack-datasource-cr-repos"
description: |-
  查询容器镜像仓库
---

# alibabacloudstack_cr_repos

根据指定过滤条件列出当前凭证权限可以访问的容器镜像服务中的镜像仓库列表。



## 示例用法

```
# 声明数据源
data "alibabacloudstack_cr_repos" "my_repos" {
  name_regex  = "my-repos"
  output_file = "my-repo-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_repos.my_repos.repos}"
}
```

## 参数参考

支持以下参数：

* `namespace` - (可选) 容器镜像命名空间名称，其中包含要查询的仓库。
* `name_regex` - (可选) 用于通过仓库名称过滤结果的正则表达式字符串。
* `enable_details` - (可选) 布尔值，默认为 false，仅导出仓库属性。如果需要属于该仓库的域名列表和标签，请设置为 true。请参阅属性中的 `tags`。

* `ids` - (可选) 匹配的容器镜像仓库列表。其元素被设置为 `names`。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 匹配的容器镜像仓库列表。其元素被设置为 `names`。
* `names` - 仓库名称列表。
* `repos` - 匹配的容器镜像仓库列表。每个元素包含以下属性：
  * `namespace` - 镜像仓库所在的容器镜像命名空间名称。
  * `name` - 容器镜像仓库名称。
  * `summary` - 仓库的一般信息。
  * `repo_type` - `PUBLIC` 或 `PRIVATE`，仓库的可见性。
  * `domain_list` - 仓库域名列表。
    * `public` - 公网端点域名。
    * `internal` - 内网端点域名，仅在某些区域可用。
    * `vpc` - VPC 端点域名。
  * `tags` - 属于此仓库的镜像标签列表。每个包含若干属性，详见 `Block Tag`。

  * `summary` - 仓库的一般信息。
  * `repo_type` - `PUBLIC` 或 `PRIVATE`，仓库的可见性。
  * `domain_list` - 仓库域名列表。

### Block Tag

* `tag` - 此镜像的标签。
* `image_id` - 此镜像的 ID。
* `digest` - 此镜像的摘要。
* `status` - 此镜像的状态。
* `image_size` - 此镜像的大小，以字节为单位。
* `image_update` - 此镜像的最后更新时间，Unix 时间戳，以纳秒为单位。
* `image_create` - 此镜像的创建时间，Unix 时间戳，以纳秒为单位。