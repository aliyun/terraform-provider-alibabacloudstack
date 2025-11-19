---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_cascade_instance"
sidebar_current: "docs-Alibabacloudstack-api-gateway-v2-cascade-instance"
description: |-
  管理API网关V2级联网关实例
---

# alibabacloudstack_api_gateway_v2_cascade_instance

管理API网关V2级联网关实例。该资源用于创建和管理级联API网关实例，支持指定级联实例ID和实例名称。**注意：该资源创建后不可修改，任何参数变更将触发资源重建**。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testAccApiGwV222547"
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name      = "${var.name}source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}



resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = var.name
  cascade_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
}
```

## 参数说明

支持以下参数，按类型排序（必填 → 变更时重建 → 可选 → 过时）：

* `cascade_instance_id` - (必填, 变更时重建) 级联实例ID。用于关联底层级联实例，必须与API网关服务中已存在的级联实例ID一致。创建后无法修改，变更将导致资源重建。
* `instance_name` - (必填, 变更时重建) 实例名称。长度为1-128字符，可包含字母、数字、短划线（-）和下划线（_）。创建后无法修改，变更将导致资源重建。
