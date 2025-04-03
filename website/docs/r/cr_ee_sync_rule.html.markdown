---
subcategory: "Container Registry (CR)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_sync_rule"
sidebar_current: "docs-alibabacloudstack-resource-cr-ee-sync-rule"
description: |-
  Provides a Alibabacloudstack resource to manage Container Registry Enterprise Edition sync rules.
---

# alibabacloudstack_cr_ee_sync_rule

This resource will help you to manager Container Registry Enterprise Edition sync rules.

For information about Container Registry Enterprise Edition sync rules and how to use it, see [Create a Sync Rule](https://www.alibabacloud.com/help/doc-detail/145280.htm)



-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_cr_ee_sync_rule" "default" {
  instance_id           = "my-source-instance-id"
  namespace_name        = "my-source-namespace"
  name                  = "test-sync-rule"
  target_region_id      = "cn-hangzhou"
  target_instance_id    = "my-target-instance-id"
  target_namespace_name = "my-target-namespace"
  tag_filter            = ".*"
  repo_name             = "my-source-repo"
  target_repo_name      = "my-target-repo"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition source instance.
* `namespace_name` - (Required, ForceNew) Name of Container Registry Enterprise Edition source namespace. It can contain 2 to 30 characters.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition sync rule. 
* `target_region_id` - (Required, ForceNew) The target region to be synchronized.
* `target_instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition target instance to be synchronized.
* `target_namespace_name` - (Required, ForceNew) Name of Container Registry Enterprise Edition target namespace to be synchronized. It can contain 2 to 30 characters.
* `tag_filter` - (Required, ForceNew) The regular expression used to filter image tags for synchronization in the source repository.
* `repo_name` - (Optional, ForceNew) Name of the source repository which should be set together with `target_repo_name`, if empty means that the synchronization scope is the entire namespace level.
* `target_repo_name` - (Optional, ForceNew) Name of the target repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of Container Registry Enterprise Edition sync rule. The value is in format `{instance_id}:{namespace_name}:{rule_id}`.
* `rule_id` - The uuid of Container Registry Enterprise Edition sync rule. 
* `sync_direction` - `FROM` or `TO`, the direction of synchronization. `FROM` means source instance, `TO` means target instance.
* `sync_scope` - `REPO` or `NAMESPACE`,the scope that the synchronization rule applies.

## Import

Container Registry Enterprise Edition sync rule can be imported using the id. Format to `{instance_id}:{namespace_name}:{rule_id}`, e.g.

```
$ terraform import alibabacloudstack_cr_ee_sync_rule.default `cri-xxx:my-namespace:crsr-yyy`
```