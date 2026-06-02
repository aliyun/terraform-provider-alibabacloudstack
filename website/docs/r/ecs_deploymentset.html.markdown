---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_deploymentset"
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
  domain              = "Default"
  granularity         = "Host"
  strategy            = "Availability"
}
```

## Argument Reference

The following arguments are supported:

* `deployment_set_name` - (Optional, ForceNew) The name of the deployment set. The name must be 2 to 128 characters in length and can contain letters, digits, colons (`:`), underscores (`_`), and hyphens (`-`).
* `description` - (Optional) The description of the deployment set. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `domain` - (Optional, ForceNew) The deployment domain. Valid values: `Default`.
* `granularity` - (Optional, ForceNew) The deployment granularity. Valid values: `Host`, `Rack`, `Switch`.
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `Availability`, `LooseDispersion`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the deployment set.

## Import

ECS Deployment Set can be imported using the DeploymentSetId, e.g.

```
$ terraform import alibabacloudstack_ecs_deployment_set.example ds-bp67acfmxazb4ph****
```