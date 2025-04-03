---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_repos"
sidebar_current: "docs-alibabacloudstack-datasource-cr-ee-repos"
description: |-
  查询容器镜像企业版仓库
---

# alibabacloudstack_cr_ee_repos
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cr_repositories`

根据指定过滤条件列出当前凭证权限可以访问的容器镜像企业版仓库列表。



## 示例用法

```
# 声明数据源
data "alibabacloudstack_cr_ee_repos" "my_repos" {
  instance_id = "cri-xx"
  name_regex  = "my-repos"
  output_file = "my-repo-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_ee_repos.my_repos.repos}"
}
```

## 参数参考

以下是支持的参数：

* `instance_id` - (必填) 容器镜像企业版实例的ID。
* `namespace` - (可选) 容器镜像企业版命名空间名称，其中包含要查询的仓库。
* `ids` - (可选) 按仓库ID过滤结果的ID列表。
* `name_regex` - (可选) 用于按仓库名称过滤结果的正则表达式字符串。
* `enable_details` - (可选) 布尔值，默认为false，仅导出仓库属性。如果需要此仓库所属的标签，请设置为true。详见`tags`在属性中。
* `names` - (可选) 仓库名称列表。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `ids` - 匹配的容器镜像企业版仓库的列表。其元素是一个仓库ID。
* `names` - 仓库名称列表。
* `repos` - 匹配的容器镜像企业版仓库列表。每个元素包含以下属性：
  * `instance_id` - 容器镜像企业版实例的ID。
  * `namespace` - 存储库所在的容器镜像企业版命名空间名称。
  * `id` - 容器镜像企业版存储库的ID。
  * `name` - 容器镜像企业版存储库的名称。
  * `summary` - 仓库的一般信息。
  * `repo_type` - `PUBLIC` 或 `PRIVATE`，仓库的可见性。
  * `tags` - 属于此仓库的镜像标签列表。每个包含几个属性，详见`标签块`。
    * `tag` - 此镜像的标签。
    * `image_id` - 此镜像的ID。
    * `digest` - 此镜像的摘要。
    * `status` - 此镜像的状态。
    * `image_size` - 此镜像的大小，以字节为单位。
    * `image_update` - 此镜像的最后更新时间，以纳秒为单位的Unix时间。
    * `image_create` - 此镜像的创建时间，以纳秒为单位的Unix时间。
* `names` - 仓库名称列表。