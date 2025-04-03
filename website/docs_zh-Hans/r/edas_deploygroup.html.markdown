---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_deploygroup"
sidebar_current: "docs-Alibabacloudstack-edas-deploygroup"
description: |- 
  编排企业级分布式应用服务（Edas）部署组
---

# alibabacloudstack_edas_deploygroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_edas_deploy_group`

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas）部署组。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-edasdeploygroupbasic4916"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name = var.name
  cluster_type = 2
  network_mode = 2
  vpc_id       = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_edas_application" "default" {
  application_name = var.name
  cluster_id      = alibabacloudstack_edas_cluster.default.id
  package_type    = "JAR"
  build_pack_id   = "15"
}

resource "alibabacloudstack_edas_deploy_group" "default" {
  app_id       = alibabacloudstack_edas_application.default.id
  group_name   = var.name
  group_type   = 2 # 流量管理启用灰度
}
```

## 参数参考

支持以下参数：

* `app_id` - (必填，变更时重建) 应用的唯一标识符。这是 EDAS 中应用程序的唯一 ID。
* `group_name` - (必填，变更时重建) 部署组的名称。它必须在同一个应用内是唯一的，并且创建后无法修改。
* `group_type` - (可选，变更时重建) 部署组的类型。有效值为：
  - `0`：默认分组。
  - `1`：流量管理未启用灰度。
  - `2`：流量管理启用灰度。此选项允许使用分阶段发布和流量管理功能。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 部署组资源的唯一标识符。其格式为 `<app_id>:<group_name>:<group_id>`。
* `group_type` - 部署组的类型。该属性反映了创建时设置的值，并表示分组的行为：
  - `0`：默认分组。
  - `1`：流量管理未启用灰度。
  - `2`：流量管理启用灰度。

此属性有助于识别部署组的配置及其在 EDAS 环境中的功能。