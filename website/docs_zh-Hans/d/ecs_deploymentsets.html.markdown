---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_deploymentsets"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-deploymentsets"
description: |- 
  查询云服务器部署集
---

# alibabacloudstack_ecs_deploymentsets
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_deployment_sets`

根据指定过滤条件列出当前凭证权限可以访问的云服务器部署集列表。

## 示例用法

以下示例展示了如何使用`alibabacloudstack_ecs_deploymentsets`数据源来查询部署集信息：

### 根据ID查询部署集

```terraform
variable "name" {	
	default = "tf-testAccDeploymentSet-12892"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = var.name
  description         = var.name
}

data "alibabacloudstack_ecs_deployment_sets" "by_ids" {	
	ids = [alibabacloudstack_ecs_deployment_set.default.id]
}

output "deployment_set_id_by_ids" {
	value = data.alibabacloudstack_ecs_deployment_sets.by_ids.sets[0].deployment_set_id
}
```

### 根据名称正则表达式查询部署集

```terraform
data "alibabacloudstack_ecs_deployment_sets" "by_name_regex" {
	name_regex = "^tf-testAccDeploymentSet"
}

output "deployment_set_id_by_name_regex" {
	value = data.alibabacloudstack_ecs_deployment_sets.by_name_regex.sets[0].deployment_set_id
}
```

## 参数参考

以下参数是支持的：

* `ids` - （可选，变更时重建）部署集ID列表。用于通过ID筛选部署集。
* `name_regex` - （可选，变更时重建）用于通过部署集名称筛选结果的正则表达式字符串。
* `deployment_set_name` - （可选，变更时重建）部署集的名称。用于精确匹配部署集名称。
* `strategy` - （可选，变更时重建）部署策略。有效值为`Availability`，表示可用性策略。

## 属性参考

除了上述参数外，还导出以下属性：

* `sets` - ECS部署集列表。每个元素包含以下属性：
  * `create_time` - 部署集的创建时间。
  * `id` - 部署集的唯一标识符（与`deployment_set_id`相同）。
  * `deployment_set_id` - 部署集ID。
  * `deployment_set_name` - 部署集的名称。
  * `description` - 部署集的描述信息。
  * `domain` - 部署域。
  * `granularity` - 部署粒度。有效值为`Host`，表示主机级部署。
  * `instance_amount` - 部署集内的实例数量。
  * `instance_ids` - 部署集内的实例ID列表。
  * `strategy` - 部署策略。有效值为`Availability`。