---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_deploymentset"
sidebar_current: "docs-Alibabacloudstack-ecs-deploymentset"
description: |- 
  Provides a ecs Deploymentset resource.
---

# alibabacloudstack_ecs_deploymentset
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_deployment_set`

Provides a ecs Deploymentset resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccecsdeploymentset3890"
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  deployment_set_name = var.name
  description         = "This is a test deployment set."
  domain             = "Default"
  granularity        = "Host"
  strategy           = "Availability"
  on_unable_to_redeploy_failed_instance = "CancelMembershipAndStart"

  tags = {
    CreatedBy = "Terraform"
    Env       = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `deployment_set_name` - (Optional, ForceNew) The name of the deployment set. The name must be 2 to 128 characters in length and can contain letters, digits, colons (`:`), underscores (`_`), and hyphens (`-`). It must start with a letter and cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the deployment set. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `domain` - (Optional, ForceNew) The deployment domain. Valid value: `Default`.
* `granularity` - (Optional, ForceNew) The deployment granularity. Valid value: `Host`.
* `on_unable_to_redeploy_failed_instance` - (Optional) The action to take when an instance fails to redeploy. Valid values:
  * `CancelMembershipAndStart`: Removes the instances from the deployment set and restarts the instances immediately after the failover is complete.
  * `KeepStopped`: Keeps the instances in the abnormal state and restarts them after ECS resources are replenished.
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid value: `Availability`.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource. Each tag consists of a key-value pair.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the deployment set. This is the same as the `deployment_set_name`.