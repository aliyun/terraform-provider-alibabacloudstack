---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_template"
sidebar_current: "docs-alibabacloudstack-resource-ack-template"
description: |-
  提供一种 ACK 模板资源。
---

# alibabacloudstack_ack_template

提供一种 ACK 模板资源。

## 示例用法

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
  description   = "一个示例 ACK 模板"
  template_type = "kubernetes"
}
```
## 参数参考
## 支持以下参数：

* `template` - (必填) 模板内容。
* `name` - (必填, ForceNew) 模板名称。
* `description` - (可选) 模板描述。
* `template_type` - (必填, ForceNew) 模板类型。
属性参考
导出以下属性：

* `id` - 模板 ID。
* `template_id` - 模板 ID。
* `template_with_hist_id` - 带历史记录的模板 ID。
* `template_hash_code_version` - 模板的哈希码版本。
* `created` - 模板创建时间。
* `acl` - 模板的 ACL。
* `version` - 模板版本。
* `tags` - 模板标签。
* `ali_uid` - 阿里云 UID。
* `updated` - 模板最后更新时间。