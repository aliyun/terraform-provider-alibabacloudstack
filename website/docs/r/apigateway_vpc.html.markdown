---
subcategory: "ApiGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_apigateway_vpc"
sidebar_current: "docs-Alibabacloudstack-apigateway-vpc"
description: |- 
  Provides a apigateway Vpc resource.
---

# alibabacloudstack_api_gateway_vpc_access
-> **NOTE:** Alias name has: `alibabacloudstack_apigateway_vpc`

Provides an API Gateway VPC access resource. This authorizes the API Gateway to access your VPC instances.

For information about Api Gateway VPC and how to use it, see [Set VPC Access](https://help.aliyun.com/document_detail/400343.html?spm=5176.10695662.1996646101.searchclickresult.67be328fV80qXE).

-> **NOTE:** Terraform will auto build VPC authorization while it uses `alibabacloudstack_api_gateway_vpc_access` to build VPC.

## Example Usage

Basic Usage

```hcl
variable "name" {
  default = "tf-testAccApiGatewayVpcAccess-2159202"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone           = data.alibabacloudstack_zones.default.ids.0
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

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the VPC authorization. It must be unique within the scope of the user's resources.
* `vpc_id` - (Required, ForceNew) The ID of the VPC that you want to authorize the API Gateway to access.
* `instance_id` - (Required, ForceNew) The ID of the ECS or Server Load Balancer instance in the VPC that you want to authorize the API Gateway to access.
* `port` - (Required, ForceNew) The port number on the instance that the API Gateway should connect to. Valid values range from 1 to 65535.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC authorization for the API Gateway. It is composed of the combination of `vpc_id`, `instance_id`, and `port`.

### Import

API Gateway VPC access can be imported using the combined ID format `<VPC_ID>:<INSTANCE_ID>:<PORT>`, e.g.,

```sh
$ terraform import alibabacloudstack_api_gateway_vpc_access.example "vpc-aswcj19ajsz:i-ajdjfsdlf:8080"
```