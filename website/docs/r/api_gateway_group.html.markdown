---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_group"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-group"
description: |-
  Provides a Alibabacloudstack Api Gateway Group Resource.
---

# alibabacloudstack_api_gateway_group

Provides an api group resource.To create an API, you must firstly create a group which is a basic attribute of the API.

For information about Api Gateway Group and how to use it, see [Create An Api Group](https://help.aliyun.com/apsara/enterprise/v_3_14_0_20210519/apigateway/apsara-developer-guide/creates-an-api-group--1.html?spm=a2c4g.14484438.10001.139)

-> **NOTE:** Terraform will auto build api group while it uses `alibabacloudstack_api_gateway_group` to build api group.

## Example Usage

Basic Usage

```
	variable "name" {
	  default = "tf_testAccGroup_4663614"
	}

	variable "description" {
	  default = "tf_testAcc api gateway description"
	}
	

resource "alibabacloudstack_api_gateway_group" "default" {
  name = "${var.name}"
  description = "${var.description}"
}

```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the api gateway group. Defaults to null.
* `description` - (Required) The description of the api gateway group. Defaults to null.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the api group of api gateway.
* `sub_domain` - (Available in 1.69.0+)	Second-level domain name automatically assigned to the API group.
* `vpc_domain` - (Available in 1.69.0+)	Second-level VPC domain name automatically assigned to the API group.

## Import

Api gateway group can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_api_gateway_group.example "ab2351f2ce904edaa8d92a0510832b91"
```
