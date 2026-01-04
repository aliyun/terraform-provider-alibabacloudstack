---
subcategory: "Container Registry"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_artifact_lifecycle_rule"
sidebar_current: "docs-Alibabacloudstack-cr-cr_ee_artifact_lifecycle_rule"
description: |-
  Manages artifact lifecycle rules for ACR Enterprise Edition to automatically clean up expired image tags and manifests.
---

# alibabacloudstack_cr_ee_artifact_lifecycle_rule

Manages artifact lifecycle rules for ACR Enterprise Edition, used to automatically clean up expired image tags and manifests.

## Example Usage

### Basic Usage

```hcl

variable "name" {
  default = "tf-testacc-cree-rule-3358413"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_namespace" "default2" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = "${var.name}2"
  auto_create        = true
  default_visibility = "PRIVATE"
}



resource "alibabacloudstack_cr_ee_attestor_lifecycle_rule" "default" {
  retention_tag_count = "30"
  tag_regexp          = "release-v.*"
  enable_delete_tag   = "true"
  namespace_name      = alibabacloudstack_cr_ee_namespace.default.name
  recent_pull_keep    = 30
  recent_push_keep    = 20
  scope               = "NAMESPACE"
  instance_id         = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `retention_tag_count` - (Required) The number of tags to retain. Specifies the number of most recently pushed tags to keep, must be greater than 0.
* `scope` - (Required) The scope of the rule. Valid values: `REPO` (repository level), `NAMESPACE` (namespace level). When set to `REPO`, `repo_name` must be specified; when set to `NAMESPACE`, `namespace_name` must be specified.
* `instance_id` - (Required, Forces new resource) The ID of the Container Registry Enterprise Edition instance. Changing the instance ID will cause the resource to be recreated.
* `enable_delete_tag` - (Optional) Whether to enable tag deletion. Default is `false`. When set to `true`, tags matching the rule will be deleted.
* `namespace_name` - (Optional) The namespace name. Required when `scope` is `NAMESPACE`.
* `recent_pull_keep` - (Optional) The number of most recently pulled images to retain. Specifies how many most recently pulled images to keep, 0 means no retention.
* `recent_push_keep` - (Optional) The number of most recently pushed images to retain. Specifies how many most recently pushed images to keep, 0 means no retention.
* `repo_name` - (Optional) The repository name. Required when `scope` is `REPO`.
* `tag_regexp` - (Optional) The regular expression for tag matching. Used to match tags to retain (e.g., `release-v.*`). Once set, the rule will be executed automatically.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `<instance_id>:<rule_id>`.
* `auto` - Whether the rule is executed automatically. Automatically set to `true` when `tag_regexp` is set.
* `create_time` - The creation time of the rule, in Unix timestamp (milliseconds).
* `enable_delete_untagged_manifest` - Whether to enable deletion of untagged manifests.
* `modified_time` - The modification time of the rule, in Unix timestamp (milliseconds).
* `rule_id` - The ID of the retention policy rule.
* `schedule` - The scheduling method. Fixed as `MANUAL`, indicating manual execution of the rule.