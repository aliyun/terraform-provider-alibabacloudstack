---
subcategory: "Alikafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_sasl_acl"
sidebar_current: "docs-alibabacloudstack-resource-alikafka-sasl_acl"
description: |-
  编排Alikafka SASL ACL资源
---

# alibabacloudstack_alikafka_sasl_acl

使用Provider配置的凭证在指定的资源集下编排Alikafka SASL ACL资源。

## 示例用法

### 基础用法

```
variable "name" {
	default = "testalikafkasslacl"
}

variable "password" {
}

resource "alibabacloudstack_alikafka_topic" "default" {
  instance_id = "cluster-private-paas-default"
  topic = "${var.name}"
  remark = "topic-remark"
}


resource "alibabacloudstack_alikafka_sasl_user" "default" {
  instance_id = "cluster-private-paas-default"
  username = "${var.name}"
  password = var.password
  type     = "scram"
}

resource "alibabacloudstack_alikafka_sasl_acl" "default" {
    instance_id =               "cluster-private-paas-default"
    username =                 "${alibabacloudstack_alikafka_sasl_user.default.username}"
    acl_resource_type =         "Topic"
    acl_resource_name =         "${alibabacloudstack_alikafka_topic.default.topic}"
    acl_resource_pattern_type = "LITERAL"
    acl_operation_type =        "Write"
}

```



## 参数参考

以下参数被支持：

* `instance_id` - (必填，变更时重建) 拥有该组的ALIKAFKA实例ID。
* `username` - (必填，变更时重建) SASL用户的用户名。长度应在1到64个字符之间。用户应为已存在的SASL用户。
* `acl_resource_type` - (必填，变更时重建) 此ACL的资源类型。资源类型只能是"Topic"和"Group"。
* `acl_resource_name` - (必填，变更时重建) 此ACL的资源名称。资源名称应为一个主题或消费者组名称。
* `acl_resource_pattern_type` - (必填，变更时重建) 此ACL的资源模式类型。资源模式支持两种类型"LITERAL"和"PREFIXED"。"LITERAL": 字面名称定义了资源的完整名称。特殊通配符"*"可以用来表示任何名称的资源。"PREFIXED": 前缀名称定义了一个资源的前缀。
* `acl_operation_type` - (必填，变更时重建) 此ACL的操作类型。操作类型只能是"Write"和"Read"。
* `host` - (可选，变更时重建) ACL的主机。

## 属性参考

以下属性被导出：

* `id` - 资源提供的`key`。其值由 `<instance_id>:<username>:<acl_resource_type>:<acl_resource_name>:<acl_resource_pattern_type>:<acl_operation_type>` 组成。
* `host` - ACL的主机。