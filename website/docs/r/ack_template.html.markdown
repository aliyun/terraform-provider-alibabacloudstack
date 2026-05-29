---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_template"
sidebar_current: "docs-Alibabacloudstack-resource-ack-template"
description: |-
  Provides a ACK Template resource.
---

# alibabacloudstack_ack_template

Provides a ACK Template resource.

## Example Usage

```hcl
resource "alibabacloudstack_ack_template" "example" {
  name          = "example-template"
  template      = <<-EOT
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      labels:
        app: test
      name: nginx-deployment-basic
      namespace: default
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: test
      template:
        metadata:
          labels:
            app: test
        spec:
          containers:
            - command:
                - sleep
                - "10000"
              image: registry.acs.inter.env128.shuguang.com/acs/busybox:1.33.1
              imagePullPolicy: IfNotPresent
              name: test
  EOT
  description   = "An example ACK template"
  template_type = "kubernetes"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the orchestration template. The name must be 1 to 63 characters in length, and can contain digits, letters, and hyphens (-). It cannot start with a hyphen (-).
* `template` - (Required) The template content in the YAML format.
* `description` - (Optional) The description of the template.
* `template_type` - (Required, ForceNew) The type of template. If the parameter is set to `kubernetes`, the template is displayed on the Templates page in the console. We recommend that you set the parameter to `kubernetes`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `template_id` - The ID of the orchestration template.
* `template_with_hist_id` - The unique ID of the template. The value remains unchanged after the template is updated.
* `template_hash_code_version` - The hash code version of the template.
* `created` - The time when the template was created.
* `acl` - The access control policy of the template.
* `version` - The version of the template.
* `tags` - The labels of the template.
* `ali_uid` - The Alibaba Cloud UID.
* `updated` - The time when the template was updated.

## Import

ACK Template can be imported using the template ID, e.g.

```
$ terraform import alibabacloudstack_ack_template.example 72d20cf8-a533-4ea9-a10d-e7630d3d2708
```