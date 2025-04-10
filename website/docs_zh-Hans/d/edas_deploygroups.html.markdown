---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_deploygroups"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-deploygroups"
description: |- 
  查询企业级分布式应用服务部署组
---

# alibabacloudstack_edas_deploygroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_deploy_groups`

根据指定过滤条件列出当前凭证权限可以访问的企业级分布式应用服务部署组列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-edas-deploy-groups7396"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = "${var.name}"
  cluster_type = 2
  network_mode = 2
  vpc_id       = "${alibabacloudstack_vpc.default.id}"
  region_id    = "cn-neimeng-env30-d01"
}

resource "alibabacloudstack_edas_application" "default" {
  application_name = "${var.name}"
  cluster_id      = "${alibabacloudstack_edas_cluster.default.id}"
  package_type    = "JAR"
  build_pack_id   = "15"
}

resource "alibabacloudstack_edas_deploy_group" "default" {
  app_id     = "${alibabacloudstack_edas_application.default.id}"
  group_name = "${var.name}"
}

data "alibabacloudstack_edas_deploy_groups" "default" {
  name_regex = "${alibabacloudstack_edas_deploy_group.default.group_name}"
  app_id     = "${alibabacloudstack_edas_application.default.id}"

  output_file = "deploygroups_output.txt"
}

output "first_group_name" {
  value = data.alibabacloudstack_edas_deploy_groups.default.groups[0].group_name
}
```

## 参数说明

以下参数是支持的：

* `app_id` - (必填，变更时重建) 要检索部署组的应用程序的ID。
* `name_regex` - (可选，变更时重建) 用于按部署组名称过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 部署组ID列表。
* `names` - 部署组名称列表。
* `groups` - 部署组列表。每个元素包含以下属性：
  * `group_id` - 部署组的ID。
  * `group_name` - 部署组的名称。长度不能超过64个字符。
  * `group_type` - 部署组的类型。有效值：
    - `0`: 默认分组。
    - `1`: 流量管理未启用灰度。
    - `2`: 流量管理启用灰度。
  * `create_time` - 创建时间的时间戳。
  * `update_time` - 更新时间的时间戳。
  * `app_id` - 应用程序的ID。
  * `cluster_id` - 集群的ID。
  * `package_version_id` - 部署包的版本ID。
  * `app_version_id` - 应用程序部署记录的版本ID。