---
subcategory: "Realtime Compute for Apache (Flink)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack:alibabacloudstack_flink_namespace"
sidebar_current: "docs-alibabacloudstack-resource-flink-namesapce"
description: |-
  Provides a Alibabacloudstack resource to manage Flink namespaces.
---

# alibabacloudstack_flink_namespace

This resource will help you to manager Flink namespaces.


## Example Usage

Basic Usage

```
variable name{
 default = "<Your NameSpace Name>"
}

resource "alibabacloudstack_ascm_user" "default" {
  display_name = var.name
  mobile_nation_code = "86"
  login_name = var.name
  login_policy_id = "1"
  role_ids = [
               "8",
               "9"
             ]
  cellphone_number = "13612345678"
  email = "${var.name}@gmail.com"
}


resource "alibabacloudstack_flink_namespace" "default" {
  owner_uid = "${alibabacloudstack_ascm_user.default.user_uid}"
  name = var.name"
  cu = "1"
  cpu_type = "Intel"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of Flink Registry namespace. 
* `cu` - (Required) Integer. Guaranteed Resources for Cpu.
* `cpu_type` - (Required) The CPU type of the resource. Valid values: `intel`.
* `owner_uid` - (Optional, ForceNew) Owner ID for this Namesapce.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Fink namespace. The value is same as its name.


## Import

Flink namespace can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_flink_namespace.default namespace_name
```