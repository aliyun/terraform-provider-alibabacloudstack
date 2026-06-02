---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_account"
description: |-
  Provides a PolarDBX Account resource.
---

# alibabacloudstack_polardbx_account

Provides a PolarDBX Account resource.

## Example Usage

```hcl
variable "name" {
  default = "accdbbind90366"
}

variable "password" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  engine_version = "5.7"
  storage = "50"
  network_type = "vpc"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  cn_node_class = "polarx.x4.medium.2e"
  cn_node_count = "2"
  dn_node_class = "mysql.n4.medium.25"
  dn_node_count = "2"
}

resource "alibabacloudstack_polardbx_account" "default" {
  instance_id  = alibabacloudstack_polardbx_instance.default.id
  account_name = var.name
  password     = var.password
  description  = "Normal user"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the PolarDBX instance. Changing this parameter will force the resource to be recreated.
* `account_name` - (Required, ForceNew) The name of the account. Changing this parameter will force the resource to be recreated. The account name must meet the following requirements:
  * Start with a lowercase letter and end with a letter or digit
  * Consist of lowercase letters, digits, or underscores
  * Be 2 to 16 characters in length
  * Cannot use reserved usernames such as root and admin
* `password` - (Required) The password of the account. This parameter is sensitive and will be masked in Terraform state files.
* `description` - (Optional) The description of the account. The length is 2 to 256 characters and cannot start with `http://` or `https://`. This attribute is returned by the API and cannot be manually set.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the resource, formatted as `<instance_id>:<account_name>`
* `description` - The description of the account
