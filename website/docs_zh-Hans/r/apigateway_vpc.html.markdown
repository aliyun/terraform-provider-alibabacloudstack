---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_vpc"
sidebar_current: "docs-Alibabacloudstack-apigateway-vpc"
description: |- 
  编排API网关下的VPC端口
---

# alibabacloudstack_api_gateway_vpc_access
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_apigateway_vpc`

使用Provider配置的凭证在指定的资源集下编排API网关下的VPC端口。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccApiGatewayVpcAccess-2159202"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.ids.0
}

data "alibabacloudstack_images" "default" {
  name_regex = "^ubuntu"
  most_recent = true
  owners = "system"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alibabacloudstack_zones.default.ids.0
}

resource "alibabacloudstack_security_group" "default" {
  name        = var.name
  description = "foo"
  vpc_id      = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  vswitch_id                  = alibabacloudstack_vswitch.default.id
  image_id                   = data.alibabacloudstack_images.default.images.0.id
  instance_type              = data.alibabacloudstack_instance_types.default.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_max_bandwidth_out = 5
  security_groups            = [alibabacloudstack_security_group.default.id]
  instance_name              = var.name
}

resource "alibabacloudstack_api_gateway_vpc_access" "default" {
  name         = var.name
  vpc_id       = alibabacloudstack_vpc.default.id
  instance_id  = alibabacloudstack_instance.default.id
  port         = "8080"
}
```

## 参数参考

支持以下参数：

* `name` - (必填，变更时重建) VPC授权的名称。它必须在用户资源范围内唯一。
* `vpc_id` - (必填，变更时重建) 您要授权API网关访问的VPC的ID。
* `instance_id` - (必填，变更时重建) 您要授权API网关访问的VPC中的ECS或服务器负载均衡器实例的ID。
* `port` - (必填，变更时重建) API网关应连接到实例上的端口号。有效值范围为1到65535。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `id` - API网关的VPC授权ID。它由`vpc_id`、`instance_id`和`port`的组合组成。
```