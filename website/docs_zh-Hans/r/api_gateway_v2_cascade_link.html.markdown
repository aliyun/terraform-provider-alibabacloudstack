---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_api_gateway_v2_cascade_link"
sidebar_current: "docs-AlibabacloudStack-api_gateway_v2_cascade_link"
description: |-
  创建和管理API网关v2版本的级联链路资源
---

# alibabacloudstack_api_gateway_v2_cascade_link

管理API网关v2版本的级联链路，用于建立源实例与级联实例之间的连接。

## 示例用法

### 基础配置

```hcl

variable "name" {
  default = "tf-testAccApiGwV237988"
}

resource "alibabacloudstack_api_gateway_v2_instance" "source" {
  instance_name      = "${var.name}-source"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_instance" "cascade" {
  instance_name      = "${var.name}-cascade"
  node_number        = 1
  instance_class     = "mini"
  broker_engine_type = "SCG"
  deploy_mode        = "custom"
}

resource "alibabacloudstack_api_gateway_v2_cascade_instance" "default" {
  instance_name       = var.name
  cascade_instance_id = alibabacloudstack_api_gateway_v2_instance.cascade.id
}


resource "alibabacloudstack_api_gateway_v2_cascade_link" "default" {
  source_instance_id      = alibabacloudstack_api_gateway_v2_instance.source.id
  source_instance_address = "10.17.94.180"
  cascade_instance_id     = alibabacloudstack_api_gateway_v2_cascade_instance.default.id
  link_name               = var.name
}
```

## 参数说明

支持以下参数：

* `cascade_instance_id` - (必填, 变更时重建) 级联实例ID。指定目标API网关实例的唯一标识符，用于建立级联关系。
* `link_name` - (必填, 变更时重建) 链路名称。自定义链路的名称，用于标识该级联链路。
* `source_instance_address` - (必填) 源实例地址。源服务实例的IP地址或域名，用于API请求路由。
* `source_instance_id` - (必填, 变更时重建) 源实例ID。源服务实例的唯一标识符，必须与源实例地址匹配。
* `cascade_instance_name` - (可选) 级联实例名称。级联API网关实例的显示名称，仅用于标识。
* `cascade_service_id` - (可选) 级联服务ID。关联的级联服务唯一标识符，用于服务级联管理。
* `source_instance_name` - (可选) 源实例名称。源服务实例的显示名称，仅用于标识。

## 属性说明

以下属性导出为资源状态的一部分：

* `id` - 级联链路的唯一ID，与`link_id`值相同。
* `link_id` - 级联链路的系统生成ID，用于唯一标识该资源。