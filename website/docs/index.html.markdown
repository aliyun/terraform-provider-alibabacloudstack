---
layout: "alibabacloudstack"
page_title: "Provider: alibabacloudstack"
sidebar_current: "docs-alibabacloudstack-index"
description: |-
  The AlibabacloudStack provider is used to interact with many resources supported by AlibabacloudStack. The provider needs to be configured with the proper credentials before it can be used.
---

# AlibabacloudStack Cloud Provider

The AlibabacloudStack Cloud provider is used to interact with the
many resources supported by AlibabacloudStack Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    alibabacloudstack = {
      source = "apsara-stack/alibabacloudstack"
      version = "1.0.8"
    }
  }
}

# Configure the AlibabacloudStack Provider
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
# Create a web server
resource "alibabacloudstack_instance" "web" {
  image_id              = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${data.alibabacloudstack_instance_types.default .instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${alibabacloudstack_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id           = "vsw-abc12345"
}

# Create security group
resource "alibabacloudstack_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id      = "vpc-abc12345"
}
```

## Authentication

The AlibabacloudStack provider accepts several ways to enter credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding `access_key`, `secret_key` , `region` ,`insecure`,`proxy` and `domain` in-line in the
alibabacloudstack provider block:

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

### Environment variables

You can provide your credentials via `ALIBABACLOUDSTACK_ACCESS_KEY`,`ALIBABACLOUDSTACK_SECRET_KEY`,
environment variables, representing your AlibabacloudStack access key and secret key respectively.
`ALIBABACLOUDSTACK_PROXY`,`ALIBABACLOUDSTACK_REGION` is also used, if applicable:

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

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the AlibabacloudStack Cloud
 `provider` block:

* `access_key` - This is the AlibabacloudStack access key. It must be provided, but
  it can also be sourced from the `ALIBABACLOUDSTACK_ACCESS_KEY` environment variable, or via
  a dynamic access key if `ecs_role_name` is specified.

* `secret_key` - This is the AlibabacloudStack secret key. It must be provided, but
  it can also be sourced from the `ALIBABACLOUDSTACK_SECRET_KEY` environment variable, or via
  a dynamic secret key if `ecs_role_name` is specified.
  
* `region` - This is the AlibabacloudStack region. It must be provided, but
  it can also be sourced from the `ALIBABACLOUDSTACK_REGION` environment variables.

* `insecure` - (Optional) Use this to Trust self-signed certificates. It's typically used to allow insecure connections.

* `resource_group_set_name` - (Optional) Use this to give resource_group_set_name for specific user organisation.

* `protocol` - (Optional) The Protocol of used by API request. Valid values: `HTTP` and `HTTPS`. Default to `HTTPS`.

* `proxy` -  (Optional) Use this to set proxy for AlibabacloudStack connection.

* `endpoints` - (Required) An `endpoints` block (documented below) to support alibabacloudstack custom endpoints.

Nested `endpoints` block supports the following:
* `ecs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.

* `slb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.

* `ess` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.

* `oss` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.


