---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_deploymentsets"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-deploymentsets"
description: |- 
  Provides a list of ecs deploymentsets owned by an alibabacloudstack account.
---

# alibabacloudstack_ecs_deploymentsets
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_deployment_sets`

This data source provides a list of ECS Deployment Sets in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage:

```terraform
data "alibabacloudstack_ecs_deploymentsets" "ids" {
  ids = ["example_id"]
}

output "ecs_deployment_set_id_1" {
  value = data.alibabacloudstack_ecs_deploymentsets.ids.sets.0.id
}

data "alibabacloudstack_ecs_deploymentsets" "nameRegex" {
  name_regex = "^my-DeploymentSet"
}

output "ecs_deployment_set_id_2" {
  value = data.alibabacloudstack_ecs_deploymentsets.nameRegex.sets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Deployment Set IDs.  
* `name_regex` - (Optional, ForceNew) A regex string used to filter results by Deployment Set name.  
* `deployment_set_name` - (Optional, ForceNew) The name of the deployment set.  
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `Availability`.  

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `sets` - A list of ECS Deployment Sets. Each element contains the following attributes:
  * `create_time` - The time when the deployment set was created.
  * `id` - The ID of the Deployment Set.
  * `deployment_set_id` - The ID of the Deployment Set.
  * `deployment_set_name` - The name of the deployment set.
  * `description` - The description of the deployment set.
  * `domain` - The deployment domain.
  * `granularity` - The deployment granularity.
  * `instance_amount` - The number of instances in the deployment set.
  * `instance_ids` - The IDs of the instances in the deployment set.
  * `strategy` - The deployment strategy.