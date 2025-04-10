---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_instances"
sidebar_current: "docs-alibabacloudstack-datasource-cr-ee-instances"
description: |-
  查询容器镜像企业版实例
---

# alibabacloudstack_cr_ee_instances

根据指定过滤条件列出当前凭证权限可以访问的容器镜像企业版实例列表。



## 示例用法

```
# 声明数据源
data "alibabacloudstack_cr_ee_instances" "my_instances" {
  name_regex  = "my-instances"
  output_file = "my-instances-json"
}

output "output" {
  value = "${data.alibabacloudstack_cr_ee_instances.my_instances.instances}"
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 按实例ID过滤结果的ID列表。
* `name_regex` - (可选) 按实例名称过滤结果的正则表达式字符串。
* `enable_details` - (可选，1.132.0版本及以上可用) 默认为 `true`。将其设置为 true 可以输出实例授权令牌。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 匹配的容器镜像企业版实例列表。其元素是一个实例的UUID。
* `names` - 实例名称列表。
* `instances` - 匹配的容器镜像企业版实例列表。每个元素包含以下属性：
  * `id` - 容器镜像企业版实例的唯一标识符（UUID）。
  * `name` - 容器镜像企业版实例的名称。
  * `region` - 容器镜像企业版实例所在的区域。
  * `specification` - 容器镜像企业版实例的规格，例如性能等级或资源配置。
  * `namespace_quota` - 实例允许创建的最大命名空间数量。
  * `namespace_usage` - 当前实例中已经创建的命名空间数量。
  * `repo_quota` - 实例允许创建的最大仓库数量。
  * `repo_usage` - 当前实例中已经创建的仓库数量。
  * `vpc_endpoints` - 在VPC网络上访问实例时使用的域名列表。
  * `public_endpoints` - 在公网网络上访问实例时使用的域名列表。
  * `authorization_token` - 登录注册表时使用的密码（仅当 `enable_details` 设置为 true 时返回）。
  * `temp_username` - 登录注册表时使用的临时用户名（仅当 `enable_details` 设置为 true 时返回）。
  * `output_file` - 数据源查询结果保存的文件名（在运行 `terraform plan` 后生成）。