---
layout: "apsarastack"
page_title: "Provider: apsarastack"
sidebar_current: "docs-apsarastack-index"
description: |-
  The ApsaraStack provider is used to interact with many resources supported by ApsaraStack. The provider needs to be configured with the proper credentials before it can be used.
---

# ApsaraStack Cloud Provider

The ApsaraStack Cloud provider is used to interact with the
many resources supported by ApsaraStack Cloud. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    apsarastack = {
      source = "apsara-stack/apsarastack"
      version = "1.0.8"
    }
  }
}

# Configure the ApsaraStack Provider
provider "apsarastack" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  insecure    =  true
  proxy      = "${var.proxy}"
  resource_group_set_name ="${var.resource_group_set_name}"
  domain = "${var.domain}"
  protocol = "HTTPS"
}


data "apsarastack_instance_types" "default" {
  cpu_core_count = 2
  memory_size    = 4
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}
# Create a web server
resource "apsarastack_instance" "web" {
  image_id              = "${data.apsarastack_images.default.images.0.id}"
  instance_type        = "${data.apsarastack_instance_types.default .instance_types.0.id}"
  system_disk_category = "cloud_efficiency"
  security_groups      = ["${apsarastack_security_group.default.id}"]
  instance_name        = "web"
  vswitch_id           = "vsw-abc12345"
}

# Create security group
resource "apsarastack_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id      = "vpc-abc12345"
}
```

## Authentication

The ApsaraStack provider accepts several ways to enter credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials

Static credentials can be provided by adding `access_key`, `secret_key` , `region` ,`insecure`,`proxy` and `domain` in-line in the
apsarastack provider block:

Usage:

```hcl
provider "apsarastack" {
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

You can provide your credentials via `APSARASTACK_ACCESS_KEY`,`APSARASTACK_SECRET_KEY`,
environment variables, representing your ApsaraStack access key and secret key respectively.
`APSARASTACK_PROXY`,`APSARASTACK_REGION` is also used, if applicable:

```hcl
provider "apsarastack" {
    endpoints {
         vpc = "${var.endpoints}"  
       }
    resource_group_set_name ="${var.resource_group_set_name}"
}
```
Usage:

```shell
$ export APSARASTACK_ACCESS_KEY="anaccesskey"
$ export APSARASTACK_SECRET_KEY="asecretkey"
$ export APSARASTACK_REGION="region"
$ export APSARASTACK_INSECURE= true
$ export APSARASTACK_PROXY= "http://IP:Port"
$ terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html)
(e.g. `alias` and `version`), the following arguments are supported in the ApsaraStack Cloud
 `provider` block:

* `access_key` - This is the ApsaraStack access key. It must be provided, but
  it can also be sourced from the `APSARASTACK_ACCESS_KEY` environment variable, or via
  a dynamic access key if `ecs_role_name` is specified.

* `secret_key` - This is the ApsaraStack secret key. It must be provided, but
  it can also be sourced from the `APSARASTACK_SECRET_KEY` environment variable, or via
  a dynamic secret key if `ecs_role_name` is specified.
  
* `region` - This is the ApsaraStack region. It must be provided, but
  it can also be sourced from the `APSARASTACK_REGION` environment variables.

* `insecure` - (Optional) Use this to Trust self-signed certificates. It's typically used to allow insecure connections.

* `resource_group_set_name` - (Optional) Use this to give resource_group_set_name for specific user organisation.

* `protocol` - (Optional) The Protocol of used by API request. Valid values: `HTTP` and `HTTPS`. Default to `HTTPS`.

* `proxy` -  (Optional) Use this to set proxy for ApsaraStack connection.

* `endpoints` - (Required) An `endpoints` block (documented below) to support apsarastack custom endpoints.

Nested `endpoints` block supports the following:
* `ecs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom ECS endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.

* `slb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SLB endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC and VPN endpoints.

* `ess` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom Autoscaling endpoints.

* `oss` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom OSS endpoints.


