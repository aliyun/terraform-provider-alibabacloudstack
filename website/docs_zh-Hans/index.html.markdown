---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  The AlibabacloudStack provider is used to interact with many resources supported by AlibabacloudStack. The provider needs to be configured with the proper credentials before it can be used.
---

# AlibabacloudStack Provider

AlibabacloudStack Provider用于Terraform管理阿里云私有云平台下的多种资源。在使用前，需配置该Provider访问云平台的正确凭证。
使用左侧导航栏查看支持的列表。

## 示例代码

```hcl
terraform {
  required_providers {
    alibabacloudstack = {
      source = "aliyun/alibabacloudstack"
      #version = "3.16.0"
    }
  }
}

# 配置 AlibabacloudStack Provider
provider "alibabacloudstack" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  insecure    =  true
  proxy      = "${var.proxy}"
  resource_group_set_name ="${var.resource_group_set_name}"
  domain = "${var.domain}"
  protocol = "HTTPS"
}


data "alibabacloudstack_instance_types" "default" {
  cpu_core_count = 2
  memory_size    = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
# 创建 Web 服务器实例
resource "alibabacloudstack_instance" "web" {
  image_id              = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${data.alibabacloudstack_instance_types.default .instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${alibabacloudstack_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id           = "vsw-abc12345"
}

# 创建安全组
resource "alibabacloudstack_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id      = "vpc-abc12345"
}
```

## 凭证配置

AlibabacloudStack Provider支持以下凭证认证方式：

- 静态凭证
- 环境变量


### 静态凭证

在模板的alibabacloudstack provider 代码块中配置 `access_key`, `secret_key` , `region` ,`insecure`,`proxy` 和 `domain` 等信息，
为AlibabacloudStack Provider平台访问凭证。

Usage:

```hcl
provider "alibabacloudstack" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  insecure    =  true
  proxy      = "${var.proxy}"
  resource_group_set_name ="${var.resource_group_set_name}"
  endpoints {
     vpc = "${var.endpoints}"  
   }
}

```

### 环境变量

通过环境变量 `ALIBABACLOUDSTACK_ACCESS_KEY`,`ALIBABACLOUDSTACK_SECRET_KEY`,
为AlibabacloudStack Provider提供平台访问凭证。
同时也可以配置`ALIBABACLOUDSTACK_PROXY`,`ALIBABACLOUDSTACK_REGION` 等环境变量。

```hcl
provider "alibabacloudstack" {
    endpoints {
         vpc = "${var.endpoints}"
       }
    resource_group_set_name ="${var.resource_group_set_name}"
}
```
Usage:

```shell
$ export ALIBABACLOUDSTACK_ACCESS_KEY="anaccesskey"
$ export ALIBABACLOUDSTACK_SECRET_KEY="asecretkey"
$ export ALIBABACLOUDSTACK_REGION="region"
$ export ALIBABACLOUDSTACK_INSECURE= true
$ export ALIBABACLOUDSTACK_PROXY= "http://IP:Port"
$ terraform plan
```

## 参数说明

除[Terraform通用参数](https://www.terraform.io/docs/configuration/providers.html)(如. `alias` 和 `version`)外, 
AlibabacloudStack 的 provider 配置块支持以下参数：

* `access_key` - (必填) 访问密钥。也可通过ALIBABACLOUDSTACK_ACCESS_KEY环境变量获取。

* `secret_key` - (必填) 秘密密钥。也可通过ALIBABACLOUDSTACK_SECRET_KEY环境变量获取。
  
* `region` - (必填) 专有云区域信息。也可通过ALIBABACLOUDSTACK_REGION环境变量获取。

* `insecure` - (可选) 允许自签名证书，用于启用不安全连接。

* `department` - (可选) 指定编排资源所隶属的组织的ID。为配置时会通过`resource_group_set_name`查找。

* `resource_group` - (可选) 指定编排资源所隶属的资源集的ID。为配置时会通过`resource_group_set_name`查找。

* `resource_group_set_name` - (可选) 指定编排资源所隶属的资源集的名称。当`resource_group_set_name`不唯一或未配置时需要配置`department`和`resource_group`。

* `protocol` - (可选) API 请求协议。可选值：HTTP/HTTPS，默认为HTTPS。

* `proxy` - (可选) 设置AlibabacloudStack连接的代理。

* `endpoints` - (可选) 自定义端点配置块，用于覆盖默认区域端点。

嵌套 endpoints 块支持以下参数：

* `ecs` - (可选) 覆盖默认 ecs 服务地址配置。

* `rds` - (可选) 覆盖默认 rds 服务地址配置。

* `slb` - (可选) 覆盖默认 slb 服务地址配置。

* `vpc` - (可选) 覆盖默认 vpc 服务地址配置。

* `ess` - (可选) 覆盖默认 ess 服务地址配置。

* `oss` - (可选) 覆盖默认 oss 服务地址配置。


