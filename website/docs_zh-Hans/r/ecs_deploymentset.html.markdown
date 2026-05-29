---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_deploymentset"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-deploymentset"
description: |-
  编排云服务器（Ecs）部署集
---

# alibabacloudstack_ecs_deploymentset
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_deployment_set`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）部署集。

## 示例用法

### 基础用法：

```terraform
variable "name" {
  default = "tf-testaccecsdeploymentset3890"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  deployment_set_name = var.name
  description         = "This is a test deployment set."
  domain              = "Default"
  granularity         = "Host"
  strategy            = "Availability"
}
```

## 参数说明

支持以下参数：

* `deployment_set_name` - (可选，变更时重建) 部署集的名称。名称必须为2到128个字符长度，可以包含字母、数字、冒号（`:`）、下划线（`_`) 和连字符（`-`）。
* `description` - (可选) 部署集的描述。描述必须为2到256个字符长度，并且不能以 `http://` 或 `https://` 开头。
* `domain` - (可选，变更时重建) 部署域。有效值：`Default`。
* `granularity` - (可选，变更时重建) 部署粒度。有效值：`Host`、`Rack`、`Switch`。
* `strategy` - (可选，变更时重建) 部署策略。有效值：`Availability`、`LooseDispersion`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 部署集的ID。

## Import

ECS 部署集可以通过 DeploymentSetId 导入，例如：

```
$ terraform import alibabacloudstack_ecs_deployment_set.example ds-bp67acfmxazb4ph****
```