---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_template"
sidebar_current: "docs-alibabacloudstack-resource-ack-template"
description: |-
  Provides a ACK Template resource.
---

# alibabacloudstack_ack_template

Provides a ACK Template resource.

## Example Usage

```hcl
resource "alibabacloudstack_ack_template" "example" {
  name          = "example-template"
  template      = <<EOF
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	labels:
		vsw: test
	name: nginx-deployment-basic
	namespace: default
	spec:
	replicas: 1
	selector:
		matchLabels:
		vsw: test
	template:
		metadata:
		labels:
			vsw: test
		spec:
		containers:
			- command:
				- sleep
				- '10000'
			image: >-
				registry.acs.inter.env128.shuguang.com/acs/busybox:1.33.1
			imagePullPolicy: IfNotPresent
			name: vsw
EOF
  description   = "An example ACK template"
  template_type = "kubernetes"
}
```
## Argument Reference
The following arguments are supported:

* `template` - (Required) The template content.
* `name` - (Required, ForceNew) The name of the template.
* `description` - (Optional) The description of the template.
* `template_type` - (Required, ForceNew) The type of the template.
Attributes Reference
The following attributes are exported:

* `id` - The ID of the template.
* `template_id` - The ID of the template.
* `template_with_hist_id` - The template ID with history.
* `template_hash_code_version` - The hash code version of the * template.
* `created` - The creation time of the template.
* `acl` - The ACL of the template.
* `version` - The version of the template.
* `tags` - The tags of the template.
* `ali_uid` - The Alibaba Cloud UID.
* `updated` - The last update time of the template.