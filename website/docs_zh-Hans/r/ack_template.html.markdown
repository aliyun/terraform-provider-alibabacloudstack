---
subcategory: "Kubernetes容器监控"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_template"
sidebar_current: "docs-Alibabacloudstack-resource-ack-template"
description: |-
  提供一种 ACK 模板资源。
---

# alibabacloudstack_ack_template

提供一种 ACK 模板资源。

## 示例用法

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
  description   = "一个示例 ACK 模板"
  template_type = "kubernetes"
}
```

## 参数参考

支持以下参数：

* `name` - (必填, ForceNew) 编排模板名称。名称长度为 1~63 个字符，可包含数字、字母和短划线（-），不能以短划线（-）开头。
* `template` - (必填) YAML 格式的模板内容。
* `description` - (可选) 模板描述。
* `template_type` - (必填, ForceNew) 模板类型。设置为 `kubernetes` 时，模板会在控制台的模板页面中显示。建议设置为 `kubernetes`。

## 属性参考

导出以下属性：

* `id` - 资源 ID。
* `template_id` - 编排模板 ID。
* `template_with_hist_id` - 模板唯一标识 ID。模板更新后该值保持不变。
* `template_hash_code_version` - 模板的哈希码版本。
* `created` - 模板创建时间。
* `acl` - 模板的访问控制策略。
* `version` - 模板版本。
* `tags` - 模板标签。
* `ali_uid` - 阿里云 UID。
* `updated` - 模板更新时间。

## Import

ACK 模板可以使用模板 ID 导入，例如：

```
$ terraform import alibabacloudstack_ack_template.example 72d20cf8-a533-4ea9-a10d-e7630d3d2708
```