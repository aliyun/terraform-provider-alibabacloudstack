---
subcategory: "Container Registry"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cr_ee_attestor_lifecycle_rules"
description: |-
  Retrieves lifecycle rule configurations for Alibaba Cloud Container Registry (ACR) Enterprise Edition.
---

# alibabacloudstack_cr_ee_attestor_lifecycle_rules

This data source queries image retention policy rules configured in Alibaba Cloud Container Registry (ACR) Enterprise Edition. It supports filtering rules by instance ID, namespace regular expressions, and other conditions, and returns the list of matching rules along with their detailed attributes.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-cree-rule-3142"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
  instance_id        = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  name               = var.name
  auto_create        = false
  default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_attestor_lifecycle_rule" "default" {
  instance_id         = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  scope               = "NAMESPACE"
  retention_tag_count = "30"
  tag_regexp          = "release-v.*"
  enable_delete_tag   = "true"
  namespace_name      = alibabacloudstack_cr_ee_namespace.default.name
  recent_pull_keep    = 30
  recent_push_keep    = 20
}

data "alibabacloudstack_cr_ee_attestor_lifecycle_rules" "default" {
  instance_id = data.alibabacloudstack_cr_ee_instances.default.instances.0.id
  ids         = ["${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}"]
}
```

## Argument Reference

The following arguments are used to filter query results. Arguments are sorted by type: required > optional (within the same type, sorted alphabetically).

* `instance_id` (String, Required): The ID of the Container Registry instance associated with the lifecycle rule. For example, `cri-private`.

* `enable_delete_untagged_manifest` (Boolean, Optional): Filters rules based on whether untagged manifest deletion is enabled. When set to `true`, only rules with this feature enabled are returned.

* `ids` (List, Optional): A list of rule IDs to filter results. Each ID is in the format `{instance_id}:{rule_id}`, for example, `cri-private:cralr-42x00j4drgt4fjr3`.

* `namespace_regex` (String, Optional): A regular expression used to filter namespace names. For example, `test.*` matches the namespace `testtf`.

## Attributes Reference

The following attributes are exported as read-only (`Computed`) properties. The `id` attribute is always listed first, followed by other top-level attributes sorted alphabetically. Nested attributes are sorted alphabetically under their parent attribute.

* `id` (String): The unique identifier for the data source, generated from a hash of the filter conditions.

* `names` (List): A list of matching namespace names. Each element is a string representing a namespace name that matches the filter conditions (e.g., `testtf`).

* `rules` (List): A list of matching lifecycle rules. Each rule object contains the following attributes:
  * `auto` (Boolean): Whether the rule executes automatically. `false` indicates manual triggering (e.g., `MANUAL` scheduling mode).
  * `create_time` (Integer): The rule creation time (Unix timestamp in milliseconds).
  * `enable_delete_tag` (Boolean): Whether tag deletion is enabled.
  * `enable_delete_untagged_manifest` (Boolean): Whether untagged manifest deletion is enabled.
  * `instance_id` (String): The ID of the associated Container Registry instance (e.g., `cri-private`).
  * `modified_time` (Integer): The last modification time of the rule (Unix timestamp in milliseconds).
  * `namespace_name` (String): The namespace name to which the rule applies (e.g., `testtf`).
  * `next_time` (Integer): The next execution time of the rule (Unix timestamp in milliseconds).
  * `recent_pull_keep` (Integer): The number of days to retain recent pull records (used for pull-time-based retention policies).
  * `recent_push_keep` (Integer): The number of days to retain recent push records (used for push-time-based retention policies).
  * `repo_name` (String): The repository name to which the rule applies (valid when the rule operates at the repository level).
  * `retention_tag_count` (Integer): The threshold for the number of tags to retain (e.g., `30` means retaining the most recent 30 tags).
  * `rule_id` (String): The unique identifier of the rule (e.g., `cralr-42x00j4drgt4fjr3`).
  * `schedule_time` (String): The rule scheduling time configuration (e.g., `MANUAL` indicates manual triggering).
  * `tag_regexp` (String): The regular expression used to match tags (e.g., `release-v.*`).